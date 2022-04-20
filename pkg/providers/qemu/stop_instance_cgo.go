//go:build cgo
// +build cgo

package qemu

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
	"strconv"
	"syscall"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *QemuProvider) StopInstance(id string) error {
	instance, err := p.GetInstance(id)
	if err != nil {
		return errors.New("retrieving instance "+id, err)
	}

	// kill qemu
	pid, err := strconv.Atoi(instance.Id)
	if err != nil {
		return errors.New("invalid instance id (should be qemu pid)", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logrus.Warn("failed finding instance, assuming instance has externally terminated", err)
	} else {
		if err := process.Signal(syscall.SIGKILL); err != nil {
			logrus.Warn("failed terminating instance, assuming instance has externally terminated", err)
		}
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

	return p.state.RemoveInstance(instance)
}
