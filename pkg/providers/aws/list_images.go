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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

const KERNEL_IMAGE_ID = "KERNEL_IMAGE_ID"

func (p *AwsProvider) ListImages() ([]*types.Image, error) {
	if len(p.state.GetImages()) < 1 {
		return []*types.Image{}, nil
	}
	imageIds := []*string{}
	for imageId := range p.state.GetImages() {
		imageIds = append(imageIds, aws.String(imageId))
	}
	param := &ec2.DescribeImagesInput{
		ImageIds: imageIds,
	}
	output, err := p.newEC2().DescribeImages(param)
	if err != nil {
		return nil, errors.New("running ec2 describe images ", err)
	}
	images := []*types.Image{}
	for _, ec2Image := range output.Images {
		imageId := *ec2Image.ImageId
		image, ok := p.state.GetImages()[imageId]
		if !ok {
			logrus.WithFields(logrus.Fields{"ec2Image": ec2Image}).Errorf("found an image that Bhojpur Kernel has no record of")
			continue
		}
		images = append(images, image)
	}
	return images, nil
}
