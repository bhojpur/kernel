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

var configureCmd = &cobra.Command{
	Use:   "configure [--provider PROVIDER-NAME]",
	Short: "A generate configuration file for daemon ('daemon.yaml')",
	Long: `An interactive command to help walk you through the process of creating or changing a configuration file for Bhojpur Kernel.
Can be used to configure an individual provider, or any number of providers.

Usage:
kernctl configure # will iterate through all possible providers and ask if user wants to configure
-or-
kernctl configure --provider PROVIDER
where provider is one of the following:
aws
gcloud
openstack
qemu
ukvm
virtualbox
vsphere
xen

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if daemonConfigFile == "" {
			daemonConfigFile = os.Getenv("HOME") + "/.bhojpur/daemon-config.yaml"
		}
		readDaemonConfig()
		reader := bufio.NewReader(os.Stdin)
		var configFunc func() error
		switch strings.ToLower(provider) {
		case "aws":
			configFunc = func() error {
				if err := doAwsConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "gcloud":
			configFunc = func() error {
				if err := doGcloudConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "openstack":
			configFunc = func() error {
				if err := doOpenstackConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "qemu":
			configFunc = func() error {
				if err := doQemuConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "ukvm":
			configFunc = func() error {
				if err := doUkvmConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "virtualbox":
			configFunc = func() error {
				if err := doVirtualboxConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "vsphere":
			configFunc = func() error {
				if err := doVsphereConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "xen":
			configFunc = func() error {
				if err := doXenConfig(reader); err != nil {
					return err
				}
				return nil
			}
		case "":
			configFunc = func() error {
				if err := doAwsConfig(reader); err != nil {
					return err
				}
				if err := doGcloudConfig(reader); err != nil {
					return err
				}
				if err := doOpenstackConfig(reader); err != nil {
					return err
				}
				if err := doQemuConfig(reader); err != nil {
					return err
				}
				if err := doUkvmConfig(reader); err != nil {
					return err
				}
				if err := doVirtualboxConfig(reader); err != nil {
					return err
				}
				if err := doVsphereConfig(reader); err != nil {
					return err
				}
				if err := doXenConfig(reader); err != nil {
					return err
				}
				return nil
			}
		}
		if err := configFunc(); err != nil {
			logrus.Fatal(err)
		}
		if err := writeDaemonConfig(); err != nil {
			logrus.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringVar(&provider, "provider", "", "<string,optional> provider to configure. if not given, Bhojpur Kernel will iterate through each possible provider to configure")
	configureCmd.Flags().StringVar(&daemonConfigFile, "f", os.Getenv("HOME")+"/.bhojpur/daemon-config.yaml", "<string, optional> output path for daemon config file")
}

func writeDaemonConfig() error {
	data, err := yaml.Marshal(daemonConfig)
	if err != nil {
		return errors.New("failed to convert config to yaml string ", err)
	}
	os.MkdirAll(filepath.Dir(daemonConfigFile), 0755)
	if err := ioutil.WriteFile(daemonConfigFile, data, 0644); err != nil {
		return errors.New("failed writing config to file "+daemonConfigFile, err)
	}
	return nil
}

func doAwsConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with AWS? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Aws) < 1 {
			daemonConfig.Providers.Aws = append(daemonConfig.Providers.Aws, config.Aws{})
		}
		if daemonConfig.Providers.Aws[0].Name == "" {
			daemonConfig.Providers.Aws[0].Name = "aws-configuration"
		}
		fmt.Printf("AWS region where to deploy unikernels [%s]: ", daemonConfig.Providers.Aws[0].Region)
		region, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		region = strings.TrimSuffix(region, "\n")
		if region != "" {
			daemonConfig.Providers.Aws[0].Region = region
		}
		fmt.Printf("AWS availability zone where to deploy unikernels (must be within region) [%s]: ", daemonConfig.Providers.Aws[0].Zone)
		zone, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		zone = strings.TrimSuffix(zone, "\n")
		if zone != "" {
			daemonConfig.Providers.Aws[0].Zone = zone
		}
	}
	return nil
}

func doGcloudConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with Google Compute Engine? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Gcloud) < 1 {
			daemonConfig.Providers.Gcloud = append(daemonConfig.Providers.Gcloud, config.Gcloud{})
		}
		if daemonConfig.Providers.Gcloud[0].Name == "" {
			daemonConfig.Providers.Gcloud[0].Name = "gcloud-configuration"
		}
		fmt.Printf("GCloud project id under which to deploy unikernels [%s]: ", daemonConfig.Providers.Gcloud[0].ProjectID)
		projectId, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		projectId = strings.TrimSuffix(projectId, "\n")
		if projectId != "" {
			daemonConfig.Providers.Gcloud[0].ProjectID = projectId
		}
		fmt.Printf("Gcloud availability zone where to deploy unikernels (must be within region) [%s]: ", daemonConfig.Providers.Gcloud[0].Zone)
		zone, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		zone = strings.TrimSuffix(zone, "\n")
		if zone != "" {
			daemonConfig.Providers.Gcloud[0].Zone = zone
		}
	}
	return nil
}

func doOpenstackConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with OpenStack? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Openstack) < 1 {
			daemonConfig.Providers.Openstack = append(daemonConfig.Providers.Openstack, config.Openstack{})
		}
		if daemonConfig.Providers.Openstack[0].Name == "" {
			daemonConfig.Providers.Openstack[0].Name = "Openstack-configuration"
		}
		fmt.Printf("OpenStack username for authentication [%s]: ", daemonConfig.Providers.Openstack[0].UserName)
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSuffix(username, "\n")
		if username != "" {
			daemonConfig.Providers.Openstack[0].UserName = username
		}
		fmt.Printf("OpenStack password for authentication [%s]: ", daemonConfig.Providers.Openstack[0].Password)
		password, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		password = strings.TrimSuffix(password, "\n")
		if password != "" {
			daemonConfig.Providers.Openstack[0].Password = password
		}
		fmt.Printf("OpenStack authentication url [%s]: ", daemonConfig.Providers.Openstack[0].AuthUrl)
		authUrl, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		authUrl = strings.TrimSuffix(authUrl, "\n")
		if authUrl != "" {
			daemonConfig.Providers.Openstack[0].AuthUrl = authUrl
		}
		fmt.Printf("OpenStack tenant id [%s]: ", daemonConfig.Providers.Openstack[0].TenantId)
		tenantId, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		tenantId = strings.TrimSuffix(tenantId, "\n")
		if tenantId != "" {
			daemonConfig.Providers.Openstack[0].TenantId = tenantId
		}
		fmt.Printf("OpenStack project name [%s]: ", daemonConfig.Providers.Openstack[0].ProjectName)
		projectName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		projectName = strings.TrimSuffix(projectName, "\n")
		if projectName != "" {
			daemonConfig.Providers.Openstack[0].ProjectName = projectName
		}
		fmt.Printf("OpenStack region id [%s]: ", daemonConfig.Providers.Openstack[0].RegionId)
		regionId, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		regionId = strings.TrimSuffix(regionId, "\n")
		if regionId != "" {
			daemonConfig.Providers.Openstack[0].RegionId = regionId
		}
		fmt.Printf("OpenStack network uuid [%s]: ", daemonConfig.Providers.Openstack[0].NetworkUUID)
		networkUUID, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		networkUUID = strings.TrimSuffix(networkUUID, "\n")
		if networkUUID != "" {
			daemonConfig.Providers.Openstack[0].NetworkUUID = networkUUID
		}
	}
	return nil
}

func doQemuConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with QEMU? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Qemu) < 1 {
			daemonConfig.Providers.Qemu = append(daemonConfig.Providers.Qemu, config.Qemu{})
		}
		if daemonConfig.Providers.Qemu[0].Name == "" {
			daemonConfig.Providers.Qemu[0].Name = "Qemu-configuration"
		}
		fmt.Print("Run QEMU unikernels in nograpic mode [y/N]? (recommended for non-graphical environments): ")
		nographic, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		nographic = strings.TrimSuffix(nographic, "\n")
		if nographic == "y" {
			daemonConfig.Providers.Qemu[0].NoGraphic = true
		}
	}
	return nil
}

func doUkvmConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with uKVM? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Ukvm) < 1 {
			daemonConfig.Providers.Ukvm = append(daemonConfig.Providers.Ukvm, config.Ukvm{})
		}
		if daemonConfig.Providers.Ukvm[0].Name == "" {
			daemonConfig.Providers.Ukvm[0].Name = "Ukvm-configuration"
		}
		fmt.Printf("Name of tap device to attach to uKVM unikernels [%s]: ", daemonConfig.Providers.Ukvm[0].Tap)
		tapDevice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		tapDevice = strings.TrimSuffix(tapDevice, "\n")
		if tapDevice != "" {
			daemonConfig.Providers.Ukvm[0].Tap = tapDevice
		}
	}
	return nil
}

func doVirtualboxConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with VirtualBox? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Virtualbox) < 1 {
			daemonConfig.Providers.Virtualbox = append(daemonConfig.Providers.Virtualbox, config.Virtualbox{})
		}
		if daemonConfig.Providers.Virtualbox[0].Name == "" {
			daemonConfig.Providers.Virtualbox[0].Name = "Virtualbox-configuration"
		}
		fmt.Printf("VirtualBox Network Type (bridged or host_only) [%s]: ", daemonConfig.Providers.Virtualbox[0].VirtualboxAdapterType)
		adapterType, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		adapterType = strings.TrimSuffix(adapterType, "\n")
		if adapterType != "" {
			daemonConfig.Providers.Virtualbox[0].VirtualboxAdapterType = config.VirtualboxAdapterType(adapterType)
		}
		fmt.Printf("Name of network adapter to attach to VirtualBox instances [%s]: ", daemonConfig.Providers.Virtualbox[0].AdapterName)
		adapterName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		adapterName = strings.TrimSuffix(adapterName, "\n")
		if adapterName != "" {
			daemonConfig.Providers.Virtualbox[0].AdapterName = adapterName
		}
	}
	return nil
}

func doVsphereConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with vSphere? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Vsphere) < 1 {
			daemonConfig.Providers.Vsphere = append(daemonConfig.Providers.Vsphere, config.Vsphere{})
		}
		if daemonConfig.Providers.Vsphere[0].Name == "" {
			daemonConfig.Providers.Vsphere[0].Name = "Vsphere-configuration"
		}
		fmt.Printf("vSphere username for authentication [%s]: ", daemonConfig.Providers.Vsphere[0].VsphereUser)
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSuffix(username, "\n")
		if username != "" {
			daemonConfig.Providers.Vsphere[0].VsphereUser = username
		}
		fmt.Printf("vSphere password for authentication [%s]: ", daemonConfig.Providers.Vsphere[0].VspherePassword)
		password, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		password = strings.TrimSuffix(password, "\n")
		if password != "" {
			daemonConfig.Providers.Vsphere[0].VspherePassword = password
		}
		fmt.Printf("vSphere authentication URL [%s]: ", daemonConfig.Providers.Vsphere[0].VsphereURL)
		authUrl, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		authUrl = strings.TrimSuffix(authUrl, "\n")
		if authUrl != "" {
			daemonConfig.Providers.Vsphere[0].VsphereURL = authUrl
		}
		fmt.Printf("vSphere datastore name [%s]: ", daemonConfig.Providers.Vsphere[0].Datastore)
		datastore, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		datastore = strings.TrimSuffix(datastore, "\n")
		if datastore != "" {
			daemonConfig.Providers.Vsphere[0].Datastore = datastore
		}
		fmt.Printf("vSphere datacenter name [%s]: ", daemonConfig.Providers.Vsphere[0].Datacenter)
		datacenter, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		datacenter = strings.TrimSuffix(datacenter, "\n")
		if datacenter != "" {
			daemonConfig.Providers.Vsphere[0].Datacenter = datacenter
		}
	}
	return nil
}

func doXenConfig(reader *bufio.Reader) error {
	fmt.Print("Do you wish to configure Bhojpur Kernel for use with Xen? [y/N]: ")
	y, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	y = strings.TrimSuffix(y, "\n")
	if y == "y" {
		if len(daemonConfig.Providers.Xen) < 1 {
			daemonConfig.Providers.Xen = append(daemonConfig.Providers.Xen, config.Xen{})
		}
		if daemonConfig.Providers.Xen[0].Name == "" {
			daemonConfig.Providers.Xen[0].Name = "Xen-configuration"
		}
		fmt.Printf("Name of xen bridge network interface to attach to Xen unikernels [%s]: ", daemonConfig.Providers.Xen[0].XenBridge)
		xenBridge, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		xenBridge = strings.TrimSuffix(xenBridge, "\n")
		if xenBridge != "" {
			daemonConfig.Providers.Xen[0].XenBridge = xenBridge
		}
		fmt.Printf("Path to PV Grub Boot Manager (see https://wiki.xen.org/wiki/PvGrub#Build for more info) [%s]: ", daemonConfig.Providers.Xen[0].KernelPath)
		pvKernel, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		pvKernel = strings.TrimSuffix(pvKernel, "\n")
		if pvKernel != "" {
			daemonConfig.Providers.Xen[0].KernelPath = pvKernel
		}
	}
	return nil
}
