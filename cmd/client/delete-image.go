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
	"os"

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rmiCmd = &cobra.Command{
	Use:     "delete-image",
	Aliases: []string{"rmi"},
	Short:   "Delete a unikernel image",
	Long: `Deletes an image.
You may specify the image by name or id.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if err := readClientConfig(); err != nil {
				return err
			}
			if imageName == "" {
				return errors.New("must specify --image", nil)
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithFields(logrus.Fields{"host": host, "force": force, "image": imageName}).Info("deleting image")
			if err := client.KernelClient(host).Images().Delete(imageName, force); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			logrus.Errorf("failed deleting image: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(rmiCmd)
	rmiCmd.Flags().StringVar(&imageName, "image", "", "<string,required> name or id of image. Bhojpur Kernel accepts a prefix of the name or id")
	rmiCmd.Flags().BoolVar(&force, "force", false, "<bool, optional> force deleting image in the case that it is running")
}
