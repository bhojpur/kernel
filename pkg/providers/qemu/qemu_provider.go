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
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/bhojpur/kernel/pkg/util/errors"
)

var debuggerTargetImageName string

type QemuProvider struct {
	config config.Qemu
	state  state.State
}

func QemuStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "qemu/state.json")

}
func qemuImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "qemu/images/")
}

func qemuInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "qemu/instances/")
}

func qemuVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "qemu/volumes/")
}

func NewQemuProvider(config config.Qemu) (*QemuProvider, error) {

	os.MkdirAll(qemuImagesDirectory(), 0777)
	os.MkdirAll(qemuInstancesDirectory(), 0777)
	os.MkdirAll(qemuVolumesDirectory(), 0777)

	if config.DebuggerPort == 0 {
		config.DebuggerPort = 3001
	}

	if err := startDebuggerListener(config.DebuggerPort); err != nil {
		return nil, errors.New("establishing debugger tcp listener", err)
	}

	p := &QemuProvider{
		config: config,
		state:  state.NewBasicState(QemuStateFile()),
	}

	return p, nil
}

func (p *QemuProvider) WithState(state state.State) *QemuProvider {
	p.state = state
	return p
}

func getImagePath(imageName string) string {
	return filepath.Join(qemuImagesDirectory(), imageName, "boot.img")
}

func getKernelPath(imageName string) string {
	return filepath.Join(qemuImagesDirectory(), imageName, "program.bin")
}

func getCmdlinePath(imageName string) string {
	return filepath.Join(qemuImagesDirectory(), imageName, "cmdline")
}

func getVolumePath(volumeName string) string {
	return filepath.Join(qemuVolumesDirectory(), volumeName, "data.img")
}
