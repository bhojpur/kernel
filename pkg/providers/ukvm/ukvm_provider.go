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
	"os"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
)

type UkvmProvider struct {
	config config.Ukvm
	state  state.State
}

func UkvmStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "ukvm/state.json")

}
func ukvmImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "ukvm/images/")
}

func ukvmInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "ukvm/instances/")
}

func ukvmVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "ukvm/volumes/")
}

func NewUkvmProvider(config config.Ukvm) (*UkvmProvider, error) {

	os.MkdirAll(ukvmImagesDirectory(), 0777)
	os.MkdirAll(ukvmInstancesDirectory(), 0777)
	os.MkdirAll(ukvmVolumesDirectory(), 0777)

	p := &UkvmProvider{
		config: config,
		state:  state.NewBasicState(UkvmStateFile()),
	}

	return p, nil
}

func (p *UkvmProvider) WithState(state state.State) *UkvmProvider {
	p.state = state
	return p
}
func getImageDir(imageName string) string {
	return filepath.Join(ukvmImagesDirectory(), imageName)
}
func getKernelPath(imageName string) string {
	return filepath.Join(ukvmImagesDirectory(), imageName, "program.bin")
}
func getUkvmPath(imageName string) string {
	return filepath.Join(ukvmImagesDirectory(), imageName, "ukvm-bin")
}

func getInstanceDir(instanceName string) string {
	return filepath.Join(ukvmInstancesDirectory(), instanceName)
}

func getInstanceLogName(instanceName string) string {
	return filepath.Join(ukvmInstancesDirectory(), instanceName, "stdout")
}

func getVolumePath(volumeName string) string {
	return filepath.Join(ukvmVolumesDirectory(), volumeName, "data.img")
}
