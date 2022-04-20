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
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/providers/xen/xenclient"
	"github.com/bhojpur/kernel/pkg/state"
)

type XenProvider struct {
	state  state.State
	client *xenclient.XenClient
}

func XenStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "xen/state.json")

}
func xenImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "xen/images/")
}

func xenInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "xen/instances/")
}

func xenVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "xen/volumes/")
}

func NewXenProvider(config config.Xen) (*XenProvider, error) {

	os.MkdirAll(xenImagesDirectory(), 0777)
	os.MkdirAll(xenInstancesDirectory(), 0777)
	os.MkdirAll(xenVolumesDirectory(), 0777)

	p := &XenProvider{
		state: state.NewBasicState(XenStateFile()),
		client: &xenclient.XenClient{
			KernelPath: config.KernelPath,
			XenBridge:  config.XenBridge,
		},
	}

	/*
		if err := p.deployInstanceListener(); err != nil {
			return nil, errors.New("deploying xen instance listener", err)
		}
	*/

	return p, nil
}

func (p *XenProvider) WithState(state state.State) *XenProvider {
	p.state = state
	return p
}

func getImagePath(imageName string) string {
	return filepath.Join(xenImagesDirectory(), imageName, "boot.img")
}

func getInstanceDir(instanceName string) string {
	return filepath.Join(xenInstancesDirectory(), instanceName)
}

func getVolumePath(volumeName string) string {
	return filepath.Join(xenVolumesDirectory(), volumeName, "data.img")
}
