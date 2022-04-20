package osv

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
	"io/ioutil"
	"os"
	"path/filepath"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const OSV_DEFAULT_SIZE = "1GB"

type dynamicProjectConfig struct {
	// Size is a string representing logical image size e.g. "10GB"
	Size string `yaml:"image_size"`
}

// CreateImageDynamic creates OSv image from project source directory and returns filepath of it.
func CreateImageDynamic(params types.CompileImageParams, useEc2Bootstrap bool) (string, error) {
	container := kutil.NewContainer("compilers-osv-dynamic").
		WithVolume(params.SourcesDir+"/", "/project_directory").
		WithEnv("MAX_IMAGE_SIZE", fmt.Sprintf("%dMB", params.SizeMB))

	logrus.WithFields(logrus.Fields{
		"params": params,
	}).Debugf("running compilers-osv-dynamic container")

	if err := container.Run(); err != nil {
		return "", errors.New("failed running compilers-osv-dynamic on "+params.SourcesDir, err)
	}

	resultFile, err := ioutil.TempFile("", "osv-dynamic.qemu.")
	if err != nil {
		return "", errors.New("failed to create tmpfile for result", err)
	}
	defer func() {
		if err != nil && !params.NoCleanup {
			os.Remove(resultFile.Name())
		}
	}()

	if err := os.Rename(filepath.Join(params.SourcesDir, "boot.qcow2"), resultFile.Name()); err != nil {
		return "", errors.New("failed to rename result file", err)
	}
	return resultFile.Name(), nil
}

// readImageSizeFromManifest parses manifest.yaml and returns image size.
// Returns default image size if anything goes wrong.
func readImageSizeFromManifest(projectDir string) kos.MegaBytes {
	config := dynamicProjectConfig{
		Size: OSV_DEFAULT_SIZE,
	}
	defaultMB, _ := kos.ParseSize(OSV_DEFAULT_SIZE)

	data, err := ioutil.ReadFile(filepath.Join(projectDir, "manifest.yaml"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err,
			"defaultSize": OSV_DEFAULT_SIZE,
		}).Warning("could not find manifest.yaml. Fallback to using default unikernel size.")
		return defaultMB
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err,
			"defaultSize": OSV_DEFAULT_SIZE,
		}).Warning("failed to parse manifest.yaml. Fallback to using default unikernel size.")
		return defaultMB
	}

	sizeMB, err := kos.ParseSize(config.Size)
	if err != nil {
		return defaultMB
	}

	return sizeMB
}
