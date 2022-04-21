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
	"os"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type openstackHandle struct {
	AuthClient *gophercloud.ProviderClient
	Region     string
}

// MergeConfWithEnv overrides configuration with environment values (inplace).
func MergeConfWithEnv(conf *config.Openstack) {
	if v := os.Getenv("OS_AUTH_URL"); v != "" {
		conf.AuthUrl = v
	}
	if v := os.Getenv("OS_USER_ID"); v != "" {
		conf.UserId = v
	}
	if v := os.Getenv("OS_USERNAME"); v != "" {
		conf.UserName = v
	}
	if v := os.Getenv("OS_PASSWORD"); v != "" {
		conf.Password = v
	}
	if v := os.Getenv("OS_TENANT_ID"); v != "" {
		conf.TenantId = v
	}
	if v := os.Getenv("OS_TENANT_NAME"); v != "" {
		conf.TenantName = v
	}
	if v := os.Getenv("OS_DOMAIN_ID"); v != "" {
		conf.DomainId = v
	}
	if v := os.Getenv("OS_DOMAIN_NAME"); v != "" {
		conf.DomainName = v
	}
	if v := os.Getenv("OS_REGION_ID"); v != "" {
		conf.RegionId = v
	}
	if v := os.Getenv("OS_REGION_NAME"); v != "" {
		conf.RegionName = v
	}
}

// validateCredentials validates presence of required credentials.
func validateCredentials(conf *config.Openstack) error {
	// Validate
	if conf.AuthUrl == "" {
		return fmt.Errorf("Argument OS_AUTH_URL needs to be set.")
	}
	if conf.UserId == "" && conf.UserName == "" {
		return fmt.Errorf("Argument OS_USER_ID or OS_USERNAME needs to be set.")
	}
	if conf.Password == "" {
		return fmt.Errorf("Argument OS_PASSWORD needs to be set.")
	}
	if conf.TenantId == "" && conf.TenantName == "" {
		return fmt.Errorf("Argument OS_TENANT_ID or OS_TENANT_NAME needs to be set.")
	}
	if conf.RegionId == "" && conf.RegionName == "" {
		return fmt.Errorf("Argument OS_REGION_ID or OS_REGION_NAME needs to be set.")
	}

	return nil
}

// getHandle builds openstackHandle object that contains information needed
// to obtain any OpenStack API client (e.g. Nova client, Glance client).
// NOTE: this function performs a HTTP request to the OpenStack Keystone service
func getHandle(conf config.Openstack) (*openstackHandle, error) {
	if err := validateCredentials(&conf); err != nil {
		return nil, err
	}
	authClient, err := openstack.AuthenticatedClient(gophercloud.AuthOptions{
		IdentityEndpoint: conf.AuthUrl,
		UserID:           conf.UserId,
		Username:         conf.UserName,
		Password:         conf.Password,
		TenantID:         conf.TenantId,
		TenantName:       conf.TenantName,
		DomainID:         conf.DomainId,
		DomainName:       conf.DomainName,
	})
	if err != nil {
		return nil, err
	}

	region := conf.RegionId
	if region == "" {
		region = conf.RegionName
	}

	return &openstackHandle{
		AuthClient: authClient,
		Region:     region,
	}, nil
}

// getNovaClient returns ServiceClient for OpenStack Nova compute service API.
func getNovaClient(handle *openstackHandle) (*gophercloud.ServiceClient, error) {
	return openstack.NewComputeV2(handle.AuthClient, gophercloud.EndpointOpts{
		Region: handle.Region,
	})
}

// getGlanceClient returns ServiceClient for OpenStack Glance image service API.
func getGlanceClient(handle *openstackHandle) (*gophercloud.ServiceClient, error) {
	return openstack.NewImageServiceV2(handle.AuthClient, gophercloud.EndpointOpts{
		Region: handle.Region,
	})
}
