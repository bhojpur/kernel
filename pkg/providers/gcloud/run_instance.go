package gcloud

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
	"encoding/json"
	"fmt"
	"time"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

func (p *GcloudProvider) RunInstance(params types.RunInstanceParams) (_ *types.Instance, err error) {
	logrus.WithFields(logrus.Fields{
		"image-id": params.ImageId,
		"mounts":   params.MntPointsToVolumeIds,
		"env":      params.Env,
	}).Infof("running instance %s", params.Name)

	var instanceId string

	defer func() {
		if err != nil {
			logrus.WithError(err).Errorf("gcloud running instance encountered an error")
			if instanceId != "" {
				if params.NoCleanup {
					logrus.Warnf("because --no-cleanup flag was provided, not cleaning up failed instance %s0", instanceId)
					return
				}
				logrus.Warnf("cleaning up instance %s", instanceId)
				p.compute().Instances.Delete(p.config.ProjectID, p.config.Zone, instanceId)
				if cleanupErr := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
					delete(instances, instanceId)
					return nil
				}); cleanupErr != nil {
					logrus.Error(errors.New("modifying instance map in state", cleanupErr))
				}
			}
		}
	}()

	image, err := p.GetImage(params.ImageId)
	if err != nil {
		return nil, errors.New("getting image", err)
	}

	if err := common.VerifyMntsInput(p, image, params.MntPointsToVolumeIds); err != nil {
		return nil, errors.New("invalid mapping for volume", err)
	}

	envData, err := json.Marshal(params.Env)
	if err != nil {
		return nil, errors.New("could not convert instance env to json", err)
	}

	//if not set, use default
	if params.InstanceMemory <= 0 {
		params.InstanceMemory = image.RunSpec.DefaultInstanceMemory
	}

	if len(envData) > 32768 {
		return nil, errors.New("total length of env metadata must be <= 32768 bytes; have json string "+string(envData), nil)
	}

	disks := []*compute.AttachedDisk{
		//boot disk
		&compute.AttachedDisk{
			AutoDelete: true,
			Boot:       true,
			//DeviceName: "sd0"
			InitializeParams: &compute.AttachedDiskInitializeParams{
				SourceImage: "global/images/" + image.Name,
			},
		},
	}

	for _, volumeId := range params.MntPointsToVolumeIds {
		disks = append(disks, &compute.AttachedDisk{
			AutoDelete: false,
			Boot:       false,
			Source:     volumeId,
		})
	}

	instanceSpec := &compute.Instance{
		Name: params.Name,
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				&compute.MetadataItems{
					Key:   "ENV_DATA",
					Value: pointerTo(string(envData)),
				},
			},
		},
		Disks:       disks,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", p.config.Zone, "g1-small"),
		NetworkInterfaces: []*compute.NetworkInterface{
			&compute.NetworkInterface{
				AccessConfigs: []*compute.AccessConfig{
					&compute.AccessConfig{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: "global/networks/default",
			},
		},
	}

	gInstance, err := p.compute().Instances.Insert(p.config.ProjectID, p.config.Zone, instanceSpec).Do()
	if err != nil {
		return nil, errors.New("creating instance on gcloud failed", err)
	}
	logrus.Infof("gcloud instance created: %+v", gInstance)

	instanceId = params.Name

	//must add instance to state before attaching volumes
	instance := &types.Instance{
		Id:             instanceId,
		Name:           params.Name,
		State:          types.InstanceState_Pending,
		Infrastructure: types.Infrastructure_GCLOUD,
		ImageId:        image.Id,
		Created:        time.Now(),
	}

	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		instances[instance.Id] = instance
		return nil
	}); err != nil {
		return nil, errors.New("modifying instance map in state", err)
	}

	logrus.WithFields(logrus.Fields{"instance": instance}).Infof("instance created succesfully")

	return instance, nil
}

func pointerTo(v string) *string {
	return &v
}
