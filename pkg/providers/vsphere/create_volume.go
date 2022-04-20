package vsphere

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
	"time"

	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *VsphereProvider) CreateVolume(params types.CreateVolumeParams) (_ *types.Volume, err error) {
	if _, volumeErr := p.GetImage(params.Name); volumeErr == nil {
		return nil, errors.New("volume already exists", nil)
	}
	c := p.getClient()

	localVmdkDir, err := ioutil.TempDir("", "localvmdkdir.")
	if err != nil {
		return nil, errors.New("creating tmp file", err)
	}
	defer os.RemoveAll(localVmdkDir)
	localVmdkFile := filepath.Join(localVmdkDir, "data.vmdk")
	logrus.WithField("raw-image", params.ImagePath).Infof("creating vmdk from raw image")
	if err := common.ConvertRawImage(types.ImageFormat_RAW, types.ImageFormat_VMDK, params.ImagePath, localVmdkFile); err != nil {
		return nil, errors.New("converting raw image to vmdk", err)
	}

	rawImageFile, err := os.Stat(localVmdkFile)
	if err != nil {
		return nil, errors.New("statting raw image file", err)
	}
	sizeMb := rawImageFile.Size() >> 20

	vsphereVolumeDir := getVolumeDatastoreDir(params.Name)
	if err := c.Mkdir(vsphereVolumeDir); err != nil {
		return nil, errors.New("creating vsphere directory for volume", err)
	}
	defer func() {
		if err != nil {
			if params.NoCleanup {
				logrus.Warnf("because --no-cleanup flag was provided, not cleaning up failed volume %s at %s", params.Name, vsphereVolumeDir)
				return
			}
			logrus.WithError(err).Warnf("creating volume failed, cleaning up volume on datastore")
			c.Rmdir(vsphereVolumeDir)
		}
	}()

	if err := c.ImportVmdk(localVmdkFile, vsphereVolumeDir); err != nil {
		return nil, errors.New("importing data.vmdk to vsphere datastore", err)
	}

	volume := &types.Volume{
		Id:             params.Name,
		Name:           params.Name,
		SizeMb:         sizeMb,
		Attachment:     "",
		Infrastructure: types.Infrastructure_VSPHERE,
		Created:        time.Now(),
	}

	if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volumes[volume.Id] = volume
		return nil
	}); err != nil {
		return nil, errors.New("modifying volume map in state", err)
	}
	return volume, nil
}
