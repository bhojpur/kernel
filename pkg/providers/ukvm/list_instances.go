package ukvm

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
	"os"
	"strconv"
	"syscall"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
)

func (p *UkvmProvider) ListInstances() ([]*types.Instance, error) {
	if len(p.state.GetInstances()) < 1 {
		return []*types.Instance{}, nil
	}

	var instances []*types.Instance
	for _, instance := range p.state.GetInstances() {
		pid, err := strconv.Atoi(instance.Id)
		if err != nil {
			return nil, errors.New("invalid id (is not a pid)", err)
		}
		if err := detectInstance(pid); err != nil {
			p.state.RemoveInstance(instance)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func detectInstance(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return errors.New("Failed to find process", err)
	}
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return errors.New(fmt.Sprintf("process.Signal on pid %d returned", pid), err)
	}
	return nil
}
