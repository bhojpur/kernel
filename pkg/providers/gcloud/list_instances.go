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
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *GcloudProvider) ListInstances() ([]*types.Instance, error) {
	if len(p.state.GetInstances()) < 1 {
		return []*types.Instance{}, nil
	}

	gInstances, err := p.compute().Instances.List(p.config.ProjectID, p.config.Zone).Do()
	if err != nil {
		return nil, errors.New("getting instance list from gcloud", err)
	}

	updatedInstances := []*types.Instance{}
	for _, instance := range p.state.GetInstances() {
		instanceFound := false
		//find instance in list
		for _, gInstance := range gInstances.Items {
			if gInstance.Name == instance.Name {
				instance.State = parseInstanceState(gInstance.Status)

				//use first network interface, skip if unavailable
				if len(gInstance.NetworkInterfaces) > 0 && len(gInstance.NetworkInterfaces[0].AccessConfigs) > 0 {
					instance.IpAddress = gInstance.NetworkInterfaces[0].AccessConfigs[0].NatIP
				}
				p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
					instances[instance.Id] = instance
					return nil
				})
				updatedInstances = append(updatedInstances, instance)
				instanceFound = true
				break
			}
		}
		if !instanceFound {
			logrus.Warnf("instance %v no longer found, cleaning it from state", instance.Name)
			p.state.RemoveInstance(instance)
		}
	}

	return updatedInstances, nil
}

func parseInstanceState(status string) types.InstanceState {
	switch status {
	case "RUNNING":
		return types.InstanceState_Running
	case "PROVISIONING":
		fallthrough
	case "STAGING":
		return types.InstanceState_Pending
	case "SUSPENDED":
		fallthrough
	case "STOPPING":
		fallthrough
	case "SUSPENDING":
		fallthrough
	case "STOPPED":
		return types.InstanceState_Stopped
	case "TERMINATED":
		return types.InstanceState_Terminated
	}
	return types.InstanceState_Unknown
}
