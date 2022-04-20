package ukvm

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
	"os"
	"path/filepath"
	"time"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *UkvmProvider) Stage(params types.StageImageParams) (_ *types.Image, err error) {
	images, err := p.ListImages()
	if err != nil {
		return nil, errors.New("retrieving image list for existing image", err)
	}
	for _, image := range images {
		if image.Name == params.Name {
			if !params.Force {
				return nil, errors.New("an image already exists with name '"+params.Name+"', try again with --force", nil)
			} else {
				logrus.WithField("image", image).Warnf("force: deleting previous image with name " + params.Name)
				if err := p.DeleteImage(image.Id, true); err != nil {
					logrus.Warn("failed to remove previously existing image", err)
				}
			}
		}
	}
	imageName := params.Name
	imageDir := getImageDir(imageName)
	logrus.Debugf("making directory: %s", imageDir)
	if err := os.MkdirAll(imageDir, 0777); err != nil {
		return nil, errors.New("creating directory for boot image", err)
	}
	defer func() {
		if err != nil && !params.NoCleanup {
			os.RemoveAll(imageDir)
		}
	}()

	kernelPath := filepath.Join(params.RawImage.LocalImagePath, "program.bin")
	if err := kos.CopyFile(kernelPath, getKernelPath(imageName)); err != nil {
		return nil, errors.New("program.bin cannot be copied", err)

	}
	ukvmPath := filepath.Join(params.RawImage.LocalImagePath, "ukvm-bin")
	if err := kos.CopyFile(ukvmPath, getUkvmPath(imageName)); err != nil {
		return nil, errors.New("ukvm-bin cannot be copied", err)
	}

	kernelPathInfo, err := os.Stat(kernelPath)
	if err != nil {
		return nil, errors.New("statting unikernel file", err)
	}
	ukvmPathInfo, err := os.Stat(ukvmPath)
	if err != nil {
		return nil, errors.New("statting ukvm file", err)
	}
	sizeMb := (ukvmPathInfo.Size() + kernelPathInfo.Size()) >> 20

	logrus.WithFields(logrus.Fields{
		"name": params.Name,
		"id":   params.Name,
		"size": sizeMb,
	}).Infof("copying raw boot image")

	image := &types.Image{
		Id:             params.Name,
		Name:           params.Name,
		RunSpec:        params.RawImage.RunSpec,
		StageSpec:      params.RawImage.StageSpec,
		SizeMb:         sizeMb,
		Infrastructure: types.Infrastructure_UKVM,
		Created:        time.Now(),
	}

	if err := p.state.ModifyImages(func(images map[string]*types.Image) error {
		images[params.Name] = image
		return nil
	}); err != nil {
		return nil, errors.New("modifying image map in state", err)
	}

	logrus.WithFields(logrus.Fields{"image": image}).Infof("image created succesfully")
	return image, nil
}
