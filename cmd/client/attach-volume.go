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

var mountPoint string

var attachCmd = &cobra.Command{
	Use:     "attach-volume",
	Aliases: []string{"attach"},
	Short:   "Attach a volume to a stopped instance",
	Long: `Attaches a volume to a stopped instance at a specified mount point.
You specify the volume by name or id.

The volume must be attached to an available mount point on the instance.
Mount points are image-specific, and are determined when the image is compiled.

For a list of mount points on the image for this instance, run kernctl images, or
kernctl describe image

If the specified mount point is occupied by another volume, the command will result
in an error
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if err := readClientConfig(); err != nil {
				return err
			}
			if volumeName == "" {
				return errors.New("must specify --volume", nil)
			}
			if instanceName == "" {
				return errors.New("must specify --instanceName", nil)
			}
			if mountPoint == "" {
				return errors.New("must specify --mountPoint", nil)
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithFields(logrus.Fields{"host": host, "instanceName": instanceName, "volume": volumeName, "mountPoint": mountPoint}).Info("attaching volume")
			if err := client.KernelClient(host).Volumes().Attach(volumeName, instanceName, mountPoint); err != nil {
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
	RootCmd.AddCommand(attachCmd)
	attachCmd.Flags().StringVar(&volumeName, "volume", "", "<string,required> name or id of volume to attach. Bhojpur Kernel accepts a prefix of the name or id")
	attachCmd.Flags().StringVar(&instanceName, "instance", "", "<string,required> name or id of instance to attach to. Bhojpur Kernel accepts a prefix of the name or id")
	attachCmd.Flags().StringVar(&mountPoint, "mountPoint", "", "<string,required> mount path for volume. this should reflect the mappings specified on the image. run 'kernctl describe-image' to see expected mount points for the image")
	attachCmd.Flags().BoolVar(&force, "force", false, "<bool, optional> force deleting volume in the case that it is running")
}
