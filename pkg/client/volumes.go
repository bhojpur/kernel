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
	"net/http"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

type volumes struct {
	kernelIP string
}

func (v *volumes) All() ([]*types.Volume, error) {
	resp, body, err := lxhttpclient.Get(v.kernelIP, "/volumes", nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var volumes []*types.Volume
	if err := json.Unmarshal(body, &volumes); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type []*types.Volume", string(body)), err)
	}
	return volumes, nil
}

func (v *volumes) Get(id string) (*types.Volume, error) {
	resp, body, err := lxhttpclient.Get(v.kernelIP, "/volumes/"+id, nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var volume types.Volume
	if err := json.Unmarshal(body, &volume); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Volume", string(body)), err)
	}
	return &volume, nil
}

func (v *volumes) Delete(id string, force bool) error {
	query := buildQuery(map[string]interface{}{
		"force": force,
	})
	resp, body, err := lxhttpclient.Delete(v.kernelIP, "/volumes/"+id+query, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (v *volumes) Create(name, dataTar, provider string, raw bool, size int, volType string, noCleanup bool) (*types.Volume, error) {
	query := buildQuery(map[string]interface{}{
		"size":       size,
		"provider":   provider,
		"type":       volType,
		"no_cleanup": noCleanup,
		"raw":        raw,
	})
	//no data provided
	var (
		resp *http.Response
		body []byte
		err  error
	)
	if dataTar == "" {
		resp, body, err = lxhttpclient.Post(v.kernelIP, "/volumes/"+name+query, nil, nil)
		if err != nil {
			return nil, errors.New("request failed", err)
		}
		if resp.StatusCode != http.StatusCreated {
			return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
		}
	} else {
		resp, body, err = lxhttpclient.PostFile(v.kernelIP, "/volumes/"+name+query, "tarfile", dataTar)
		if err != nil {
			return nil, errors.New("request failed", err)
		}
		if resp.StatusCode != http.StatusCreated {
			return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
		}
	}
	var volume types.Volume
	if err := json.Unmarshal(body, &volume); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Volume", string(body)), err)
	}
	return &volume, nil
}

func (v *volumes) Attach(id, instanceId, mountPoint string) error {
	query := buildQuery(map[string]interface{}{
		"mount": mountPoint,
	})
	resp, body, err := lxhttpclient.Post(v.kernelIP, "/volumes/"+id+"/attach/"+instanceId+query, nil, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (v *volumes) Detach(id string) error {
	resp, body, err := lxhttpclient.Post(v.kernelIP, "/volumes/"+id+"/detach", nil, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}
