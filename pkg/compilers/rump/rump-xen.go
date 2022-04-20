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
)

func CreateImageXen(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error) {
	return createImageXen(kernel, args, mntPoints, bakedEnv, noCleanup, false)
}

func CreateImageXenAddStub(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error) {
	return createImageXen(kernel, args, mntPoints, bakedEnv, noCleanup, true)
}

func createImageXen(kernel, args string, mntPoints, bakedEnv []string, noCleanup, addStub bool) (*types.RawImage, error) {
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
	volIndex := 0
	// add root -> sda1 mapping
	res.RunSpec.DeviceMappings = append(res.RunSpec.DeviceMappings, types.DeviceMapping{MountPoint: "/", DeviceName: "/dev/sda1"})

	bootBlk := blk{
		Source:     "etfs",
		Path:       "sda1",
		FSType:     "blk",
		MountPoint: "/bootpart",
	}
	c.Blk = append(c.Blk, bootBlk)

	for _, mntPoint := range mntPoints {
		// start from sdb; sda is for root.
		volIndex++
		deviceMapped := fmt.Sprintf("sd%c1", 'a'+volIndex)
		blk := blk{
			Source:     "etfs",
			Path:       deviceMapped,
			FSType:     "blk",
			MountPoint: mntPoint,
		}

		c.Blk = append(c.Blk, blk)
		res.RunSpec.DeviceMappings = append(res.RunSpec.DeviceMappings, types.DeviceMapping{MountPoint: mntPoint, DeviceName: "/dev/" + deviceMapped})
	}

	// aws network
	c.Net = &net{
		If:     "xenif0",
		Cloner: "true",
		Type:   "inet",
		Method: DHCP,
	}

	cmdline, err := toRumpJson(c)
	if err != nil {
		return nil, err
	}
	imgFile, err := compilers.BuildBootableImage(kernel, cmdline, false, noCleanup)

	if err != nil {
		return nil, err
	}

	res.LocalImagePath = imgFile
	res.StageSpec = types.StageSpec{
		ImageFormat:           types.ImageFormat_RAW,
		XenVirtualizationType: types.XenVirtualizationType_Paravirtual,
	}
	res.RunSpec.DefaultInstanceMemory = 1024 // todo: investigate why rump uses so much memory on xen

	return res, nil
}
