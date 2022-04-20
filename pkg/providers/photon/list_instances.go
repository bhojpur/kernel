package photon

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
	"github.com/bhojpur/kernel/pkg/util/errors"

	"github.com/bhojpur/kernel/pkg/types"
)

func (p *PhotonProvider) ListInstances() ([]*types.Instance, error) {
	if len(p.state.GetInstances()) < 1 {
		return []*types.Instance{}, nil
	}

	var instances []*types.Instance
	for _, instance := range p.state.GetInstances() {

		vm, err := p.client.VMs.Get(instance.Id)
		if err != nil {
			return nil, errors.New("retrieving vm for instance id "+instance.Id, err)
		}

		// TODO: get ip..

		switch vm.State {
		case "STARTED":
			instance.State = types.InstanceState_Running
		case "CREATING":
			instance.State = types.InstanceState_Pending
		case "STOPPED":
			fallthrough
		case "SUSPENDED":
			fallthrough
		default:
			instance.State = types.InstanceState_Stopped
			break
		}
		err = p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
			instances[instance.Id] = instance
			return nil
		})
		if err != nil {
			return nil, errors.New("saving instance to state", err)
		}

		instances = append(instances, instance)
	}

	return instances, nil
}
