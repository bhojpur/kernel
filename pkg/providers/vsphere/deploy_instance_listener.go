package vsphere

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"io/ioutil"
	"os"
	"time"

	"path/filepath"
	"strings"

	"github.com/bhojpur/kernel/pkg/compilers/rump"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

var timeout = time.Second * 10
var instanceListenerData = "InstanceListenerData"

func (p *VsphereProvider) deployInstanceListener() (err error) {
	logrus.Infof("checking if Bhojpur Kernel instance listener is alive...")
	if instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout); err == nil {
		logrus.Infof("Bhojpur Kernel instance listener is alive with IP %s", instanceListenerIp)
		return nil
	}
	logrus.Infof("cannot contact Bhojpur Kernel instance listener... cleaning up previous if it exists..")
	c := p.getClient()
	c.PowerOffVm(VsphereKernelInstanceListener)
	c.DestroyVm(VsphereKernelInstanceListener)
	logrus.Infof("compiling new Bhojpur Kernel instance listener")
	sourceDir, err := ioutil.TempDir("", "vsphereinstancelistener.")
	if err != nil {
		return errors.New("creating temp dir for Bhojpur Kernel instance listener source", err)
	}
	defer os.RemoveAll(sourceDir)
	rawImage, err := common.CompileInstanceListener(sourceDir, instanceListenerPrefix, "compilers-rump-go-hw", rump.CreateImageVmware, true)
	if err != nil {
		return errors.New("compiling Bhojpur Kernel instance listener source to unikernel", err)
	}
	logrus.Infof("staging new Bhojpur Kernel instance listener image")
	c.Rmdir(getImageDatastoreDir(VsphereKernelInstanceListener))
	params := types.StageImageParams{
		Name:     VsphereKernelInstanceListener,
		RawImage: rawImage,
		Force:    true,
	}
	image, err := p.Stage(params)
	if err != nil {
		return errors.New("building bootable vSphere image for Bhojpur Kernel instsance listener", err)
	}
	defer func() {
		if err != nil {
			p.DeleteImage(image.Id, true)
		}
	}()

	if err := p.runInstanceListener(image); err != nil {
		return errors.New("launching instance of Bhojpur Kernel instance listener", err)
	}
	return nil
}

func (p *VsphereProvider) runInstanceListener(image *types.Image) (err error) {
	logrus.WithFields(logrus.Fields{
		"image-id": image.Id,
	}).Infof("running instance of Bhojpur Kernel instance listener")

	newVolume := false
	instanceListenerVol, err := p.GetVolume(instanceListenerData)
	if err != nil {
		newVolume = true
		imagePath, err := util.BuildEmptyDataVolume(10)
		if err != nil {
			return errors.New("failed creating raw data volume", err)
		}
		defer os.RemoveAll(filepath.Dir(imagePath))

		params := types.CreateVolumeParams{
			Name:      instanceListenerData,
			ImagePath: imagePath,
		}
		instanceListenerVol, err = p.CreateVolume(params)
		if err != nil {
			return errors.New("creating data vol for instance listener", err)
		}
	}

	c := p.getClient()

	instanceDir := getInstanceDatastoreDir(VsphereKernelInstanceListener)
	defer func() {
		if err != nil {
			logrus.WithError(err).Warnf("error encountered, ensuring vm and disks are destroyed")
			p.DetachVolume(instanceListenerVol.Id)
			c.PowerOffVm(VsphereKernelInstanceListener)
			c.DestroyVm(VsphereKernelInstanceListener)
			c.Rmdir(instanceDir)
			if newVolume {
				p.DeleteVolume(instanceListenerVol.Id, true)
			}
			c.Rmdir(getVolumeDatastorePath(instanceListenerData))
		}
	}()

	logrus.Debugf("creating vSphere vm")

	if err := c.CreateVm(VsphereKernelInstanceListener, image.RunSpec.DefaultInstanceMemory, image.RunSpec.VsphereNetworkType, p.config.NetworkLabel); err != nil {
		return errors.New("creating vm", err)
	}

	logrus.Debugf("copying base boot image to instance dir")
	instanceBootImagePath := instanceDir + "/boot.vmdk"
	if err := c.CopyFile(getImageDatastorePath(image.Name), instanceBootImagePath); err != nil {
		return errors.New("copying boot.vmdk", err)
	}
	if err := c.CopyFile(strings.TrimSuffix(getImageDatastorePath(image.Name), ".vmdk")+"-flat.vmdk", strings.TrimSuffix(instanceBootImagePath, ".vmdk")+"-flat.vmdk"); err != nil {
		return errors.New("copying boot-flat.vmdk", err)
	}
	if err := c.AttachDisk(VsphereKernelInstanceListener, instanceBootImagePath, 0, image.RunSpec.StorageDriver); err != nil {
		return errors.New("attaching boot vol to instance", err)
	}

	controllerPort, err := common.GetControllerPortForMnt(image, "/data")
	if err != nil {
		return errors.New("getting controller port for mnt point", err)
	}
	logrus.Infof("attaching %s to %s on controller port %v", instanceListenerVol.Id, VsphereKernelInstanceListener, controllerPort)
	if err := c.AttachDisk(VsphereKernelInstanceListener, getVolumeDatastorePath(instanceListenerVol.Name), controllerPort, image.RunSpec.StorageDriver); err != nil {
		return errors.New("attaching disk to vm", err)
	}
	if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volume, ok := volumes[instanceListenerVol.Id]
		if !ok {
			return errors.New("no record of "+volume.Id+" in the state", nil)
		}
		volume.Attachment = instanceListenerVol.Id
		return nil
	}); err != nil {
		return errors.New("modifying volumes in state", err)
	}

	logrus.Debugf("powering on vm")
	if err := c.PowerOnVm(VsphereKernelInstanceListener); err != nil {
		return errors.New("powering on vm", err)
	}

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, time.Second*60)
	if err != nil {
		return errors.New("failed to retrieve Bhojpur Kernel instance listener IP. is Bhojpur Kernel instance listener running?", err)
	}

	vm, err := c.GetVm(VsphereKernelInstanceListener)
	if err != nil {
		return errors.New("getting vm info from vSphere", err)
	}

	instanceId := vm.Config.UUID
	instance := &types.Instance{
		Id:             instanceId,
		Name:           VsphereKernelInstanceListener,
		State:          types.InstanceState_Pending,
		IpAddress:      instanceListenerIp,
		Infrastructure: types.Infrastructure_VSPHERE,
		ImageId:        image.Id,
		Created:        time.Now(),
	}

	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		instances[instance.Id] = instance
		return nil
	}); err != nil {
		return errors.New("modifying instance map in state", err)
	}
	logrus.WithField("instance", instance).Infof("instance created successfully")

	return nil
}
