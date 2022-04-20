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
	. "github.com/bhojpur/kernel/pkg/client"

	"fmt"
	"strings"
	"time"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/tests/helpers"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

const (
	test_go_app        = "test_go_app"
	test_rump_java_app = "test_rump_java_app"
	test_java_app      = "test_java_app"
	test_jar_app       = "test_jar_app"
	test_nodejs_app    = "test_nodejs_app"
	test_python3_app   = "test_python3_app"
)

var ipTimeout = time.Second * 180

var _ = Describe("Bhojpur Kernel Functionality", func() {
	daemonUrl := "127.0.0.1:3000"
	var c = KernelClient(daemonUrl)
	Describe("instances", func() {
		Describe("All()", func() {
			var image *types.Image
			var volume *types.Volume
			AfterEach(func() {
				if image != nil {
					c.Images().Delete(image.Id, true)
				}
				if volume != nil {
					c.Volumes().Delete(volume.Id, true)
				}
			})
			Context("instances exist", func() {
				Describe("Run()", func() {
					imagesWithVolumes := []string{
						test_go_app,
						test_python3_app,
						test_nodejs_app,
					}
					imagesWithoutVolumes := []string{
						test_java_app,
						test_jar_app,
					}
					providersWithVolumes := []string{}
					//TODO: aws should support mounts
					providersWithoutVolumes := []string{}
					if len(cfg.Providers.Virtualbox) > 0 {
						providersWithVolumes = append(providersWithVolumes, "virtualbox")
					}
					if len(cfg.Providers.Vsphere) > 0 {
						providersWithVolumes = append(providersWithVolumes, "vsphere")
					}
					if len(cfg.Providers.Aws) > 0 {
						providersWithoutVolumes = append(providersWithoutVolumes, "aws")
					}
					if len(cfg.Providers.Xen) > 0 {
						providersWithoutVolumes = append(providersWithoutVolumes, "xen")
					}
					entries := []table.TableEntry{}
					for _, imageName := range imagesWithVolumes {
						for _, provider := range providersWithVolumes {
							entries = append(entries, table.Entry(imageName+" on "+provider, imageName, true, provider))
						}
						for _, provider := range providersWithoutVolumes {
							entries = append(entries, table.Entry(imageName+" on "+provider, imageName, false, provider))
						}
					}
					for _, imageName := range imagesWithoutVolumes {
						for _, provider := range providersWithVolumes {
							entries = append(entries, table.Entry(imageName+" on "+provider, imageName, false, provider))
						}
						for _, provider := range providersWithoutVolumes {
							entries = append(entries, table.Entry(imageName+" on "+provider, imageName, false, provider))
						}
					}
					logrus.WithField("entries", entries).WithField("imageNames", append(imagesWithVolumes, imagesWithoutVolumes...)).WithField("providers", providersWithVolumes).Infof("ENTRIES TO TEST")
					Context("Build() then Run()", func() {
						table.DescribeTable("running images", func(imageName string, withVolume bool, provider string) {
							compiler := ""
							switch {
							case strings.Contains(imageName, "go"):
								logrus.Infof("found image type GO: %s", imageName)
								compiler = fmt.Sprintf("rump-go-%s", provider)
								break
							case strings.Contains(imageName, "nodejs"):
								logrus.Infof("found image type NODE: %s", imageName)
								compiler = fmt.Sprintf("rump-nodejs-%s", provider)
								break
							case strings.Contains(imageName, "python"):
								logrus.Infof("found image type PYTHON: %s", imageName)
								compiler = fmt.Sprintf("rump-python-%s", provider)
								break
							case strings.Contains(imageName, "war"):
								fallthrough
							case strings.Contains(imageName, "jar"):
								fallthrough
							case strings.Contains(imageName, "java"):
								logrus.Infof("found image type JAVA: %s", imageName)
								compiler = fmt.Sprintf("osv-java-%s", provider)
								break
							default:
								logrus.Panic("unknown image name " + imageName)
							}
							//vsphere -> vmware for compilers
							compiler = strings.Replace(compiler, "vsphere", "vmware", -1)
							compiler = strings.Replace(compiler, "aws", "xen", -1)
							mounts := []string{}
							mountPointsToVols := map[string]string{}
							var err error
							if withVolume {
								mounts = append(mounts, "/data")
								volume, err = helpers.CreateTestDataVolume(daemonUrl, "test_volume_"+imageName, provider)
								Expect(err).ToNot(HaveOccurred())
								mountPointsToVols["/data"] = volume.Id
							}
							image, err = helpers.BuildTestImage(daemonUrl, imageName, compiler, provider, mounts)
							Expect(err).ToNot(HaveOccurred())
							instanceName := imageName
							noCleanup := false
							env := map[string]string{"KEY": "VAL"}
							memoryMb := 256
							instance, err := c.Instances().Run(instanceName, image.Name, mountPointsToVols, env, memoryMb, noCleanup, false)
							Expect(err).ToNot(HaveOccurred())
							instanceIp, err := helpers.WaitForIp(daemonUrl, instance.Id, ipTimeout)
							Expect(err).ToNot(HaveOccurred())
							testInstancePing(instanceIp)
							testInstanceEnv(instanceIp)
							if withVolume {
								testInstanceMount(instanceIp)
							}
						}, entries...)
					})
				})
			})
		})
	})
})
