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

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var show bool

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Configure Bhojpur Kernel daemon URL for CLI client commands",
	Long: `Sets the host url of the Bhojpur Kernel daemon for CLI commands.
If running Bhojpur Kernel locally, use 'kernctl target --host localhost'

args:
--host: <string, required>: host/ip address of the host running the Bhojpur Kernel daemon
--port: <int, optional>: port the daemon is running on (default: 3000)

--show: <bool,optional>: shows the current target that is set`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if show {
				if err := readClientConfig(); err != nil {
					return err
				}
				logrus.Infof("Current target: %s", clientConfig.Host)
				return nil
			}
			if host == "" {
				return errors.New("--host must be set for target", nil)
			}
			if err := setClientConfig(host, port); err != nil {
				return errors.New("failed to save target to config file", err)
			}
			logrus.Infof("target set: %s:%v", host, port)
			return nil
		}(); err != nil {
			logrus.Errorf("failed running target: %v", err)
			os.Exit(-1)
		}
	},
}

func setClientConfig(host string, port int) error {
	data, err := yaml.Marshal(config.ClientConfig{Host: fmt.Sprintf("%s:%v", host, port)})
	if err != nil {
		return errors.New("failed to convert config to yaml string ", err)
	}
	os.MkdirAll(filepath.Dir(clientConfigFile), 0755)
	if err := ioutil.WriteFile(clientConfigFile, data, 0644); err != nil {
		return errors.New("failed writing config to file "+clientConfigFile, err)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(targetCmd)
	targetCmd.Flags().BoolVar(&show, "show", false, "<bool,optional>: shows the current target that is set")
}
