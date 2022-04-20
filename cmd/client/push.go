package cmd

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
	"fmt"
	"io/ioutil"

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an image to a Bhojpur Kernel Image Repository",
	Long: `
Example usage:
kernctl push --image myImage

Requires that you first authenticate to a Bhojpur Kernel image repository with 'kernctl login'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := readClientConfig(); err != nil {
			logrus.Fatal(err)
		}
		c, err := getHubConfig()
		if err != nil {
			logrus.Fatal(err)
		}
		if imageName == "" {
			logrus.Fatal("--image must be set")
		}
		if host == "" {
			host = clientConfig.Host
		}
		if err := client.KernelClient(host).Images().Push(c, imageName); err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(imageName + " pushed")
	},
}

func getHubConfig() (config.HubConfig, error) {
	var c config.HubConfig
	data, err := ioutil.ReadFile(hubConfigFile)
	if err != nil {
		return c, errors.New("reading "+hubConfigFile, err)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, errors.New("failed to convert config from yaml", err)
	}
	return c, nil
}

func init() {
	RootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVar(&imageName, "image", "", "<string,required> image to push")
}
