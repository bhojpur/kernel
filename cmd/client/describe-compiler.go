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
	"os"

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var describeCompilerCmd = &cobra.Command{
	Use:   "describe-compiler",
	Short: "Describe compiler usage",
	Long: `Describes compiler usage.
You must provide triple base-language-provider to specify what compiler to describe.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			// Determine host.
			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithField("host", host).Info("listing providers")

			// Validate input.
			if base == "" {
				return errors.New("--base must be set", nil)
			}
			if lang == "" {
				return errors.New("--language must be set", nil)
			}
			if provider == "" {
				return errors.New("--provider must be set", nil)
			}
			logrus.WithFields(logrus.Fields{
				"base":     base,
				"lang":     lang,
				"provider": provider,
			}).Info("describe compiler")

			// Ask daemon.
			description, err := client.KernelClient(host).DescribeCompiler(base, lang, provider)
			if err != nil {
				return err
			}

			// Print result to the console.
			fmt.Println(description)

			return nil
		}(); err != nil {
			logrus.Errorf("failed describing compiler: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(describeCompilerCmd)
	describeCompilerCmd.Flags().StringVar(&base, "base", "", "<string,required> name of the unikernel base to use")
	describeCompilerCmd.Flags().StringVar(&lang, "language", "", "<string,required> language the unikernel source is written in")
	describeCompilerCmd.Flags().StringVar(&provider, "provider", "", "<string,required> name of the target infrastructure to compile for")
}
