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
	"net/url"
	"strings"

	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

type client struct {
	kernelIP string
}

func KernelClient(kernelIP string) *client {
	return &client{kernelIP: kernelIP}
}

func (c *client) Images() *images {
	return &images{kernelIP: c.kernelIP}
}

func (c *client) Instances() *instances {
	return &instances{kernelIP: c.kernelIP}
}

func (c *client) Volumes() *volumes {
	return &volumes{kernelIP: c.kernelIP}
}

func (c *client) AvailableCompilers() ([]string, error) {
	resp, body, err := lxhttpclient.Get(c.kernelIP, "/available_compilers", nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var compilers []string
	if err := json.Unmarshal(body, &compilers); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return compilers, nil
}

func (c *client) AvailableProviders() ([]string, error) {
	resp, body, err := lxhttpclient.Get(c.kernelIP, "/available_providers", nil)
	if err != nil {
		return nil, errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var compilers []string
	if err := json.Unmarshal(body, &compilers); err != nil {
		return nil, errors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Image", string(body)), err)
	}
	return compilers, nil
}

func (c *client) DescribeCompiler(base string, lang string, provider string) (string, error) {
	query := buildQuery(map[string]interface{}{
		"base":     base,
		"lang":     lang,
		"provider": provider,
	})
	resp, body, err := lxhttpclient.Get(c.kernelIP, "/describe_compiler"+query, nil)
	if err != nil {
		return "", errors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return string(body), nil
}

func buildQuery(params map[string]interface{}) string {
	queryArray := []string{}
	for key, val := range params {
		queryArray = append(queryArray, url.QueryEscape(fmt.Sprintf("%s", key))+"="+url.QueryEscape(fmt.Sprintf("%v", val)))
	}
	queryString := "?" + strings.Join(queryArray, "&")
	return queryString
}
