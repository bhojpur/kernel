package photon

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

	"github.com/bhojpur/kernel/pkg/compilers/rump"
	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/photon-controller-go-sdk/photon"
)

var timeout = time.Second * 10

const PhotonKernelInstanceListener = "PhotonKernelInstanceListener"
const instanceListenerPrefix = "kernel_photon"

func (p *PhotonProvider) DeployInstanceListener(config config.Photon) error {
	logrus.Infof("checking if Bhojpur Kernel instance listener is alive...")
	if instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout); err == nil {
		logrus.Infof("instance listener is alive with IP %s", instanceListenerIp)
		return nil
	}
	logrus.Infof("cannot contact Bhojpur Kernel instance listener... cleaning up previous if it exists..")
	vms, err := p.client.Projects.GetVMs(config.ProjectId, &photon.VmGetOptions{
		Name: PhotonKernelInstanceListener,
	})
	if err != nil {
		return errors.New("getting photon vm list", err)
	}
	for _, vm := range vms.Items {
		if vm.Name == PhotonKernelInstanceListener {
			task, err := p.client.VMs.Stop(vm.ID)
			if err != nil {
				return errors.New("Stopping vm", err)
			}
			task, _ = p.waitForTaskSuccess(task)
			p.client.VMs.Delete(vm.ID)
			break
		}
	}
	logrus.Infof("compiling new Bhojpur Kernel instance listener")
	sourceDir, err := ioutil.TempDir("", "photon.instancelistener.")
	if err != nil {
		return errors.New("creating temp dir for instance listener source", err)
	}
	defer os.RemoveAll(sourceDir)
	rawImage, err := common.CompileInstanceListener(sourceDir, instanceListenerPrefix, "compilers-rump-go-hw", rump.CreateImageVmware, false)
	if err != nil {
		return errors.New("compiling instance listener source to unikernel", err)
	}
	defer os.Remove(rawImage.LocalImagePath)
	logrus.Infof("staging new instance listener image")
	//delete old image if it exists
	if err := p.deleteOldImage(); err != nil {
		logrus.Warn("failed removing previous image", err)
	}

	params := types.StageImageParams{
		Name:     PhotonKernelInstanceListener,
		RawImage: rawImage,
		Force:    true,
	}
	image, err := p.Stage(params)
	if err != nil {
		return errors.New("building bootable virtualbox image for instsance listener", err)
	}
	defer func() {
		if err != nil {
			p.DeleteImage(image.Id, true)
		}
	}()

	if err := p.runInstanceListener(image); err != nil {
		return errors.New("launching instance of instance listener", err)
	}
	return nil
}

func (p *PhotonProvider) deleteOldImage() error {
	if err := p.DeleteImage(PhotonKernelInstanceListener, true); err != nil {
		return nil
	}
	images, err := p.client.Images.GetAll()
	if err != nil {
		return errors.New("retrieving photon image list", err)
	}
	for _, image := range images.Items {
		if image.Name == PhotonKernelInstanceListener {
			task, err := p.client.Images.Delete(image.ID)
			if err != nil {
				return errors.New("Delete image", err)
			}
			_, err = p.waitForTaskSuccess(task)
			if err != nil {
				return errors.New("Delete image", err)
			}
		}
	}
	return errors.New("previous image not found", err)
}

func (p *PhotonProvider) runInstanceListener(image *types.Image) (err error) {
	vmflavor, err := p.getKernelFlavor("vm")
	if err != nil {
		return errors.New("can't get vm flavor", err)
	}

	diskflavor, err := p.getKernelFlavor("ephemeral-disk")
	if err != nil {
		return errors.New("can't get disk flavor", err)
	}

	disk := photon.AttachedDisk{
		Flavor:   diskflavor.Name,
		Kind:     "ephemeral-disk",
		Name:     "bootdisk-" + image.Id,
		BootDisk: true,
	}

	vmspec := &photon.VmCreateSpec{
		Flavor:        vmflavor.Name,
		SourceImageID: image.Id,
		Name:          PhotonKernelInstanceListener,
		Affinities:    nil,
		AttachedDisks: []photon.AttachedDisk{disk},
	}

	task, err := p.client.Projects.CreateVM(p.projectId, vmspec)

	if err != nil {
		return errors.New("Creating vm", err)
	}

	task, err = p.waitForTaskSuccess(task)

	if err != nil {
		return errors.New("Waiting for create vm", err)
	}

	instanceId := task.Entity.ID
	task, err = p.client.VMs.Start(instanceId)
	if err != nil {
		return errors.New("Starting vm", err)
	}

	task, err = p.waitForTaskSuccess(task)
	if err != nil {
		return errors.New("Starting vm", err)
	}

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, time.Minute*5)
	if err != nil {
		return errors.New("failed to retrieve instance listener ip. is Bhojpur Kernel instance listener running?", err)
	}

	instance := &types.Instance{
		Id:             instanceId,
		Name:           PhotonKernelInstanceListener,
		State:          types.InstanceState_Running,
		IpAddress:      instanceListenerIp,
		Infrastructure: types.Infrastructure_PHOTON,
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
