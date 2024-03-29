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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var clientConfigFile, hubConfigFile, host string
var port int

var RootCmd = &cobra.Command{
	Use:   "kernctl",
	Short: "The unikernel compilation, deployment, and management tool",
	Long: `Bhojpur Kernel is a tool for compiling application source code
into bootable disk images. Bhojpur Kernel also runs and manages unikernel
instances across infrastructures.
Create a client configuration file with 'kernctl target'
You may set a custom client configuration file
with the global flag --client-config=<path>`,
}

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&clientConfigFile, "client-config", getHomeDir()+"/.bhojpur/client-config.yaml", "client config file")
	RootCmd.PersistentFlags().StringVar(&hubConfigFile, "hub-config", getHomeDir()+"/.bhojpur/hub-config.yaml", "hub config file")
	RootCmd.PersistentFlags().StringVar(&host, "host", "", "<string, optional>: host/ip address of the host running the Bhojpur Kernel daemon")
	targetCmd.Flags().IntVar(&port, "port", 3000, "<int, optional>: port the daemon is running on (default: 3000)")
}

var clientConfig config.ClientConfig

func readClientConfig() error {
	data, err := ioutil.ReadFile(clientConfigFile)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read client configuration file at " + clientConfigFile + `
Try setting your config with 'kernctl target --host HOST_URL'`)
		return err
	}
	data = bytes.Replace(data, []byte("\n"), []byte{}, -1)
	if err := yaml.Unmarshal(data, &clientConfig); err != nil {
		logrus.WithError(err).Errorf("failed to parse client configuration yaml at " + clientConfigFile + `
Please ensure config file contains valid yaml.'\n
Try setting your config with 'kernctl target --host HOST_URL'`)
		return err
	}
	return nil
}

type imageSlice []*types.Image

func (p imageSlice) Len() int           { return len(p) }
func (p imageSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p imageSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p imageSlice) Sort() { sort.Sort(p) }

func printImages(images ...*types.Image) {
	sortedImages := make(imageSlice, len(images))
	for i, image := range images {
		sortedImages[i] = image
	}
	sortedImages.Sort()
	fmt.Printf("%-20s %-20s %-15s %-30s %-6s %-20s\n", "NAME", "ID", "INFRASTRUCTURE", "CREATED", "SIZE(MB)", "MOUNTPOINTS")
	for _, image := range sortedImages {
		printImage(image)
	}
}

func printImage(image *types.Image) {
	for i, deviceMapping := range image.RunSpec.DeviceMappings {
		//ignore root device mount point
		if deviceMapping.MountPoint == "/" {
			image.RunSpec.DeviceMappings = append(image.RunSpec.DeviceMappings[:i], image.RunSpec.DeviceMappings[i+1:]...)
		}
	}
	if len(image.RunSpec.DeviceMappings) == 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d \n", image.Name, image.Id, image.Infrastructure, image.Created.String(), image.SizeMb)
	} else if len(image.RunSpec.DeviceMappings) > 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d %-20.20s\n", image.Name, image.Id, image.Infrastructure, image.Created.String(), image.SizeMb, image.RunSpec.DeviceMappings[0].MountPoint)
		if len(image.RunSpec.DeviceMappings) > 1 {
			for i := 1; i < len(image.RunSpec.DeviceMappings); i++ {
				fmt.Printf("%102s\n", image.RunSpec.DeviceMappings[i].MountPoint)
			}
		}
	}
}

type userImageSlice []*types.UserImage

func (p userImageSlice) Len() int           { return len(p) }
func (p userImageSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p userImageSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p userImageSlice) Sort() { sort.Sort(p) }

func printUserImages(images ...*types.UserImage) {
	sortedImages := make(userImageSlice, len(images))
	for i, image := range images {
		sortedImages[i] = image
	}
	sortedImages.Sort()
	fmt.Printf("%-20s %-20s %-15s %-30s %-6s %-20s\n", "NAME", "OWNER", "INFRASTRUCTURE", "CREATED", "SIZE(MB)", "MOUNTPOINTS")
	for _, image := range sortedImages {
		printUserImage(image)
	}
}

func printUserImage(image *types.UserImage) {
	for i, deviceMapping := range image.RunSpec.DeviceMappings {
		//ignore root device mount point
		if deviceMapping.MountPoint == "/" {
			image.RunSpec.DeviceMappings = append(image.RunSpec.DeviceMappings[:i], image.RunSpec.DeviceMappings[i+1:]...)
		}
	}
	if len(image.RunSpec.DeviceMappings) == 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d \n", image.Name, image.Owner, image.Infrastructure, image.Created.String(), image.SizeMb)
	} else if len(image.RunSpec.DeviceMappings) > 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d %-20.20s\n", image.Name, image.Owner, image.Infrastructure, image.Created.String(), image.SizeMb, image.RunSpec.DeviceMappings[0].MountPoint)
		if len(image.RunSpec.DeviceMappings) > 1 {
			for i := 1; i < len(image.RunSpec.DeviceMappings); i++ {
				fmt.Printf("%102s\n", image.RunSpec.DeviceMappings[i].MountPoint)
			}
		}
	}
}

type instanceSlice []*types.Instance

func (p instanceSlice) Len() int           { return len(p) }
func (p instanceSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p instanceSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p instanceSlice) Sort() { sort.Sort(p) }

func printInstances(instances ...*types.Instance) {
	sortedInstances := make(instanceSlice, len(instances))
	for i, instance := range instances {
		sortedInstances[i] = instance
	}
	sortedInstances.Sort()
	fmt.Printf("%-15s %-20s %-14s %-30s %-20s %-15s %-12s\n",
		"NAME", "ID", "INFRASTRUCTURE", "CREATED", "IMAGE", "IPADDRESS", "STATE")
	for _, instance := range sortedInstances {
		printInstance(instance)
	}
}

func printInstance(instance *types.Instance) {
	fmt.Printf("%-15.15s %-20.20s %-14.14s %-30.30s %-20.20v %-15.15s %-12.12s\n",
		instance.Name, instance.Id, instance.Infrastructure, instance.Created.String(), instance.ImageId, instance.IpAddress, instance.State)
}

func printVolumes(volume ...*types.Volume) {
	fmt.Printf("%-15.15s %-15.15s %-14.14s %-30.30s %-20.20v %-12.12s\n",
		"NAME", "ID", "INFRASTRUCTURE", "CREATED", "ATTACHED-INSTANCE", "SIZE(MB)")
	for _, volume := range volume {
		printVolume(volume)
	}
}

func printVolume(volume *types.Volume) {
	fmt.Printf("%-15.15s %-15.15s %-14.14s %-30.30s %-20.20v %-12.12d\n",
		volume.Name, volume.Id, volume.Infrastructure, volume.Created.String(), volume.Attachment, volume.SizeMb)
}
