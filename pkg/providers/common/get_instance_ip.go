package common

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

	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/layer-x/layerx-commons/lxhttpclient"
)

func GetInstanceIp(listenerIp string, listenerPort int, instanceId string) (string, error) {
	_, body, err := lxhttpclient.Get(fmt.Sprintf("%s:%v", listenerIp, listenerPort), "/instances", nil)
	if err != nil {
		return "", errors.New("http GET on instance listener", err)
	}
	var instanceIpMap map[string]string
	if err := json.Unmarshal(body, &instanceIpMap); err != nil {
		return "", errors.New("unmarshalling response ("+string(body)+") to map", err)
	}
	ip, ok := instanceIpMap[instanceId]
	if !ok {
		return "", errors.New("instance "+instanceId+" not found in map: "+fmt.Sprintf("%v", instanceIpMap), err)
	}
	return ip, nil
}
