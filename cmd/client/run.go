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
	"strings"

	"bufio"
	"net"

	"github.com/bhojpur/kernel/pkg/client"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var instanceName, imageName string
var volumes, envPairs []string
var instanceMemory, debugPort int

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a unikernel instance from a compiled image",
	Long: `Deploys a running instance from a Bhojpur Kernel compiled unikernel disk image.
The instance will be deployed on the provider the image was compiled for.
e.g. if the image was compiled for virtualbox, Bhojpur Kernel will attempt to deploy
the image on the configured virtualbox environment.

'kernctl run' requires a Bhojpur Kernel managed volume (see 'kernctl volumes' and 'kernctl create volume')
to be attached and mounted to each mount point specified at image compilation time.
This means that if the image was compiled with two mount points, /data1 and /data2,
'kernctl run' requires 2 available volumes to be attached to the instance at runtime, which
must be specified with the flags --vol SOME_VOLUME_NAME:/data1 --vol ANOTHER_VOLUME_NAME:/data2
If no mount points are required for the image, volumes cannot be attached.

environment variables can be set at runtime through the use of the -env flag.

Example usage:
	kernctl run --instanceName newInstance --imageName myImage --vol myVol:/mount1 --vol yourVol:/mount2 --env foo=bar --env another=one --instanceMemory 1234

	# will create and run an instance of myImage on the provider environment myImage is compiled for
	# instance will be named newInstance
	# instance will attempt to mount Bhojpur Kernel managed volume myVol to /mount1
	# instance will attempt to mount Bhojpur Kernel managed volume yourVol to /mount2
	# instance will boot with env variable 'foo' set to 'bar'
	# instance will boot with env variable 'another' set to 'one'
	# instance will get 1234 MB of memory

	# note that run must take exactly one --vol argument for each mount point defined in the image specification
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := func() error {
			if instanceName == "" {
				return errors.New("--instanceName must be set", nil)
			}
			if imageName == "" {
				return errors.New("--imageName must be set", nil)
			}
			if err := readClientConfig(); err != nil {
				return err
			}
			if host == "" {
				host = clientConfig.Host
			}

			mountPointsToVols := make(map[string]string)
			for _, vol := range volumes {
				pair := strings.SplitN(vol, ":", 2)
				if len(pair) != 2 {
					return errors.New(fmt.Sprintf("invalid format for vol flag: %s", vol), nil)
				}
				volId := pair[0]
				mnt := pair[1]
				mountPointsToVols[mnt] = volId
			}

			env := make(map[string]string)
			for _, e := range envPairs {
				pair := strings.Split(e, "=")
				if len(pair) != 2 {
					return errors.New(fmt.Sprintf("invalid format for env flag: %s", e), nil)
				}
				key := pair[0]
				val := pair[1]
				env[key] = val
			}

			logrus.WithFields(logrus.Fields{
				"instanceName": instanceName,
				"imageName":    imageName,
				"env":          env,
				"mounts":       mountPointsToVols,
				"host":         host,
			}).Infof("running kernctl run")
			instance, err := client.KernelClient(host).Instances().Run(instanceName, imageName, mountPointsToVols, env, instanceMemory, noCleanup, debugMode)
			if err != nil {
				return errors.New("running image failed: %v", err)
			}
			printInstances(instance)
			if debugMode {
				logrus.Infof("attaching debugger to instance %s ...", instance.Name)
				connectDebugger()
			}
			return nil
		}(); err != nil {
			logrus.Errorf("failed running instance: %v", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&instanceName, "instanceName", "", "<string,required> name to give the instance. must be unique")
	runCmd.Flags().StringVar(&imageName, "imageName", "", "<string,required> image to use")
	runCmd.Flags().StringSliceVar(&envPairs, "env", []string{}, "<string,repeated> set any number of environment variables for the instance. must be in the format KEY=VALUE")
	runCmd.Flags().StringSliceVar(&volumes, "vol", []string{}, `<string,repeated> each --vol flag specifies one volume id and the corresponding mount point to attach
	to the instance at boot time. volumes must be attached to the instance for each mount point expected by the image.
	run 'kernctl image <image_name>' to see the mount points required for the image.
	specified in the format 'volume_id:mount_point'`)
	runCmd.Flags().IntVar(&instanceMemory, "instanceMemory", 0, "<int, optional> amount of memory (in MB) to assign to the instance. if none is given, the provider default will be used")
	runCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "<bool, optional> for debugging; do not clean up artifacts for instances that fail to launch")
	runCmd.Flags().BoolVar(&debugMode, "debug-mode", false, "<bool, optional> runs the instance in Debug mode so GDB can be attached. Currently only supported on QEMU provider")
	runCmd.Flags().IntVar(&debugPort, "debug-port", 3001, "<int, optional> target port for debugger tcp connections. used in conjunction with --debug-mode")
}

func connectDebugger() {
	addr := fmt.Sprintf("%v:%v", strings.Split(host, ":")[0], debugPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logrus.Errorf("failed to initiate tcp connection: %v", err)
		os.Exit(-1)
	}
	if _, err := conn.Write([]byte("GET / HTTP/1.0\r\n\r\n")); err != nil {
		logrus.Errorf("failed to initialize debgger connection: %v", err)
		os.Exit(-1)
	}

	go func() {
		reader := bufio.NewReader(conn)
		for {
			data, err := reader.ReadBytes('\n')
			if err != nil {
				logrus.Errorf("disconnected from debugger: %v", err)
				os.Exit(0)
			}
			fmt.Print(string(data))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	logrus.Infof("Connected to %s", host)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			logrus.Errorf("failed reading stdin: %v", err)
			os.Exit(-1)
		}
		if _, err := conn.Write(data); err != nil {
			logrus.Errorf("writing to tcp connection: %v", err)
			os.Exit(-1)
		}
	}
}
