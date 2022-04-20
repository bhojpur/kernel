package firecracker

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
	"sync"

	firecrackersdk "github.com/firecracker-microvm/firecracker-go-sdk"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
)

type FirecrackerProvider struct {
	config config.Firecracker
	state  state.State

	runningMachines map[string]*firecrackersdk.Machine
	mapLock         sync.RWMutex
}

func FirecrackerStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "firecracker/state.json")

}
func firecrackerImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "firecracker/images/")
}

func firecrackerInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "firecracker/instances/")
}

func firecrackerVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "firecracker/volumes/")
}

func NewProvider(config config.Firecracker) (*FirecrackerProvider, error) {

	os.MkdirAll(firecrackerImagesDirectory(), 0777)
	os.MkdirAll(firecrackerInstancesDirectory(), 0777)
	os.MkdirAll(firecrackerVolumesDirectory(), 0777)

	p := &FirecrackerProvider{
		config:          config,
		state:           state.NewBasicState(FirecrackerStateFile()),
		runningMachines: map[string]*firecrackersdk.Machine{},
	}

	return p, nil
}

func (p *FirecrackerProvider) WithState(state state.State) *FirecrackerProvider {
	p.state = state
	return p
}

func getImagePath(imageName string) string {
	return filepath.Join(firecrackerImagesDirectory(), imageName, "boot.img")
}

func getVolumePath(volumeName string) string {
	return filepath.Join(firecrackerVolumesDirectory(), volumeName, "data.img")
}

func getInstanceDir(instanceName string) string {
	return filepath.Join(firecrackerInstancesDirectory(), instanceName)
}

func getImageDir(imageName string) string {
	return filepath.Join(firecrackerImagesDirectory(), imageName)
}
