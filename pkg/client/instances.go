package client

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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bhojpur/kernel/pkg/daemon"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

type instances struct {
	kernelIP string
}

func (i *instances) All() ([]*types.Instance, error) {
	resp, body, err := lxhttpclient.Get(i.kernelIP, "/instances", nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var instances []*types.Instance
	if err := json.Unmarshal(body, &instances); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type []*types.Instance", string(body)), err)
	}
	return instances, nil
}

func (i *instances) Get(id string) (*types.Instance, error) {
	resp, body, err := lxhttpclient.Get(i.kernelIP, "/instances/"+id, nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var instance types.Instance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Instance", string(body)), err)
	}
	return &instance, nil
}

func (i *instances) Delete(id string, force bool) error {
	query := buildQuery(map[string]interface{}{
		"force": force,
	})
	resp, body, err := lxhttpclient.Delete(i.kernelIP, "/instances/"+id+query, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (i *instances) GetLogs(id string) (string, error) {
	resp, body, err := lxhttpclient.Get(i.kernelIP, "/instances/"+id+"/logs", nil)
	if err != nil {
		return "", errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return string(body), nil
}

func (i *instances) AttachLogs(id string, deleteOnDisconnect bool) (io.ReadCloser, error) {
	query := buildQuery(map[string]interface{}{
		"follow": true,
		"delete": deleteOnDisconnect,
	})
	resp, err := lxhttpclient.GetAsync(i.kernelIP, "/instances/"+id+"/logs"+query, nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v", resp.StatusCode), err)
	}
	return resp.Body, nil
}

func (i *instances) Run(instanceName, imageName string, mountPointsToVols, env map[string]string, memoryMb int, noCleanup, debugMode bool) (*types.Instance, error) {
	runInstanceRequest := daemon.RunInstanceRequest{
		InstanceName: instanceName,
		ImageName:    imageName,
		Mounts:       mountPointsToVols,
		Env:          env,
		MemoryMb:     memoryMb,
		NoCleanup:    noCleanup,
		DebugMode:    debugMode,
	}
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/instances/run", nil, runInstanceRequest)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	var instance types.Instance
	if err := json.Unmarshal(body, &instance); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Instance", string(body)), err)
	}
	return &instance, nil
}

func (i *instances) Start(id string) error {
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/instances/"+id+"/start", nil, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (i *instances) Stop(id string) error {
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/instances/"+id+"/stop", nil, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}
