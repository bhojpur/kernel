package openstack

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
	"time"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/sirupsen/logrus"
)

const DEFAULT_INSTANCE_DISKMB int = 10 * 1024 // 10 GB

func (p *OpenstackProvider) RunInstance(params types.RunInstanceParams) (_ *types.Instance, err error) {
	// return nil, errors.New("not yet supportded for openstack", nil)

	logrus.WithFields(logrus.Fields{
		"image-id": params.ImageId,
		"mounts":   params.MntPointsToVolumeIds,
		"env":      params.Env,
	}).Infof("running instance %s", params.Name)

	clientNova, err := p.newClientNova()
	if err != nil {
		return nil, err
	}

	image, err := p.GetImage(params.ImageId)
	if err != nil {
		return nil, errors.New("failed to get image", err)
	}

	// If not set, use default.
	if params.InstanceMemory <= 0 {
		params.InstanceMemory = image.RunSpec.DefaultInstanceMemory
	}

	// Pick flavor.
	minDiskMB := image.RunSpec.MinInstanceDiskMB
	if minDiskMB <= 0 {
		// TODO(miha-plesko): raise error here, since compiler should set MinInstanceDiskMB.
		// This commit adds field MinInstanceDiskMB to the RunSpec, but ATM non of the existing
		// compilers actually set it (so it's always zero). This field should be set at compile time
		// since only then compiler is actually aware of the logical size of the disk.
		// Raise error here after compiler is updated.
		minDiskMB = DEFAULT_INSTANCE_DISKMB
	}
	flavor, err := pickFlavor(clientNova, minDiskMB, params.InstanceMemory)
	if err != nil {
		return nil, errors.New("failed to pick flavor", err)
	}

	// Run instance.
	serverId, err := launchServer(clientNova, params.Name, flavor.Name, image.Name, p.config.NetworkUUID)
	if err != nil {
		return nil, errors.New("failed to run instance", err)
	}

	instance := &types.Instance{
		Id:             serverId,
		Name:           params.Name,
		State:          types.InstanceState_Pending,
		Infrastructure: types.Infrastructure_OPENSTACK,
		ImageId:        image.Id,
		Created:        time.Now(),
	}

	// Update state.
	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		instances[instance.Id] = instance
		return nil
	}); err != nil {
		return nil, errors.New("failed to modify instance map in state", err)
	}

	logrus.WithFields(logrus.Fields{"instance": instance}).Infof("instance created succesfully")

	return instance, nil
}

// launchServer launches single server of given image and returns it's id.
func launchServer(clientNova *gophercloud.ServiceClient, name, flavorName, imageName, networkUUID string) (string, error) {
	resp := servers.Create(clientNova, servers.CreateOpts{
		Name:       name,
		FlavorName: flavorName,
		ImageName:  imageName,
		Networks: []servers.Network{
			servers.Network{UUID: networkUUID},
		},
	})

	if resp.Err != nil {
		return "", resp.Err
	}

	server, err := resp.Extract()
	return server.ID, err
}
