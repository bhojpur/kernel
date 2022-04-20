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
	"strings"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

type images struct {
	kernelIP string
}

func (i *images) All() ([]*types.Image, error) {
	resp, body, err := lxhttpclient.Get(i.kernelIP, "/images", nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var images []*types.Image
	if err := json.Unmarshal(body, &images); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type []*types.Image", string(body)), err)
	}
	return images, nil
}

func (i *images) Get(id string) (*types.Image, error) {
	resp, body, err := lxhttpclient.Get(i.kernelIP, "/images/"+id, nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var image types.Image
	if err := json.Unmarshal(body, &image); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return &image, nil
}

func (i *images) Build(name, sourceTar, base, lang, provider, args string, mounts []string, force, noCleanup bool) (*types.Image, error) {
	query := buildQuery(map[string]interface{}{
		"base":       base,
		"lang":       lang,
		"provider":   provider,
		"args":       args,
		"mounts":     strings.Join(mounts, ","),
		"force":      force,
		"no_cleanup": noCleanup,
	})
	resp, body, err := lxhttpclient.PostFile(i.kernelIP, "/images/"+name+"/create"+query, "tarfile", sourceTar)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	var image types.Image
	if err := json.Unmarshal(body, &image); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return &image, nil
}

func (i *images) Delete(id string, force bool) error {
	query := buildQuery(map[string]interface{}{
		"force": force,
	})
	resp, body, err := lxhttpclient.Delete(i.kernelIP, "/images/"+id+query, nil)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (i *images) Push(c config.HubConfig, imageName string) error {
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/images/push/"+imageName, nil, c)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (i *images) Pull(c config.HubConfig, imageName, provider string, force bool) error {
	query := buildQuery(map[string]interface{}{
		"provider": provider,
		"force":    force,
	})
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/images/pull/"+imageName+query, nil, c)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (i *images) RemoteDelete(c config.HubConfig, imageName string) error {
	resp, body, err := lxhttpclient.Post(i.kernelIP, "/images/remote-delete/"+imageName, nil, c)
	if err != nil {
		return errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}
