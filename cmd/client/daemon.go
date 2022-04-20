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
	"os"
	"path/filepath"

	"net/url"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/daemon"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var daemonRuntimeFolder, daemonConfigFile, logFile string
var debugMode, trace bool

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Runs the Bhojpur Kernel backend (daemon)",
	Long: `Run this command to start the Bhojpur Kernel daemon process.
This should normally be left running as a long-running background process.
The daemon requires that docker is installed and running on the your system.
Necessary docker containers must be built for the daemon to work properly;
Run 'make' in the Bhojpur Kernel root directory to build all necessary containers.

Daemon also requires a configuration file with credentials and configuration info
for desired providers.

Example usage:
	kernctl daemon --f ./my-config.yaml --port 12345 --debug --trace --logfile logs.txt

	 # will start the daemon using config file at my-config.yaml
	 # running on port 12345
	 # debug mode activated
	 # trace mode activated
	 # outputting logs to logs.txt
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {

			// set Bhojpur Kernel home
			config.Internal.KernelHome = daemonRuntimeFolder

			if daemonConfigFile == "" {
				daemonConfigFile = filepath.Join(config.Internal.KernelHome, "daemon-config.yaml")
			}

			if err := readDaemonConfig(); err != nil {
				return err
			}

			//don't print vsphere password
			redactions := []string{}
			for _, vsphereConfig := range daemonConfig.Providers.Vsphere {
				redactions = append(redactions, vsphereConfig.VspherePassword, url.QueryEscape(vsphereConfig.VspherePassword))
			}
			logrus.SetFormatter(&kutil.RedactedTextFormatter{
				Redactions: redactions,
			})

			if debugMode {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if trace {
				logrus.AddHook(&kutil.AddTraceHook{true})
			}
			if logFile != "" {
				os.Create(logFile)
				f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
				if err != nil {
					return errors.New(fmt.Sprintf("failed to open log file %s for writing", logFile), err)
				}
				logrus.AddHook(&kutil.TeeHook{f})
			}

			logrus.WithField("config", daemonConfig).Info("Bhojpur Kernel daemon started")
			d, err := daemon.NewKernelDaemon(daemonConfig)
			if err != nil {
				return errors.New("daemon failed to initialize", err)
			}
			d.Run(port)
			return nil
		}(); err != nil {
			logrus.Errorf("running daemon failed: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVar(&daemonRuntimeFolder, "d", getHomeDir()+"/.bhojpur/", "daemon runtime folder - where state is stored. (default is $HOME/.bhojpur/)")
	daemonCmd.Flags().StringVar(&daemonConfigFile, "f", "", "daemon config file (default is {RuntimeFolder}/daemon-config.yaml)")
	daemonCmd.Flags().IntVar(&port, "port", 3000, "<int, optional> listening port for daemon")
	daemonCmd.Flags().BoolVar(&debugMode, "debug", false, "<bool, optional> more verbose logging for the daemon")
	daemonCmd.Flags().BoolVar(&trace, "trace", false, "<bool, optional> add stack trace to daemon logs")
	daemonCmd.Flags().StringVar(&logFile, "logfile", "", "<string, optional> output logs to file (in addition to stdout)")
}

var daemonConfig config.DaemonConfig

func readDaemonConfig() error {
	data, err := ioutil.ReadFile(daemonConfigFile)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read daemon configuration file at " + daemonConfigFile + `\n
		See documentation at https://github.com/bhojpur/kernel/pkg/util for creating daemon config.'`)
		return errors.New(errMsg, err)
	}
	if err := yaml.Unmarshal(data, &daemonConfig); err != nil {
		errMsg := fmt.Sprintf("failed to parse daemon configuration yaml at " + daemonConfigFile + `\n
		Please ensure config file contains valid yaml.'`)
		return errors.New(errMsg, err)
	}
	return nil
}
