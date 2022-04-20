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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *QemuProvider) PullImage(params types.PullImagePararms) error {
	images, err := p.ListImages()
	if err != nil {
		return errors.New("retrieving image list for existing image", err)
	}
	for _, image := range images {
		if image.Name == params.ImageName {
			if !params.Force {
				return errors.New("an image already exists with name '"+params.ImageName+"', try again with --force", nil)
			} else {
				logrus.WithField("image", image).Warnf("force: deleting previous image with name " + params.ImageName)
				if err := p.DeleteImage(image.Id, true); err != nil {
					logrus.Warn(errors.New("failed removing previously existing image", err))
				}
			}
		}
	}

	tmpImage, err := ioutil.TempFile("", "tmp-pull-image-"+params.ImageName)
	if err != nil {
		return errors.New("creating tmp file", err)
	}
	defer os.RemoveAll(tmpImage.Name())
	image, err := common.PullImage(params.Config, params.ImageName, tmpImage)
	if err != nil {
		return errors.New("pulling image", err)
	}
	imagePath := getImagePath(image.Name)
	os.MkdirAll(filepath.Dir(imagePath), 0755)
	if err := os.Rename(tmpImage.Name(), imagePath); err != nil {
		return errors.New("renaming tmp image to "+imagePath, err)
	}

	if err := p.state.ModifyImages(func(images map[string]*types.Image) error {
		images[image.Name] = image
		return nil
	}); err != nil {
		return errors.New("modifying image map in state", err)
	}
	logrus.Infof("image %v pulled successfully from %v", err)
	return nil
}
