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
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/rackspace/gophercloud/pagination"
)

func (p *OpenstackProvider) ListImages() ([]*types.Image, error) {
	// Return immediately if no image is managed by Bhojpur Kernel.
	managedImages := p.state.GetImages()
	if len(managedImages) < 1 {
		return []*types.Image{}, nil
	}

	clientGlance, err := p.newClientGlance()
	if err != nil {
		return nil, err
	}

	return fetchImages(clientGlance, managedImages)
}

func fetchImages(clientGlance *gophercloud.ServiceClient, managedImages map[string]*types.Image) ([]*types.Image, error) {
	result := []*types.Image{}

	pager := images.List(clientGlance, nil)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}

		for _, i := range imageList {
			// Filter out images that Bhojpur Kernel is not aware of.
			image, ok := managedImages[i.ID]
			if !ok {
				continue
			}
			result = append(result, image)
		}

		return true, nil
	})

	return result, nil
}
