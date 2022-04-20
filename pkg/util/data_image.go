//go:build !container-binary
// +build !container-binary

package util

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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func BuildRawDataImageWithType(dataTar io.ReadCloser, size kos.MegaBytes, volType string, usePartitionTables bool) (string, error) {
	buildDir, err := ioutil.TempDir("", ".raw_data_image_folder.")
	if err != nil {
		return "", errors.New("creating tmp build folder", err)
	}
	defer os.RemoveAll(buildDir)

	dataFolder := filepath.Join(buildDir, "data")
	err = os.Mkdir(dataFolder, 0755)
	if err != nil {
		return "", errors.New("creating tmp data folder", err)
	}

	if err := kos.ExtractTar(dataTar, dataFolder); err != nil {
		return "", errors.New("extracting data tar", err)
	}

	container := NewContainer("image-creator").Privileged(true).WithVolume("/dev/", "/dev/").
		WithVolume(buildDir+"/", "/opt/vol")

	tmpResultFile, err := ioutil.TempFile(buildDir, "data.image.result.img.")
	if err != nil {
		return "", err
	}
	tmpResultFile.Close()
	args := []string{"-o", filepath.Base(tmpResultFile.Name())}

	if size > 0 {
		args = append(args, "-p", fmt.Sprintf("%v", usePartitionTables),
			"-v", fmt.Sprintf("%s,%v", filepath.Base(dataFolder), size.ToBytes()))
	} else {
		args = append(args, "-p", fmt.Sprintf("%v", usePartitionTables),
			"-v", filepath.Base(dataFolder),
		)
	}
	args = append(args, "-t", volType)

	logrus.WithFields(logrus.Fields{
		"command": args,
	}).Debugf("running image-creator container")

	if err = container.Run(args...); err != nil {
		return "", errors.New("failed running image-creator on "+dataFolder, err)
	}

	resultFile, err := ioutil.TempFile("", "data-volume-creator-result.img.")
	if err != nil {
		return "", err
	}
	resultFile.Close()
	if err := os.Rename(tmpResultFile.Name(), resultFile.Name()); err != nil {
		return "", errors.New("renaming "+tmpResultFile.Name()+" to "+resultFile.Name(), err)
	}

	return resultFile.Name(), nil
}

func BuildRawDataImage(dataTar io.ReadCloser, size kos.MegaBytes, usePartitionTables bool) (string, error) {
	return BuildRawDataImageWithType(dataTar, size, "ext2", usePartitionTables)
}
func BuildEmptyDataVolumeWithType(size kos.MegaBytes, volType string) (string, error) {

	if size < 1 {
		return "", errors.New("must specify size > 0", nil)
	}
	dataFolder, err := ioutil.TempDir("", "empty.data.folder.")
	if err != nil {
		return "", errors.New("creating tmp build folder", err)
	}
	defer os.RemoveAll(dataFolder)

	buildDir := filepath.Dir(dataFolder)

	container := NewContainer("image-creator").Privileged(true).WithVolume("/dev/", "/dev/").
		WithVolume(buildDir+"/", "/opt/vol")

	tmpResultFile, err := ioutil.TempFile(buildDir, "data.image.result.img.")
	if err != nil {
		return "", err
	}
	tmpResultFile.Close()
	args := []string{"-v", fmt.Sprintf("%s,%v", filepath.Base(dataFolder), size.ToBytes()), "-o", filepath.Base(tmpResultFile.Name())}
	args = append(args, "-t", volType)

	logrus.WithFields(logrus.Fields{
		"command": args,
	}).Debugf("running image-creator container")
	if err := container.Run(args...); err != nil {
		return "", errors.New("failed running image-creator on "+dataFolder, err)
	}

	resultFile, err := ioutil.TempFile("", "empty-data-volume-creator-result.img.")
	if err != nil {
		return "", err
	}
	resultFile.Close()
	if err := os.Rename(tmpResultFile.Name(), resultFile.Name()); err != nil {
		return "", errors.New("renaming "+tmpResultFile.Name()+" to "+resultFile.Name(), err)
	}

	return resultFile.Name(), nil
}

func BuildEmptyDataVolume(size kos.MegaBytes) (string, error) {
	return BuildEmptyDataVolumeWithType(size, "ext2")
}
