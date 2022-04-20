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

	"io/ioutil"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/sirupsen/logrus"
)

func CreateImageQemu(kernel string, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error) {
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
	c = setRumpCmdLine(c, "program.bin", argv, false)

	bootBlk := blk{
		Source:     "dev",
		Path:       "/dev/ld0e",
		FSType:     "blk",
		MountPoint: "/bootpart",
	}
	c.Blk = append(c.Blk, bootBlk)

	res := &types.RawImage{}
	res.RunSpec.Compiler = compilers.Rump

	for i, mntPoint := range mntPoints {
		deviceMapped := fmt.Sprintf("ld%ca", '1'+i)
		blk := blk{
			Source:     "dev",
			Path:       "/dev/" + deviceMapped,
			FSType:     "blk",
			MountPoint: mntPoint,
		}

		c.Blk = append(c.Blk, blk)
		logrus.Debugf("adding mount point to image: %s:%s", mntPoint, deviceMapped)
		res.RunSpec.DeviceMappings = append(res.RunSpec.DeviceMappings,
			types.DeviceMapping{MountPoint: mntPoint, DeviceName: deviceMapped})
	}

	// virtualbox network
	c.Net = &net{
		If:     "vioif0",
		Type:   "inet",
		Method: DHCP,
	}

	cmdline, err := toRumpJson(c)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("writing rump json config: %s", cmdline)

	imgFile, err := compilers.BuildBootableImage(kernel, cmdline, true, noCleanup)
	if err != nil {
		return nil, err
	}

	//copy kernel for qemu
	if err := kos.CopyFile(kernel, filepath.Join(filepath.Dir(imgFile), "program.bin")); err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filepath.Join(filepath.Dir(imgFile), "cmdline"), []byte(cmdline), 0644); err != nil {
		return nil, err
	}

	res.LocalImagePath = imgFile
	res.StageSpec.ImageFormat = types.ImageFormat_RAW
	res.RunSpec.DefaultInstanceMemory = 512
	return res, nil

}
