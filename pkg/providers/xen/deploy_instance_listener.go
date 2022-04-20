package xen

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

	"fmt"

	"github.com/bhojpur/kernel/pkg/compilers/rump"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/xen/xenclient"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

const (
	instanceListenerPrefix    = "bhojpur_xen"
	XenKernelInstanceListener = "XenKernelInstanceListener"
)

var timeout = time.Second * 10
var instanceListenerData = "InstanceListenerData"

func (p *XenProvider) deployInstanceListener() error {
	logrus.Infof("checking if Bhojpur Kernel instance listener is alive...")
	if instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout); err == nil {
		logrus.Infof("Bhojpur Kernel instance listener is alive with IP %s", instanceListenerIp)
		return nil
	}
	logrus.Infof("cannot contact Bhojpur Kernel instance listener... cleaning up previous if it exists..")
	p.client.DestroyVm(XenKernelInstanceListener)
	logrus.Infof("compiling new Bhojpur Kernel instance listener")
	sourceDir, err := ioutil.TempDir("", "xen.instancelistener.")
	if err != nil {
		return errors.New("creating temp dir for Bhojpur Kernel instance listener source", err)
	}
	defer os.RemoveAll(sourceDir)
	rawImage, err := common.CompileInstanceListener(sourceDir, instanceListenerPrefix, "compilers-rump-go-xen", rump.CreateImageXen, false)
	if err != nil {
		return errors.New("compiling instance listener source to unikernel", err)
	}
	defer os.Remove(rawImage.LocalImagePath)
	logrus.Infof("staging new Bhojpur Kernel instance listener image")
	os.RemoveAll(getImagePath(XenKernelInstanceListener))
	params := types.StageImageParams{
		Name:     XenKernelInstanceListener,
		RawImage: rawImage,
		Force:    true,
	}
	image, err := p.Stage(params)
	if err != nil {
		return errors.New("building bootable xen image for Bhojpur Kernel instsance listener", err)
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

func (p *XenProvider) runInstanceListener(image *types.Image) (err error) {
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
		defer os.RemoveAll(imagePath)
		createVolumeParams := types.CreateVolumeParams{
			Name:      instanceListenerData,
			ImagePath: imagePath,
		}

		instanceListenerVol, err = p.CreateVolume(createVolumeParams)
		if err != nil {
			return errors.New("creating data vol for instance listener", err)
		}
		defer func() {
			if err != nil {
				p.DeleteVolume(instanceListenerVol.Id, true)
			}
		}()
	}

	defer func() {
		if err != nil {
			logrus.WithError(err).Warnf("error encountered, ensuring vm and disks are destroyed")
			p.DetachVolume(instanceListenerVol.Id)
			p.client.DestroyVm(XenKernelInstanceListener)
			os.RemoveAll(getInstanceDir(XenKernelInstanceListener))
			if newVolume {
				os.RemoveAll(getVolumePath(instanceListenerData))
			}
		}
	}()

	logrus.Debugf("creating xen vm")

	xenParams := xenclient.CreateVmParams{
		Name:      XenKernelInstanceListener,
		Memory:    image.RunSpec.DefaultInstanceMemory,
		BootImage: getImagePath(image.Name),
		VmDir:     getInstanceDir(XenKernelInstanceListener),
		DataVolumes: []xenclient.VolumeConfig{
			xenclient.VolumeConfig{
				ImagePath:  getVolumePath(instanceListenerVol.Name),
				DeviceName: "sdb1",
			},
		},
	}

	os.MkdirAll(getInstanceDir(XenKernelInstanceListener), 0755)

	if err := p.client.CreateVm(xenParams); err != nil {
		return errors.New("creating vm", err)
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

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, time.Minute*5)
	if err != nil {
		return errors.New("failed to retrieve instance listener ip. is Bhojpur Kernel instance listener running?", err)
	}

	doms, err := p.client.ListVms()
	if err != nil {
		return errors.New("getting vm info from xen", err)
	}
	instanceId := XenKernelInstanceListener
	for _, d := range doms {
		if d.Config.CInfo.Name == XenKernelInstanceListener {
			instanceId = fmt.Sprintf("%d", d.Domid)
			break
		}
	}

	instance := &types.Instance{
		Id:             instanceId,
		Name:           XenKernelInstanceListener,
		State:          types.InstanceState_Pending,
		IpAddress:      instanceListenerIp,
		Infrastructure: types.Infrastructure_VIRTUALBOX,
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
