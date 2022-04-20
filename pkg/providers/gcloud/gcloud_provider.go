package gcloud

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
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"
)

func GcloudStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "gcloud", "state.json")
}

type GcloudProvider struct {
	config     config.Gcloud
	state      state.State
	computeSvc *compute.Service
	storageSvc *storage.Service
}

func NewGcloudProvier(config config.Gcloud) (*GcloudProvider, error) {
	logrus.Infof("state file: %s", GcloudStateFile())

	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, errors.New("failed to start default client", err)
	}
	computeService, err := compute.New(client)
	if err != nil {
		return nil, errors.New("failed to start compute client", err)
	}

	storageSevice, err := storage.New(client)
	if err != nil {
		return nil, errors.New("failed to start storage client", err)
	}

	return &GcloudProvider{
		config:     config,
		state:      state.NewBasicState(GcloudStateFile()),
		computeSvc: computeService,
		storageSvc: storageSevice,
	}, nil
}

func (p *GcloudProvider) WithState(state state.State) *GcloudProvider {
	p.state = state
	return p
}

func (p *GcloudProvider) compute() *compute.Service {
	return p.computeSvc
}

func (p *GcloudProvider) storage() *storage.Service {
	return p.storageSvc
}
