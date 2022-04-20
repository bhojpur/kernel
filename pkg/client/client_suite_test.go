package client_test

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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bhojpur/kernel/pkg/daemon"
	"github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/tests/helpers"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/sirupsen/logrus"
)

var cfg = helpers.NewTestConfig()

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	var d *daemon.KernelDaemon
	var tmpKernel helpers.TempKernelHome
	BeforeSuite(func() {
		if os.Getenv("DEBUG_OFF") != "1" {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if os.Getenv("MAKE_CONTAINERS") == "1" {
			if err := helpers.MakeContainers(helpers.GetProjectRoot()); err != nil {
				logrus.Panic(err)
			}
		}

		tmpKernel.SetupKernel()
		var err error
		d, err = daemon.NewKernelDaemon(cfg)
		if err != nil {
			logrus.Panic(err)
		}
		go d.Run(3000)

	})
	AfterSuite(func() {
		//if err := helpers.RemoveContainers(projectRoot); err != nil {
		//	logrus.Panic(err)
		//}
		defer tmpKernel.TearDownKernel()
		err := d.Stop()
		if err != nil {
			logrus.Panic(err)
		}
	})
	RunSpecs(t, "Client Suite")
}

func testInstancePing(instanceIp string) {
	testInstanceEndpoint(instanceIp, "/ping_test", "pong")
}

func testInstanceEnv(instanceIp string) {
	testInstanceEndpoint(instanceIp, "/env_test", "VAL")
}

func testInstanceMount(instanceIp string) {
	testInstanceEndpoint(instanceIp, "/mount_test", "test_data")
}

func testInstanceEndpoint(instanceIp, path, expectedResponse string) {
	var resp *http.Response
	var body []byte
	var err error
	err = util.Retry(10, 2*time.Second, func() error {
		resp, body, err = lxhttpclient.Get(instanceIp+":8080", path, nil)
		return err
	})
	logrus.WithFields(logrus.Fields{
		"resp": resp,
		"body": string(body),
		"err":  err,
	}).Debugf("got resp")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	var testResponse struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &testResponse)
	Expect(err).ToNot(HaveOccurred())
	Expect(testResponse.Message).To(ContainSubstring(expectedResponse))
}
