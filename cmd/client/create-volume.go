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
	"io/ioutil"
	"os"
	"strings"

	"github.com/bhojpur/kernel/pkg/client"
	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var data string
var size int
var volumeType string
var rawVolume bool

const (
	VolTypeExt2 = "ext2"
)

var cvCmd = &cobra.Command{
	Use:   "create-volume",
	Short: "Create a Bhojpur Kernel managed data volume",
	Long: `Create a data volume which can be attached to and detached from
Bhojpur Kernel managed instances.

Volumes can be created from a directory, which will copy the contents
of the directory onto the voume. Empty volume can also be created.

Volumes will persist after instances are deleted, allowing application data
to be persisted beyond the lifecycle of individual instances.

If specifying a data folder (with --data), specifying a size for the volume is
not necessary. Bhojpur Kernel will automatically size the volume to fit the data provided.
A larger volume can be requested with the --size flag.

If no data directory is provided, --size is a required parameter to specify the
desired size for the empty volume to be createad.

Volumes are created for a specific provider, specified with the --provider flag.
Volumes can only be attached to instances of the same provider type.
To see a list of available providers, run 'kernctl providers'

Volume names must be unique. If a volume exists with the same name, you will be
required to remove the volume with 'kernctl delete-volume' before the new volume
can be created.

--size parameter uses MB

Example usage:
	kernctl create-volume --name myVolume --data ./myApp/data --provider aws

	# will create an EBS-backed AWS volume named myVolume using the data found in ./myApp/src,
	# the size will be either 1GB (the default minimum size on AWS) or greater, if the size of the
	volume is greater


Another example (empty volume):
	kernctl create-volume -name anotherVolume --size 500 -provider vsphere

	# will create a 500mb sparse vmdk file and upload it to the vsphere datastore,
	where it can be attached to a vsphere instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if name == "" {
				return errors.New("--name must be set", nil)
			}
			if data == "" && size == 0 {
				return errors.New("either --data or --size must be set", nil)
			}
			if provider == "" {
				return errors.New("--provider must be set", nil)
			}
			if volumeType == "" {
				volumeType = VolTypeExt2
			} else {
				volumeType = strings.ToLower(volumeType)
			}

			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithFields(logrus.Fields{
				"name":       name,
				"data":       data,
				"size":       size,
				"provider":   provider,
				"host":       host,
				"volumeType": volumeType,
			}).Infof("creating volume")
			if data != "" {
				dataTar, err := ioutil.TempFile("", "data.tar.gz.")
				if err != nil {
					logrus.WithError(err).Error("failed to create tmp tar file")
				}
				if false {
					defer os.Remove(dataTar.Name())
				}
				if err := kos.Compress(data, dataTar.Name()); err != nil {
					return errors.New("failed to tar data", err)
				}
				data = dataTar.Name()
				logrus.Infof("Data packaged as tarball: %s\n", dataTar.Name())
			}

			volume, err := client.KernelClient(host).Volumes().Create(name, data, provider, rawVolume, size, volumeType, noCleanup)

			if err != nil {
				return errors.New("creatinv volume image failed", err)
			}
			printVolumes(volume)
			return nil
		}(); err != nil {
			logrus.Errorf("create-volume failed: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(cvCmd)
	cvCmd.Flags().StringVar(&name, "name", "", "<string,required> name to give the volume. must be unique")
	cvCmd.Flags().StringVar(&data, "data", "", "<string,special> path to data folder (or file if --raw is provided). optional if --size is provided")
	cvCmd.Flags().BoolVar(&rawVolume, "raw", false, "<bool,optional> if true then then data is expected to be a file that will be used as is. if false (default) data should point to a folder which will be turned into a volume.")
	cvCmd.Flags().IntVar(&size, "size", 0, "<int,special> size to create volume in MB. optional if --data is provided")
	cvCmd.Flags().StringVar(&provider, "provider", "", "<string,required> name of the target infrastructure to compile for")
	cvCmd.Flags().StringVar(&volumeType, "type", "", "<string,optional> FS type of the volume. ext2 or FAT are supported. defaults to ext2")

	cvCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "<bool, optional> for debugging; do not clean up artifacts for volumes that fail to build")
}
