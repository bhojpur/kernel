package xen

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

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *XenProvider) CreateVolume(params types.CreateVolumeParams) (_ *types.Volume, err error) {
	if _, volumeErr := p.GetImage(params.Name); volumeErr == nil {
		return nil, errors.New("volume already exists", nil)
	}

	volumePath := getVolumePath(params.Name)
	if err := os.MkdirAll(filepath.Dir(volumePath), 0755); err != nil {
		return nil, errors.New("creating directory for volume file", err)
	}
	defer func() {
		if err != nil {
			if params.NoCleanup {
				logrus.Warnf("because --no-cleanup flag was provided, not cleaning up failed volume %s at %s", params.Name, volumePath)
			} else {
				os.RemoveAll(filepath.Dir(volumePath))
			}
		}
	}()
	logrus.WithField("raw-image", params.ImagePath).Infof("creating volume from raw image")

	rawImageFile, err := os.Stat(params.ImagePath)
	if err != nil {
		return nil, errors.New("statting raw image file", err)
	}
	sizeMb := rawImageFile.Size() >> 20

	if err := os.Rename(params.ImagePath, volumePath); err != nil {
		return nil, errors.New("copying raw image from "+params.ImagePath+"to "+volumePath, err)
	}

	volume := &types.Volume{
		Id:             params.Name,
		Name:           params.Name,
		SizeMb:         sizeMb,
		Attachment:     "",
		Infrastructure: types.Infrastructure_XEN,
		Created:        time.Now(),
	}

	err = p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volumes[volume.Id] = volume
		return nil
	})
	if err != nil {
		return nil, errors.New("modifying volume map in state", err)
	}
	return volume, nil

}
