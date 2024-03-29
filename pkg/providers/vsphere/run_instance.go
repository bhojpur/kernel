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
	"strings"
	"time"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/sirupsen/logrus"
)

func (p *VsphereProvider) RunInstance(params types.RunInstanceParams) (_ *types.Instance, err error) {
	logrus.WithFields(logrus.Fields{
		"image-id": params.ImageId,
		"mounts":   params.MntPointsToVolumeIds,
		"env":      params.Env,
	}).Infof("running instance %s", params.Name)

	if _, err := p.GetInstance(params.Name); err == nil {
		return nil, errors.New("instance with name "+params.Name+" already exists. virtualbox provider requires unique names for instances", nil)
	}

	image, err := p.GetImage(params.ImageId)
	if err != nil {
		return nil, errors.New("getting image", err)
	}

	if err := common.VerifyMntsInput(p, image, params.MntPointsToVolumeIds); err != nil {
		return nil, errors.New("invalid mapping for volume", err)
	}

	instanceDir := getInstanceDatastoreDir(params.Name)

	portsUsed := []int{}

	c := p.getClient()

	defer func() {
		if err != nil {
			if params.NoCleanup {
				logrus.Warnf("because --no-cleanup flag was provided, not cleaning up failed instance %s001", params.Name)
				return
			}
			logrus.WithError(err).Warnf("error encountered, ensuring vm and disks are destroyed")
			c.PowerOffVm(params.Name)
			for _, portUsed := range portsUsed {
				c.DetachDisk(params.Name, portUsed, image.RunSpec.StorageDriver)
			}
			c.DestroyVm(params.Name)
			c.Rmdir(instanceDir)
		}
	}()

	logrus.Debugf("creating vsphere vm")

	//if not set, use default
	if params.InstanceMemory <= 0 {
		params.InstanceMemory = image.RunSpec.DefaultInstanceMemory
	}

	if err := c.CreateVm(params.Name, params.InstanceMemory, image.RunSpec.VsphereNetworkType, p.config.NetworkLabel); err != nil {
		return nil, errors.New("creating vm", err)
	}

	logrus.Debugf("powering on vm to assign mac addr")
	if err := c.PowerOnVm(params.Name); err != nil {
		return nil, errors.New("failed to power on vm to assign mac addr", err)
	}

	vm, err := c.GetVm(params.Name)
	if err != nil {
		return nil, errors.New("failed to retrieve vm info after create", err)
	}

	macAddr := ""
	if vm.Config.Hardware.Device != nil {
		for _, device := range vm.Config.Hardware.Device {
			if len(device.MacAddress) > 0 {
				macAddr = device.MacAddress
				break
			}
		}
	}
	if macAddr == "" {
		logrus.WithFields(logrus.Fields{"vm": vm}).Warnf("vm found, cannot identify mac addr")
		return nil, errors.New("could not find mac addr on vm", nil)
	}
	if err := c.PowerOffVm(params.Name); err != nil {
		return nil, errors.New("failed to power off vm after retrieving mac addr", err)
	}

	logrus.Debugf("copying base boot vmdk to instance dir")
	instanceBootImagePath := instanceDir + "/boot.vmdk"
	if err := c.CopyFile(getImageDatastorePath(image.Name), instanceBootImagePath); err != nil {
		return nil, errors.New("copying base boot.vmdk", err)
	}
	if err := c.CopyFile(strings.TrimSuffix(getImageDatastorePath(image.Name), ".vmdk")+"-flat.vmdk", strings.TrimSuffix(instanceBootImagePath, ".vmdk")+"-flat.vmdk"); err != nil {
		return nil, errors.New("copying base boot-flat.vmdk", err)
	}
	if err := c.AttachDisk(params.Name, instanceBootImagePath, 0, image.RunSpec.StorageDriver); err != nil {
		return nil, errors.New("attaching boot vol to instance", err)
	}

	for mntPoint, volumeId := range params.MntPointsToVolumeIds {
		volume, err := p.GetVolume(volumeId)
		if err != nil {
			return nil, errors.New("getting volume", err)
		}
		controllerPort, err := common.GetControllerPortForMnt(image, mntPoint)
		if err != nil {
			return nil, errors.New("getting controller port for mnt point", err)
		}
		if err := c.AttachDisk(params.Name, getVolumeDatastorePath(volume.Name), controllerPort, image.RunSpec.StorageDriver); err != nil {
			return nil, errors.New("attaching disk to vm", err)
		}
		portsUsed = append(portsUsed, controllerPort)
	}

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout)
	if err != nil {
		return nil, errors.New("failed to retrieve instance listener ip. is Bhojpur Kernel instance listener running?", err)
	}

	logrus.Debugf("sending env to listener")
	if _, _, err := lxhttpclient.Post(instanceListenerIp+":3000", "/set_instance_env?mac_address="+macAddr, nil, params.Env); err != nil {
		return nil, errors.New("sending instance env to listener", err)
	}

	logrus.Debugf("powering on vm")
	if err := c.PowerOnVm(params.Name); err != nil {
		return nil, errors.New("powering on vm", err)
	}

	instanceId := vm.Config.UUID

	instance := &types.Instance{
		Id:             instanceId,
		Name:           params.Name,
		State:          types.InstanceState_Pending,
		IpAddress:      "",
		Infrastructure: types.Infrastructure_VSPHERE,
		ImageId:        image.Id,
		Created:        time.Now(),
	}

	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		instances[instance.Id] = instance
		return nil
	}); err != nil {
		return nil, errors.New("modifying instance map in state", err)
	}

	logrus.WithField("instance", instance).Infof("instance created successfully")

	return instance, nil
}
