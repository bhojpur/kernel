package qemu

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
	"os/exec"
	"strings"
	"time"

	"io/ioutil"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *QemuProvider) RunInstance(params types.RunInstanceParams) (_ *types.Instance, err error) {
	logrus.WithFields(logrus.Fields{
		"image-id": params.ImageId,
		"mounts":   params.MntPointsToVolumeIds,
		"env":      params.Env,
	}).Infof("running instance %s", params.Name)

	if _, err := p.GetInstance(params.Name); err == nil {
		return nil, errors.New("instance with name "+params.Name+" already exists. qemu provider requires unique names for instances", nil)
	}

	image, err := p.GetImage(params.ImageId)
	if err != nil {
		return nil, errors.New("getting image", err)
	}

	if err := common.VerifyMntsInput(p, image, params.MntPointsToVolumeIds); err != nil {
		return nil, errors.New("invalid mapping for volume", err)
	}

	volumeIdInOrder := make([]string, len(params.MntPointsToVolumeIds))

	for mntPoint, volumeId := range params.MntPointsToVolumeIds {

		controllerPort, err := common.GetControllerPortForMnt(image, mntPoint)
		if err != nil {
			return nil, err
		}
		volumeIdInOrder[controllerPort] = volumeId
	}

	logrus.Debugf("creating qemu vm")

	volImagesInOrder, err := p.getVolumeImages(volumeIdInOrder)
	if err != nil {
		return nil, errors.New("can't get volumes", err)
	}

	volArgs := volPathToQemuArgs(volImagesInOrder)

	if params.InstanceMemory == 0 {
		params.InstanceMemory = image.RunSpec.DefaultInstanceMemory
	}

	qemuArgs := []string{"-m", fmt.Sprintf("%v", params.InstanceMemory), "-net",
		"nic,model=virtio,netdev=mynet0", "-netdev", "user,id=mynet0,net=192.168.76.0/24,dhcpstart=192.168.76.9",
	}

	cmdlinedata, err := ioutil.ReadFile(getCmdlinePath(image.Name))
	if err != nil {
		logrus.Debugf("cmdLine not found, assuming classic bootloader")
		qemuArgs = append(qemuArgs, "-drive", fmt.Sprintf("file=%s,format=raw,if=ide", getImagePath(image.Name)))
	} else {
		// inject env for rump:
		cmdline := string(cmdlinedata)
		if compilers.CompilerType(image.RunSpec.Compiler).Base() == compilers.Rump {
			cmdline = injectEnv(cmdline, params.Env)
		}

		// qemu escape
		cmdline = strings.Replace(cmdline, ",", ",,", -1)

		if _, err := os.Stat(getImagePath(image.Name)); err == nil {
			qemuArgs = append(qemuArgs, "-device", "virtio-blk-pci,id=blk0,drive=hd0")
			qemuArgs = append(qemuArgs, "-drive", fmt.Sprintf("file=%s,format=qcow2,if=none,id=hd0", getImagePath(image.Name)))
		}

		qemuArgs = append(qemuArgs, "-kernel", getKernelPath(image.Name))
		qemuArgs = append(qemuArgs, "-append", cmdline)
	}

	if params.DebugMode {
		logrus.Debugf("running instance in debug mode.\nattach Bhojpur Kernel debugger to port :%v", p.config.DebuggerPort)
		qemuArgs = append(qemuArgs, "-s", "-S")
		debuggerTargetImageName = image.Name
	}

	if p.config.NoGraphic {
		qemuArgs = append(qemuArgs, "-nographic", "-vga", "none")
	}

	qemuArgs = append(qemuArgs, volArgs...)
	cmd := exec.Command("qemu-system-x86_64", qemuArgs...)

	util.LogCommand(cmd, true)

	if err := cmd.Start(); err != nil {
		return nil, errors.New("can't start qemu - make sure it's in your path.", nil)
	}

	var instanceIp string

	instance := &types.Instance{
		Id:             fmt.Sprintf("%d", cmd.Process.Pid),
		Name:           params.Name,
		State:          types.InstanceState_Running,
		IpAddress:      instanceIp,
		Infrastructure: types.Infrastructure_QEMU,
		ImageId:        image.Id,
		Created:        time.Now(),
	}

	if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
		instances[instance.Id] = instance
		return nil
	}); err != nil {
		return nil, errors.New("modifying instance map in state", err)
	}

	logrus.WithField("instance", instance).Infof("instance created successfully")

	return instance, nil
}

func (p *QemuProvider) getVolumeImages(volumeIdInOrder []string) ([]string, error) {

	var volPath []string
	for _, v := range volumeIdInOrder {
		v, err := p.GetVolume(v)
		if err != nil {
			return nil, err
		}
		volPath = append(volPath, getVolumePath(v.Name))
	}
	return volPath, nil
}

func volPathToQemuArgs(volPaths []string) []string {
	var res []string
	for _, v := range volPaths {
		res = append(res, "-drive", fmt.Sprintf("if=virtio,file=%s,format=qcow2", v))
	}
	return res
}

func injectEnv(cmdline string, env map[string]string) string {
	// rump json is not really json so we can't parse it
	var envRumpJson []string
	for key, value := range env {
		envRumpJson = append(envRumpJson, fmt.Sprintf("\"env\": \"%s=%s\"", key, value))
	}

	cmdline = cmdline[:len(cmdline)-2] + "," + strings.Join(envRumpJson, ",") + "}}"
	return cmdline
}
