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
	"fmt"
	"strings"
	"time"

	"os"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/xen/xenclient"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *XenProvider) RunInstance(params types.RunInstanceParams) (_ *types.Instance, err error) {
	logrus.WithFields(logrus.Fields{
		"image-id": params.ImageId,
		"mounts":   params.MntPointsToVolumeIds,
		"env":      params.Env,
	}).Infof("running instance %s", params.Name)

	if _, err := p.GetInstance(params.Name); err == nil {
		return nil, errors.New("instance with name "+params.Name+" already exists. xen provider requires unique names for instances", nil)
	}

	image, err := p.GetImage(params.ImageId)
	if err != nil {
		return nil, errors.New("getting image", err)
	}

	if err := common.VerifyMntsInput(p, image, params.MntPointsToVolumeIds); err != nil {
		return nil, errors.New("invalid mapping for volume", err)
	}

	volumeIdToDevice := make(map[string]string)

	// till we support pv without boot device, we need a boot device..
	bootmapping := "sda1"
	for _, mapping := range image.RunSpec.DeviceMappings {
		if mapping.MountPoint == "/" {
			bootmapping = removeDevFromDeviceName(mapping.DeviceName)
			break
		}
	}

	for mntPoint, volumeId := range params.MntPointsToVolumeIds {
		for _, mapping := range image.RunSpec.DeviceMappings {
			if mntPoint == mapping.MountPoint {
				volumeIdToDevice[volumeId] = mapping.DeviceName
				break
			}
		}
	}

	logrus.Debugf("creating xen vm")

	// TODO add support for boot drive mapping.

	var dataVolumes []xenclient.VolumeConfig
	for volid, deviceName := range volumeIdToDevice {
		volPath, err := p.getVolPath(volid)
		if err != nil {
			return nil, errors.New("failed to get volume path", err)
		}
		dataVolumes = append(dataVolumes, xenclient.VolumeConfig{
			ImagePath:  volPath,
			DeviceName: removeDevFromDeviceName(deviceName),
		})
	}

	if err := os.MkdirAll(getInstanceDir(params.Name), 0755); err != nil {
		return nil, errors.New("failed to create instance dir", err)
	}

	//if not set, use default
	if params.InstanceMemory <= 0 {
		params.InstanceMemory = image.RunSpec.DefaultInstanceMemory
	}

	xenParams := xenclient.CreateVmParams{
		Name:           params.Name,
		Memory:         params.InstanceMemory,
		BootImage:      getImagePath(image.Name),
		BootDeviceName: bootmapping,
		VmDir:          getInstanceDir(params.Name),
		DataVolumes:    dataVolumes,
	}

	if err := p.client.CreateVm(xenParams); err != nil {
		return nil, errors.New("creating xen domain", err)
	}

	instanceId := params.Name
	if doms, err := p.client.ListVms(); err == nil {
		for _, d := range doms {
			if d.Config.CInfo.Name == params.Name {
				instanceId = fmt.Sprintf("%d", d.Domid)
				break
			}
		}
	}

	var instanceIp string

	instance := &types.Instance{
		Id:             instanceId,
		Name:           params.Name,
		State:          types.InstanceState_Running,
		IpAddress:      instanceIp,
		Infrastructure: types.Infrastructure_XEN,
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

func (p *XenProvider) getVolPath(volId string) (string, error) {

	v, err := p.GetVolume(volId)
	if err != nil {
		return "", err
	}
	return getVolumePath(v.Name), nil

}

func removeDevFromDeviceName(devName string) string {

	const prefix = "/dev/"

	if strings.HasPrefix(devName, prefix) {
		devName = devName[len(prefix):]
	}

	return devName
}
