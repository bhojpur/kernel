package virtualbox

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
	"path/filepath"
	"strconv"

	"github.com/bhojpur/kernel/pkg/providers/virtualbox/virtualboxclient"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VirtualboxProvider) DetachVolume(id string) error {
	volume, err := p.GetVolume(id)
	if err != nil {
		return errors.New("retrieving volume "+id, err)
	}
	if volume.Attachment == "" {
		return errors.New("volume has no attachment", nil)
	}
	instanceId := volume.Attachment
	instance, err := p.GetInstance(instanceId)
	if err != nil {
		return errors.New("retrieving instance "+instanceId, err)
	}
	image, err := p.GetImage(instance.ImageId)
	if err != nil {
		return errors.New("retrieving image "+instance.ImageId, err)
	}
	vm, err := virtualboxclient.GetVm(instance.Id)
	if err != nil {
		return errors.New("retreiving vm from virtualbox", err)
	}
	var controllerKey string
	for _, device := range vm.Devices {
		if filepath.Clean(device.DiskFile) == filepath.Clean(getVolumePath(volume.Name)) {
			controllerKey = device.ControllerKey
		}
	}
	if controllerKey == "" {
		return errors.New("could not find device attached to "+instance.Name+" that matches volume "+getVolumePath(volume.Name), nil)
	}

	controllerPort, err := strconv.Atoi(controllerKey)
	if err != nil {
		return errors.New("could not convert "+controllerKey+" to int", err)
	}
	logrus.Debugf("using storage controller %s", image.RunSpec.StorageDriver)

	if err := virtualboxclient.DetachDisk(instance.Id, controllerPort, image.RunSpec.StorageDriver); err != nil {
		return errors.New("detaching disk from vm", err)
	}

	if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volume, ok := volumes[volume.Id]
		if !ok {
			return errors.New("no record of "+volume.Id+" in the state", nil)
		}
		volume.Attachment = ""
		return nil
	}); err != nil {
		return errors.New("modifying volume map in state", err)
	}
	return nil
}
