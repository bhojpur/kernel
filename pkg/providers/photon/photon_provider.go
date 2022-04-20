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
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/bhojpur/kernel/pkg/util/errors"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/vmware/photon-controller-go-sdk/photon"
)

type PhotonProvider struct {
	config    config.Photon
	state     state.State
	u         *url.URL
	client    *photon.Client
	projectId string
}

func PhotonStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "photon/state.json")

}
func photonImagesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "photon/images/")
}

func photonInstancesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "photon/instances/")
}

func photonVolumesDirectory() string {
	return filepath.Join(config.Internal.KernelHome, "photon/volumes/")
}

func NewPhotonProvider(config config.Photon) (*PhotonProvider, error) {

	os.MkdirAll(photonImagesDirectory(), 0755)
	os.MkdirAll(photonInstancesDirectory(), 0755)
	os.MkdirAll(photonVolumesDirectory(), 0755)

	p := &PhotonProvider{
		config: config,
		state:  state.NewBasicState(PhotonStateFile()),
	}

	p.client = photon.NewClient(p.config.PhotonURL, "", nil)
	p.projectId = p.config.ProjectId
	_, err := p.client.Status.Get()
	if err != nil {
		return nil, err
	}

	if err := p.DeployInstanceListener(config); err != nil /*&& !strings.Contains(err.Error(), "already exists")*/ {
		return nil, errors.New("deploying photon instance listener", err)
	}

	return p, nil
}

func (p *PhotonProvider) WithState(state state.State) *PhotonProvider {
	p.state = state
	return p
}

func (p *PhotonProvider) waitForTaskSuccess(task *photon.Task) (*photon.Task, error) {
	task, err := p.client.Tasks.WaitTimeout(task.ID, 30*time.Minute)
	if err != nil {
		return nil, errors.New("error waiting for task creating photon image", err)
	}

	if task.State != "COMPLETED" {
		return nil, errors.New("Error with task "+task.ID, nil)
	}

	return task, nil
}
