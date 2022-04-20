package photon

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
	"github.com/vmware/photon-controller-go-sdk/photon"
)

func createVmdk(params types.StageImageParams, workVmdk func(file string) (string, error)) (string, int64, error) {

	localVmdkDir, err := ioutil.TempDir("", "vmdkdir.")
	if err != nil {
		return "", 0, errors.New("creating tmp file", err)
	}
	defer os.RemoveAll(localVmdkDir)
	localVmdkFile := filepath.Join(localVmdkDir, "boot.vmdk")

	logrus.WithField("raw-image", params.RawImage).Infof("creating boot volume from raw image")
	if err := common.ConvertRawToNewVmdk(params.RawImage.LocalImagePath, localVmdkFile); err != nil {
		return "", 0, errors.New("converting raw image to vmdk", err)
	}

	rawImageFile, err := os.Stat(localVmdkFile)
	if err != nil {
		return "", 0, errors.New("statting raw image file", err)
	}
	sizeMb := rawImageFile.Size() >> 20

	logrus.WithFields(logrus.Fields{
		"name": params.Name,
		"id":   params.Name,
		"size": sizeMb,
	}).Infof("importing base vmdk for unikernel image")

	imgId, err := workVmdk(localVmdkFile)
	return imgId, sizeMb, err

}

func (p *PhotonProvider) Stage(params types.StageImageParams) (_ *types.Image, err error) {
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
					logrus.Warn(errors.New("failed removing previously existing image", err))
				}
			}
		}
	}

	// create vmdk
	imgId, sizeMb, err := createVmdk(params, func(vmdkFile string) (string, error) {
		options := &photon.ImageCreateOptions{
			ReplicationType: "EAGER",
		}
		task, err := p.client.Images.CreateFromFile(vmdkFile, options)
		if err != nil {
			return "", errors.New("error creating photon image", err)
		}

		task, err = p.waitForTaskSuccess(task)
		if err != nil {
			return "", errors.New("error waiting for task creating photon image", err)
		}

		return task.Entity.ID, nil
	})
	if err != nil {
		return nil, errors.New("importing base boot.vmdk to photon", err)
	}

	// upload images
	image := &types.Image{
		Id:             imgId,
		Name:           params.Name,
		StageSpec:      params.RawImage.StageSpec,
		RunSpec:        params.RawImage.RunSpec,
		SizeMb:         sizeMb,
		Infrastructure: types.Infrastructure_PHOTON,
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
