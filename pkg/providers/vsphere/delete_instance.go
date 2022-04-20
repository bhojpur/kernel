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
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VsphereProvider) DeleteInstance(id string, force bool) error {
	instance, err := p.GetInstance(id)
	if err != nil {
		return errors.New("retrieving instance "+id, err)
	}
	if instance.State == types.InstanceState_Running {
		if force {
			if err := p.StopInstance(instance.Id); err != nil {
				return errors.New("stopping instance for deletion", err)
			}
		} else {
			return errors.New("instance "+instance.Id+"is still running. try again with --force or power off instance first", err)
		}
	}
	image, err := p.GetImage(instance.ImageId)
	if err != nil {
		return errors.New("getting image for instance", err)
	}
	volumesToDetach := []*types.Volume{}
	volumes, err := p.ListVolumes()
	if err != nil {
		return errors.New("getting volume list", err)
	}
	for _, volume := range volumes {
		if volume.Attachment == instance.Id {
			logrus.Debugf("detaching volume: %v", volume)
			volumesToDetach = append(volumesToDetach, volume)
		}
	}

	c := p.getClient()
	for controllerPort, deviceMapping := range image.RunSpec.DeviceMappings {
		if deviceMapping.MountPoint != "/" {
			if err := c.DetachDisk(instance.Id, controllerPort, image.RunSpec.StorageDriver); err != nil {
				return errors.New("detaching volume from instance", err)
			}
		}
	}
	err = c.DestroyVm(instance.Name)
	if err != nil {
		return errors.New("failed to terminate instance "+instance.Id, err)
	}
	return p.state.RemoveInstance(instance)
}
