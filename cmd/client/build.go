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

	"github.com/bhojpur/kernel/pkg/client"
	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var name, sourcePath, base, lang, provider, runArgs string
var mountPoints []string
var force, noCleanup bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a unikernel image from source code files",
	Long: `Compiles source files into a runnable unikernel image.

Images must be compiled using a specific compiler. Compilers are composed of three elements:
 **Base**: The Unikernel Base to build the image on. These include different unikernel implementations, such as [rumprun](https://github.com/rumpkernel/rumprun), [IncludeOS](https://gitter.im/hioa-cs/IncludeOS), etc.
 **Language**: The language/runtime the image should be built with. e.g. for a Golang project, specify '-language go'. Languages supported depend on the unikernel base.
 **Provider**: The (cloud/hypervisor) provider the image should be built for. Supported providers depend on the unikernel base.

In order to see a list of all supported base/language/provider combinations, run 'kernctl compilers'

Images must be compiled for a specific provider, specified with the '--provider' flag
To see a list of available providers, run 'unik providers'

A unikernel base that is compatible with the provider must be specified with the '--base' flag.
A language runtime that is compatible with the base must be specified with the '--language' flag.
To see a table of all compatible base-language-provider combinations, run 'kernctl compilers'

If you wish to attach volumes to instances of an image, the image must be compiled in advance
with a list of the expected mount points. e.g. for an application that reads from a '/data' folder,
the unikernel should be compiled with the flag '--mount /data'

Runtime arguments to be passed to your unikernel must also be specified at compile time.
You can specify arguments as a single string passed to the '--args' flag

Image names must be unique. If an image exists with the same name, you can force overwriting with the --force flag

Example usage:
	kernctl build --name myUnikernel --path ./myApp/src --base rump --language go --provider aws --mountpoint /foo --mountpoint /bar --args 'arg1 arg2 arg3' --force

	* will create a Go unikernel named myUnikernel using the sources found in ./myApp/src,
	* compiled using rumprun targeting AWS infrastructure,
	* expecting a volume to be mounted at /foo at runtime,
	* expecting another volume to be mounted at /bar at runtime,
	* passing 'arg1 arg2 arg3' as arguments to the application when it is run,
	* and deleting any previous existing instances and image for the name myUnikernel before compiling

Another example (using only the required parameters):
	kernctl build -name anotherUnikernel -path ./anotherApp/src --base includeos --language cpp --provider virtualbox
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if name == "" {
				return errors.New("--name must be set", nil)
			}
			if sourcePath == "" {
				return errors.New("--path must be set", nil)
			}
			if base == "" {
				return errors.New("--base must be set", nil)
			}
			if lang == "" {
				return errors.New("--language must be set", nil)
			}
			if provider == "" {
				return errors.New("--provider must be set", nil)
			}
			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}
			logrus.WithFields(logrus.Fields{
				"name":        name,
				"path":        sourcePath,
				"base":        base,
				"language":    lang,
				"provider":    provider,
				"args":        runArgs,
				"mountPoints": mountPoints,
				"force":       force,
				"host":        host,
			}).Infof("running Bhojpur Kernel build")
			sourceTar, err := ioutil.TempFile("", "sources.tar.gz.")
			if err != nil {
				logrus.WithError(err).Error("failed to create tmp tar file")
			}
			defer os.Remove(sourceTar.Name())
			if err := kos.Compress(sourcePath, sourceTar.Name()); err != nil {
				return errors.New("failed to tar sources", err)
			}
			logrus.Infof("Application packaged as tarball: %s\n", sourceTar.Name())
			image, err := client.KernelClient(host).Images().Build(name, sourceTar.Name(), base, lang, provider, runArgs, mountPoints, force, noCleanup)
			if err != nil {
				return errors.New("building image failed", err)
			}
			printImages(image)
			return nil
		}(); err != nil {
			logrus.Errorf("build failed: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&name, "name", "", "<string,required> name to give the unikernel. must be unique")
	buildCmd.Flags().StringVar(&sourcePath, "path", "", "<string,required> path to root application sources folder")
	buildCmd.Flags().StringVar(&base, "base", "", "<string,required> name of the unikernel base to use")
	buildCmd.Flags().StringVar(&lang, "language", "", "<string,required> language the unikernel source is written in")
	buildCmd.Flags().StringVar(&provider, "provider", "", "<string,required> name of the target infrastructure to compile for")
	buildCmd.Flags().StringVar(&runArgs, "args", "", "<string,optional> to be passed to the unikernel at runtime")
	buildCmd.Flags().StringSliceVar(&mountPoints, "mountpoint", []string{}, "<string,repeated> specify up to 8 mount points for volumes")
	buildCmd.Flags().BoolVar(&force, "force", false, "<bool, optional> force overwriting a previously existing")
	buildCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "<bool, optional> for debugging; do not clean up artifacts for images that fail to build")
}
