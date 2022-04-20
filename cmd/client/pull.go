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

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull an image from a Bhojpur Kernel Image Repository",
	Long: `
Example usage:
kernctl pull --image theirImage --provider virtualbox|qemu|xen

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
		if provider == "" {
			logrus.Fatal("--provider must be set")
		}
		if host == "" {
			host = clientConfig.Host
		}
		if err := client.KernelClient(host).Images().Pull(c, imageName, provider, force); err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(imageName + " pulled")
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
	pullCmd.Flags().StringVar(&imageName, "image", "", "<string,required> image to pull")
	pullCmd.Flags().StringVar(&provider, "provider", "", "<string,required> name of the provider the image is built for")
	pullCmd.Flags().BoolVar(&force, "force", false, "<bool,optional> force overwriting local image of the same name")
}
