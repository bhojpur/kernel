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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Bhojpur Kernel Repository to pull & push images",
	Run: func(cmd *cobra.Command, args []string) {
		defaultUrl := "http://hub.bhojpur.net"
		reader := bufio.NewReader(os.Stdin)
		if err := func() error {
			fmt.Printf("Bhojpur Kernel Hub Repository URL [%v]: ", defaultUrl)
			url, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			url = strings.Trim(url, "\n")
			if len(url) < 1 {
				url = defaultUrl
			}
			fmt.Printf("Username: ")
			user, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			fmt.Printf("Password: ")
			pass, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if err := setHubConfig(url, strings.Trim(user, "\n"), strings.Trim(pass, "\n")); err != nil {
				return err
			}
			fmt.Printf("using url %v\n", url)
			return nil
		}(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

func setHubConfig(url, user, pass string) error {
	data, err := yaml.Marshal(config.HubConfig{URL: url, Username: user, Password: pass})
	if err != nil {
		return errors.New("failed to convert config to yaml string ", err)
	}
	os.MkdirAll(filepath.Dir(hubConfigFile), 0755)
	if err := ioutil.WriteFile(hubConfigFile, data, 0644); err != nil {
		return errors.New("failed writing config to file "+clientConfigFile, err)
	}
	return nil
}
