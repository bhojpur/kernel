package openstack

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

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func (p *OpenstackProvider) DeleteImage(id string, force bool) error {
	image, err := p.GetImage(id)
	if err != nil {
		return errors.New("retrieving image", err)
	}

	// Delete instances of this image.
	instances, err := p.ListInstances()
	if err != nil {
		return errors.New("failed to retrieve list of instances", err)
	}
	for _, instance := range instances {
		if instance.ImageId == image.Id {
			if !force {
				return fmt.Errorf("instance '%s' found which uses image '%s'! Try again with --force.", instance.Id, image.Id)
			} else {
				err = p.DeleteInstance(instance.Id, true)
				if err != nil {
					return errors.New(fmt.Sprintf("failed to delete instance '%s' which uses image '%s'", instance.Id, image.Id), err)
				}
			}
		}
	}

	clientGlance, err := p.newClientGlance()
	if err != nil {
		return err
	}

	if err := deleteImage(clientGlance, image.Id); err != nil {
		return errors.New(fmt.Sprintf("failed to delete image '%s'", image.Id), err)
	}

	// Update state.
	if err := p.state.ModifyImages(func(imageList map[string]*types.Image) error {
		delete(imageList, image.Id)
		return nil
	}); err != nil {
		return errors.New("failed to modify image map in state", err)
	}
	return nil
}

// deleteImage deletes image from OpenStack.
func deleteImage(clientGlance *gophercloud.ServiceClient, imageId string) error {
	return images.Delete(clientGlance, imageId).Err
}
