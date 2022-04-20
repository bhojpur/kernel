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
	"strings"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/sirupsen/logrus"
)

func (p *OpenstackProvider) ListInstances() ([]*types.Instance, error) {
	// Return immediately if no instance is managed by Bhojpur Kernel.
	managedInstances := p.state.GetInstances()
	if len(managedInstances) < 1 {
		return []*types.Instance{}, nil
	}

	clientNova, err := p.newClientNova()
	if err != nil {
		return nil, err
	}

	instList, err := fetchInstances(clientNova, managedInstances)

	// Update state.
	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		// Clear everything.
		for k := range instances {
			delete(instances, k)
		}

		// Add fetched instances.
		for _, inst := range instList {
			instances[inst.Id] = inst
		}
		return nil
	}); err != nil {
		return nil, errors.New("failed to modify instance map in state", err)
	}

	return instList, nil
}

// fetchInstances fetches a list of instances runnign on OpenStack and returns a list of
// those that are managed by Bhojpur Kernel.
func fetchInstances(clientNova *gophercloud.ServiceClient, managedInstances map[string]*types.Instance) ([]*types.Instance, error) {
	var result []*types.Instance = make([]*types.Instance, 0)

	pagerServers := servers.List(clientNova, servers.ListOpts{})
	pagerServers.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		for _, s := range serverList {
			// Filter out instances that Bhojpur Kernel is not aware of.
			instance, ok := managedInstances[s.ID]
			if !ok {
				continue
			}

			// Interpret instance state and filter out instance with bad state.
			if state := parseInstanceState(s.Status); state == types.InstanceState_Terminated {
				continue
			} else {
				instance.State = state
			}

			// Update fields.
			instance.Name = s.Name
			instance.IpAddress = s.AccessIPv4

			result = append(result, instance)
		}

		return true, nil
	})
	return result, nil
}

func parseInstanceState(serverState string) types.InstanceState {
	// http://docs.openstack.org/developer/nova/vmstates.html#vm-states-and-possible-commands
	switch strings.ToLower(serverState) {
	case "active", "rescued":
		return types.InstanceState_Running
	case "building":
		return types.InstanceState_Pending
	case "paused":
		return types.InstanceState_Paused
	case "suspended":
		return types.InstanceState_Suspended
	case "shutoff", "stopped", "soft_deleted":
		return types.InstanceState_Stopped
	case "hard_deleted":
		return types.InstanceState_Terminated
	case "error":
		return types.InstanceState_Error
	}

	logrus.WithFields(logrus.Fields{
		"serverState": serverState,
	}).Infof("Received unknown instance state")

	return types.InstanceState_Unknown
}
