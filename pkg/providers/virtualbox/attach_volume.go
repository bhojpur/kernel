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
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/virtualbox/virtualboxclient"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VirtualboxProvider) AttachVolume(id, instanceId, mntPoint string) error {
	volume, err := p.GetVolume(id)
	if err != nil {
		return errors.New("retrieving volume "+id, err)
	}
	instance, err := p.GetInstance(instanceId)
	if err != nil {
		return errors.New("retrieving instance "+instanceId, err)
	}
	image, err := p.GetImage(instance.ImageId)
	if err != nil {
		return errors.New("retrieving image for instance", err)
	}
	controllerPort, err := common.GetControllerPortForMnt(image, mntPoint)
	if err != nil {
		return errors.New("getting controller port for mnt point", err)
	}
	logrus.Debugf("using storage controller %s", image.RunSpec.StorageDriver)

	if err := virtualboxclient.RefreshDiskUUID(getVolumePath(volume.Name)); err != nil {
		return errors.New("refreshing disk uuid", err)
	}

	if err := virtualboxclient.AttachDisk(instance.Id, getVolumePath(volume.Name), controllerPort, image.RunSpec.StorageDriver); err != nil {
		return errors.New("attaching disk to vm", err)
	}

	if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volume, ok := volumes[volume.Id]
		if !ok {
			return errors.New("no record of "+volume.Id+" in the state", nil)
		}
		volume.Attachment = instance.Id
		return nil
	}); err != nil {
		return errors.New("modifying volumes in state", err)
	}
	return nil
}
