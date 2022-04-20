package common

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
	"encoding/json"
	"os"
	"path/filepath"

	"io/ioutil"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func ConvertRawImage(sourceFormat, targetFormat types.ImageFormat, inputFile, outputFile string) error {
	targetFormatName := string(targetFormat)
	if targetFormat == types.ImageFormat_VHD {
		targetFormatName = "vpc" //for some reason qemu calls VHD disks vpc
	}
	dir := filepath.Dir(inputFile)
	outDir := filepath.Dir(outputFile)

	container := kutil.NewContainer("qemu-util").WithVolume(dir, "/bhojpur/input").
		WithVolume(outDir, "/bhojpur/output")

	args := []string{"qemu-img", "convert", "-f", string(sourceFormat), "-O", targetFormatName}
	if targetFormat == types.ImageFormat_VMDK {
		args = append(args, "-o", "compat6")
	}

	//this needs to be done because docker produces files as root. argh!!!
	tmpOutputFile, err := ioutil.TempFile(outDir, "convert.image.result.")
	if err != nil {
		return errors.New("temp file for root user", err)
	}
	tmpOutputFile.Close()
	defer os.Remove(tmpOutputFile.Name())

	args = append(args, "/bhojpur/input/"+filepath.Base(inputFile), "/bhojpur/output/"+filepath.Base(tmpOutputFile.Name()))

	logrus.WithField("command", args).Debugf("running command")
	if err := container.Run(args...); err != nil {
		return errors.New("failed converting raw image to "+string(targetFormat), err)
	}

	if err := kos.CopyFile(tmpOutputFile.Name(), outputFile); err != nil {
		return errors.New("copying tmp result to final result", err)
	}

	return nil
}

func fixVmdk(vmdkFile string) error {
	file, err := os.OpenFile(vmdkFile, os.O_RDWR, 0)
	if err != nil {
		return errors.New("can't open vmdk", err)
	}
	defer file.Close()

	var buffer [1024]byte

	n, err := file.Read(buffer[:])
	if err != nil {
		return errors.New("can't read vmdk", err)
	}
	if n < len(buffer) {
		return errors.New("bad vmdk", err)
	}

	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		return errors.New("can't seek vmdk", err)
	}

	result := bytes.Replace(buffer[:], []byte("# The disk Data Base"), []byte("# The Disk Data Base"), 1)

	_, err = file.Write(result)
	if err != nil {
		return errors.New("can't write vmdk", err)
	}

	return nil
}

func ConvertRawToNewVmdk(inputFile, outputFile string) error {

	dir := filepath.Dir(inputFile)
	outDir := filepath.Dir(outputFile)

	container := kutil.NewContainer("euranova/ubuntu-vbox").WithVolume(dir, dir).
		WithVolume(outDir, outDir)

	args := []string{
		"VBoxManage", "convertfromraw", inputFile, outputFile, "--format", "vmdk", "--variant", "Stream"}

	logrus.WithField("command", args).Debugf("running command")
	if err := container.Run(args...); err != nil {
		return errors.New("failed converting raw image to vmdk", err)
	}

	err := fixVmdk(outputFile)
	if err != nil {
		return errors.New("failed converting raw image to vmdk", err)
	}

	return nil
}

func GetVirtualImageSize(imageFile string, imageFormat types.ImageFormat) (int64, error) {
	formatName := string(imageFormat)
	if imageFormat == types.ImageFormat_VHD {
		formatName = "vpc" //for some reason qemu calls VHD disks vpc
	}
	dir := filepath.Dir(imageFile)

	container := kutil.NewContainer("qemu-util").WithVolume(dir, dir)
	args := []string{"qemu-img", "info", "--output", "json", "-f", formatName, imageFile}

	logrus.WithField("command", args).Debugf("running command")
	out, err := container.CombinedOutput(args...)
	if err != nil {
		return -1, errors.New("failed getting image info", err)
	}
	var info imageInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return -1, errors.New("parsing "+string(out)+" to json", err)
	}
	return info.VirtualSize, nil
}

type imageInfo struct {
	VirtualSize int64  `json:"virtual-size"`
	Filename    string `json:"filename"`
	ClusterSize int    `json:"cluster-size"`
	Format      string `json:"format"`
	ActualSize  int    `json:"actual-size"`
	DirtyFlag   bool   `json:"dirty-flag"`
}
