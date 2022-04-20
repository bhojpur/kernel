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
	"path/filepath"
	"strings"

	"time"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func VirtualboxStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "virtualbox/state.json")
}
func virtualboxImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "virtualbox/images/")
}
func virtualboxInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "virtualbox/instances/")
}
func virtualboxVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "virtualbox/volumes/")
}

const VboxKernelInstanceListener = "VboxKernelInstanceListener"
const instanceListenerPrefix = "bhojpur_virtualbox"

type VirtualboxProvider struct {
	config             config.Virtualbox
	state              state.State
	instanceListenerIp string
}

func NewVirtualboxProvider(config config.Virtualbox) (*VirtualboxProvider, error) {
	os.MkdirAll(virtualboxImagesDirectory(), 0755)
	os.MkdirAll(virtualboxInstancesDirectory(), 0755)
	os.MkdirAll(virtualboxVolumesDirectory(), 0755)

	p := &VirtualboxProvider{
		config: config,
		state:  state.NewBasicState(VirtualboxStateFile()),
	}

	if err := p.deployInstanceListener(config); err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, errors.New("deploying virtualbox instance listener", err)
	}

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout)
	if err != nil {
		return nil, errors.New("failed to retrieve instance listener ip. is Bhojpur Kernel instance listener running?", err)
	}

	p.instanceListenerIp = instanceListenerIp

	// begin update instances cycle
	go func() {
		for {
			if err := p.syncState(); err != nil {
				logrus.Error("error updatin virtualbox state:", err)
			}
			time.Sleep(time.Second)
		}
	}()

	return p, nil
}

func (p *VirtualboxProvider) WithState(state state.State) *VirtualboxProvider {
	p.state = state
	return p
}

func getImagePath(imageName string) string {
	return filepath.Join(virtualboxImagesDirectory(), imageName, "boot.vmdk")
}

func getInstanceDir(instanceName string) string {
	return filepath.Join(virtualboxInstancesDirectory(), instanceName)
}

func getVolumePath(volumeName string) string {
	return filepath.Join(virtualboxVolumesDirectory(), volumeName, "data.vmdk")
}
