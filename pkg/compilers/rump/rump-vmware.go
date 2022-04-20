package rump

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
	"strings"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/sirupsen/logrus"
)

func CreateImageVmware(kernel string, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error) {
	return createImageVmware(kernel, args, mntPoints, bakedEnv, noCleanup, false)
}

func CreateImageVmwareAddStub(kernel string, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error) {
	return createImageVmware(kernel, args, mntPoints, bakedEnv, noCleanup, true)
}

func createImageVmware(kernel string, args string, mntPoints, bakedEnv []string, noCleanup, addStub bool) (*types.RawImage, error) {
	// create rump config
	var c rumpConfig
	if bakedEnv != nil {
		c.Env = make(map[string]string)
		for i, pair := range bakedEnv {
			c.Env[fmt.Sprintf("env%d", i)] = pair
		}
	}

	argv := []string{}
	if args != "" {
		argv = strings.Split(args, " ")
	}
	c = setRumpCmdLine(c, "program.bin", argv, addStub)

	res := &types.RawImage{}
	// add root -> sd0 mapping
	res.RunSpec.DeviceMappings = append(res.RunSpec.DeviceMappings, types.DeviceMapping{MountPoint: "/", DeviceName: "sd0"})

	bootBlk := blk{
		Source:     "dev",
		Path:       "/dev/sd0e", // no disk label on the boot partition; so partition e is used.
		FSType:     "blk",
		MountPoint: "/bootpart",
	}
	c.Blk = append(c.Blk, bootBlk)

	for i, mntPoint := range mntPoints {
		deviceMapped := fmt.Sprintf("sd%ca", '1'+i)
		blk := blk{
			Source:     "dev",
			Path:       "/dev/" + deviceMapped,
			FSType:     "blk",
			MountPoint: mntPoint,
		}

		c.Blk = append(c.Blk, blk)
		res.RunSpec.DeviceMappings = append(res.RunSpec.DeviceMappings, types.DeviceMapping{MountPoint: mntPoint, DeviceName: deviceMapped})
	}

	// aws network
	c.Net = &net{
		If:     "wm0",
		Type:   "inet",
		Method: DHCP,
	}

	cmdline, err := toRumpJson(c)
	if err != nil {
		return nil, err
	}

	imgFile, err := compilers.BuildBootableImage(kernel, cmdline, true, noCleanup)
	if err != nil {
		return nil, err
	}

	res.LocalImagePath = imgFile
	res.StageSpec.ImageFormat = types.ImageFormat_RAW
	res.RunSpec.StorageDriver = types.StorageDriver_SCSI
	res.RunSpec.VsphereNetworkType = types.VsphereNetworkType_E1000
	res.RunSpec.DefaultInstanceMemory = 256
	logrus.WithField("runspec", res.RunSpec).Infof("created raw vmware image")
	return res, nil
}
