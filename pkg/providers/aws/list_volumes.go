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

func (p *AwsProvider) ListVolumes() ([]*types.Volume, error) {
	if len(p.state.GetVolumes()) < 1 {
		return []*types.Volume{}, nil
	}
	volumeIds := []*string{}
	for volumeId := range p.state.GetVolumes() {
		volumeIds = append(volumeIds, aws.String(volumeId))
	}
	param := &ec2.DescribeVolumesInput{
		VolumeIds: volumeIds,
	}
	output, err := p.newEC2().DescribeVolumes(param)
	if err != nil {
		return nil, errors.New("running ec2 describe volumes ", err)
	}
	volumes := []*types.Volume{}
	for _, ec2Volume := range output.Volumes {
		volumeId := *ec2Volume.VolumeId
		if volumeId == "" {
			continue
		}
		volume, ok := p.state.GetVolumes()[volumeId]
		if !ok {
			logrus.WithFields(logrus.Fields{"ec2Volume": ec2Volume}).Errorf("found a volume that Bhojpur Kernel has no record of")
			continue
		}
		if len(ec2Volume.Attachments) > 0 {
			if len(ec2Volume.Attachments) > 1 {
				return nil, errors.New("ec2 reports volume to have >1 attachments. wut", nil)
			}
			volume.Attachment = *ec2Volume.Attachments[0].InstanceId
		} else {
			volume.Attachment = ""
		}
		if err := p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
			volumes[volume.Id] = volume
			return nil
		}); err != nil {
			return nil, errors.New("modifying volume map in state", err)
		}
		volumes = append(volumes, volume)
	}
	return volumes, nil
}
