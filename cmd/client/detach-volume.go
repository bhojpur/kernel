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

var detachCmd = &cobra.Command{
	Use:     "detach-volume",
	Aliases: []string{"detach"},
	Short:   "Detach an attached volume from a stopped instance",
	Long: `Detaches a volume to a stopped instance at a specified mount point.
You specify the volume by name or id.

After detaching the volume, the volume can be mounted to another instance.

If the instance is not stopped, detach will result in an error.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if err := readClientConfig(); err != nil {
				return err
			}
			if volumeName == "" {
				return errors.New("must specify --volume", nil)
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithFields(logrus.Fields{"host": host, "volume": volumeName}).Info("detaching volume")
			if err := client.KernelClient(host).Volumes().Detach(volumeName); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			logrus.Errorf("failed deleting volume: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(detachCmd)
	detachCmd.Flags().StringVar(&volumeName, "volume", "", "<string,required> name or id of volume to detach. Bhojpur Kernel accepts a prefix of the name or id")
}
