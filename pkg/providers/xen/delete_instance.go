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
	"os"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
)

func (p *XenProvider) DeleteInstance(id string, force bool) error {
	instance, err := p.GetInstance(id)
	if err != nil {
		return errors.New("retrieving instance "+id, err)
	}

	if err := p.client.DestroyVm(instance.Name); err != nil {
		return errors.New("xen api destroy call", err)
	}

	if err := os.RemoveAll(getInstanceDir(instance.Name)); err != nil {
		return errors.New("deleting instance dir", err)
	}

	volumesToDetach := []*types.Volume{}
	volumes, err := p.ListVolumes()
	if err != nil {
		return errors.New("getting volume list", err)
	}
	for _, volume := range volumes {
		if volume.Attachment == instance.Id {
			volumesToDetach = append(volumesToDetach, volume)
		}
	}
	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		delete(instances, instance.Id)
		return nil
	}); err != nil {
		return errors.New("modifying image map in state", err)
	}
	return p.state.RemoveInstance(instance)
}
