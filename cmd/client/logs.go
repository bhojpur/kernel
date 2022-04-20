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
	"bufio"
	"fmt"
	"os"

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var follow, deleteOnDisconnect bool

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "retrieve the logs (stdout) of a unikernel instance",
	Long: `Retrieves logs from a running unikernel instance.

Cannot be used on an instance in powered-off state.
Use the --follow flag to attach to the instance's stdout
Use --delete in combination with --follow to force automatic instance
deletion when the HTTP connection to the instance is broken (by client
disconnect). The --delete flag is typically intended for use with
orchestration software such as cluster managers which may require
a persistent http connection managed instances.

You may specify the instance by name or id.

Example usage:
	kernctl logs --instancce myInstance

	# will return captured stdout from myInstance since boot time

	kernctl logs --instance myInstance --follow --delete

	# will open an HTTP connection between the CLI and Bhojpur Kernel
	backend which streams stdout from the instance to the client
	# when the client disconnects (i.e. with Ctrl+C) Bhojpur Kernel will
	automatically power down and terminate the instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}
			if instanceName == "" {
				return errors.New("must specify --instance", nil)
			}
			if follow {
				logrus.WithFields(logrus.Fields{"host": host, "instance": instanceName}).Info("attaching to instance")
				r, err := client.KernelClient(host).Instances().AttachLogs(instanceName, deleteOnDisconnect)
				if err != nil {
					return err
				}
				reader := bufio.NewReader(r)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					if line != "\n" {
						fmt.Printf(line)
					}
				}
			} else {
				logrus.WithFields(logrus.Fields{"host": host, "instance": instanceName}).Info("getting instance logs")
				data, err := client.KernelClient(host).Instances().GetLogs(instanceName)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", string(data))
			}
			return nil
		}(); err != nil {
			logrus.Errorf("failed retrieving instance logs: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVar(&instanceName, "instance", "", "<string,required> name or id of instance. Bhojpur Kernel accepts a prefix of the name or id")
	logsCmd.Flags().BoolVar(&follow, "follow", false, "<bool,optional> follow stdout of instance as it is printed")
	logsCmd.Flags().BoolVar(&deleteOnDisconnect, "delete", false, "<bool,optional> use this flag with the --follow flag to trigger automatic deletion of instance after client closes the http connection")
}
