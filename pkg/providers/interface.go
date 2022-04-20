package providers

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
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
)

type Provider interface {
	GetConfig() ProviderConfig
	//Images
	Stage(params types.StageImageParams) (*types.Image, error)
	ListImages() ([]*types.Image, error)
	GetImage(nameOrIdPrefix string) (*types.Image, error)
	DeleteImage(id string, force bool) error
	//Instances
	RunInstance(params types.RunInstanceParams) (*types.Instance, error)
	ListInstances() ([]*types.Instance, error)
	GetInstance(nameOrIdPrefix string) (*types.Instance, error)
	DeleteInstance(id string, force bool) error
	StartInstance(id string) error
	StopInstance(id string) error
	GetInstanceLogs(id string) (string, error)
	//Volumes
	CreateVolume(params types.CreateVolumeParams) (*types.Volume, error)
	ListVolumes() ([]*types.Volume, error)
	GetVolume(nameOrIdPrefix string) (*types.Volume, error)
	DeleteVolume(id string, force bool) error
	AttachVolume(id, instanceId, mntPoint string) error
	DetachVolume(id string) error
	//Hub
	PullImage(params types.PullImagePararms) error
	PushImage(params types.PushImagePararms) error
	RemoteDeleteImage(params types.RemoteDeleteImagePararms) error
}

type ProviderConfig struct {
	UsePartitionTables bool
}

type Providers map[string]Provider

func (providers Providers) Keys() []string {
	keys := []string{}
	for providerType := range providers {
		keys = append(keys, providerType)
	}
	return keys
}

func (providers Providers) ProviderForImage(imageId string) (Provider, error) {
	for _, provider := range providers {
		_, err := provider.GetImage(imageId)
		if err == nil {
			return provider, nil
		}
	}
	return nil, errors.New("image "+imageId+" not found", nil)
}

func (providers Providers) ProviderForInstance(instanceId string) (Provider, error) {
	for _, provider := range providers {
		_, err := provider.GetInstance(instanceId)
		if err == nil {
			return provider, nil
		}
	}
	return nil, errors.New("instance "+instanceId+" not found", nil)
}

func (providers Providers) ProviderForVolume(volumeId string) (Provider, error) {
	for _, provider := range providers {
		_, err := provider.GetVolume(volumeId)
		if err == nil {
			return provider, nil
		}
	}
	return nil, errors.New("volume "+volumeId+" not found", nil)
}
