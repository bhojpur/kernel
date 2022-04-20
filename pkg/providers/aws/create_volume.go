package aws

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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *AwsProvider) CreateVolume(params types.CreateVolumeParams) (*types.Volume, error) {
	logrus.WithField("raw-image", params.ImagePath).WithField("az", p.config.Zone).Infof("creating data volume from raw image")
	s3svc := p.newS3()
	ec2svc := p.newEC2()
	imageFile, err := os.Stat(params.ImagePath)
	if err != nil {
		return nil, errors.New("stat image file", err)
	}
	volumeId, err := createDataVolumeFromRawImage(s3svc, ec2svc, params.ImagePath, imageFile.Size(), types.ImageFormat_RAW, p.config.Zone)
	if err != nil {
		return nil, errors.New("creating aws boot volume", err)
	}
	tagVolumeInput := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(volumeId),
		},
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(params.Name),
			},
		},
	}
	if _, err := ec2svc.CreateTags(tagVolumeInput); err != nil {
		return nil, errors.New("tagging volume", err)
	}

	rawImageFile, err := os.Stat(params.ImagePath)
	if err != nil {
		return nil, errors.New("statting raw image file", err)
	}
	sizeMb := rawImageFile.Size() >> 20

	volume := &types.Volume{
		Id:             volumeId,
		Name:           params.Name,
		SizeMb:         sizeMb,
		Attachment:     "",
		Infrastructure: types.Infrastructure_AWS,
		Created:        time.Now(),
	}

	if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		volumes[volume.Id] = volume
		return nil
	}); err != nil {
		return nil, errors.New("modifying volume map in state", err)
	}

	return nil, nil
}
func (p *AwsProvider) CreateEmptyVolume(name string, size int) (*types.Volume, error) {
	return nil, nil
}
