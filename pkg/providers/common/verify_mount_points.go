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
	"github.com/bhojpur/kernel/pkg/providers"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func VerifyMntsInput(p providers.Provider, image *types.Image, mntPointsToVolumeIds map[string]string) error {
	for _, deviceMapping := range image.RunSpec.DeviceMappings {
		if deviceMapping.MountPoint == "/" {
			//ignore boot mount point
			continue
		}
		_, ok := mntPointsToVolumeIds[deviceMapping.MountPoint]
		if !ok {
			logrus.WithFields(logrus.Fields{"required-device-mappings": image.RunSpec.DeviceMappings}).Errorf("requied mount point missing: %s", deviceMapping.MountPoint)
			return errors.New("required mount point missing from input", nil)
		}
	}
	for mntPoint, volumeId := range mntPointsToVolumeIds {
		mntPointExists := false
		for _, deviceMapping := range image.RunSpec.DeviceMappings {
			if deviceMapping.MountPoint == mntPoint {
				mntPointExists = true
				break
			}
		}
		if !mntPointExists {
			return errors.New("mount point "+mntPoint+" does not exist for image "+image.Id, nil)
		}
		_, err := p.GetVolume(volumeId)
		if err != nil {
			return errors.New("could not find volume "+volumeId, err)
		}
	}
	return nil
}
