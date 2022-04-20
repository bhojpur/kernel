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

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "List available Bhojur Kernel managed volumes",
	Long: `Lists all available Bhojpur Kernel managed volumes across providers.

ATTACHED-INSTANCE gives the instance ID of the instance a volume
is attached to, if any. Only volumes that have no attachment are
available to be attached to an instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithField("host", host).Info("listing volumes")
			volumes, err := client.KernelClient(host).Volumes().All()
			if err != nil {
				return errors.New("listing volumes failed", err)
			}
			printVolumes(volumes...)
			return nil
		}(); err != nil {
			logrus.Errorf("failed listing volumes: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(volumesCmd)
}
