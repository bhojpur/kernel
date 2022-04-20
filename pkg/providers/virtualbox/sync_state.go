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
	"os"
	"strings"
	"time"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/virtualbox/virtualboxclient"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VirtualboxProvider) syncState() error {
	if len(p.state.GetInstances()) < 1 {
		return nil
	}
	for _, instance := range p.state.GetInstances() {
		vm, err := virtualboxclient.GetVm(instance.Name)
		if err != nil {
			if strings.Contains(err.Error(), "Could not find a registered machine") {
				logrus.Warnf("instance found in state that is no longer registered to Virtualbox")
				os.RemoveAll(getInstanceDir(instance.Name))
				p.state.RemoveInstance(instance)
				continue
			}
			return errors.New("retrieving vm for instance id "+instance.Name, err)
		}
		macAddr := vm.MACAddr

		if vm.Running {
			instance.State = types.InstanceState_Running
		} else {
			instance.State = types.InstanceState_Stopped
		}

		var ipAddress string
		kutil.Retry(3, time.Duration(500*time.Millisecond), func() error {
			if instance.Name == VboxKernelInstanceListener {
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
