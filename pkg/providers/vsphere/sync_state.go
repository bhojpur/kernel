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
	"time"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/vsphere/vsphereclient"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VsphereProvider) syncState() error {
	if len(p.state.GetInstances()) < 1 {
		return nil
	}
	c := p.getClient()
	vms := []*vsphereclient.VirtualMachine{}
	for instanceId := range p.state.GetInstances() {
		vm, err := c.GetVmByUuid(instanceId)
		if err != nil {
			return errors.New("getting vm info for "+instanceId, err)
		}
		vms = append(vms, vm)
	}
	for _, vm := range vms {
		//we use mac address as the vm id
		macAddr := ""
		for _, device := range vm.Config.Hardware.Device {
			if len(device.MacAddress) > 0 {
				macAddr = device.MacAddress
				break
			}
		}
		if macAddr == "" {
			logrus.WithFields(logrus.Fields{"vm": vm}).Warnf("vm found, cannot identify mac addr")
			continue
		}

		instanceId := vm.Config.UUID
		instance, ok := p.state.GetInstances()[instanceId]
		if !ok {
			continue
		}

		switch vm.Summary.Runtime.PowerState {
		case "poweredOn":
			instance.State = types.InstanceState_Running
			break
		case "poweredOff":
		case "suspended":
			instance.State = types.InstanceState_Stopped
			break
		default:
			instance.State = types.InstanceState_Unknown
			break
		}

		var ipAddress string
		kutil.Retry(3, time.Duration(500*time.Millisecond), func() error {
			if instance.Name == VsphereKernelInstanceListener {
				ipAddress = p.instanceListenerIp
			} else {
				var err error
				ipAddress, err = common.GetInstanceIp(p.instanceListenerIp, 3000, macAddr)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
			if _, ok := instances[instance.Id]; ok {
				instances[instance.Id].IpAddress = ipAddress
				instances[instance.Id].State = instance.State
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
