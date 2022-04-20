package config

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

type DaemonConfig struct {
	Providers Providers `yaml:"providers"`
	Version   string    `yaml:"version"`
}

type Providers struct {
	Aws         []Aws         `yaml:"aws"`
	Gcloud      []Gcloud      `yaml:"gcloud"`
	Vsphere     []Vsphere     `yaml:"vsphere"`
	Virtualbox  []Virtualbox  `yaml:"virtualbox"`
	Qemu        []Qemu        `yaml:"qemu"`
	Photon      []Photon      `yaml:"photon"`
	Xen         []Xen         `yaml:"xen"`
	Openstack   []Openstack   `yaml:"openstack"`
	Ukvm        []Ukvm        `yaml:"ukvm"`
	Firecracker []Firecracker `yaml:"firecracker"`
}

type Aws struct {
	Name   string `yaml:"name"`
	Region string `yaml:"region"`
	Zone   string `yaml:"zone"`
}

type Gcloud struct {
	Name      string `yaml:"name"`
	ProjectID string `yaml:"project_id"`
	Zone      string `yaml:"zone"`
}

type Vsphere struct {
	Name            string `yaml:"name"`
	VsphereUser     string `yaml:"vsphere_user"`
	VspherePassword string `yaml:"vsphere_password"`
	VsphereURL      string `yaml:"vsphere_url"`
	Datastore       string `yaml:"datastore"`
	Datacenter      string `yaml:"datacenter"`
	NetworkLabel    string `yaml:"network"`
}

type Photon struct {
	Name      string `yaml:"name"`
	PhotonURL string `yaml:"photon_url"`
	ProjectId string `yaml:"project_id"`
}

type Virtualbox struct {
	Name                  string                `yaml:"name"`
	AdapterName           string                `yaml:"adapter_name"`
	VirtualboxAdapterType VirtualboxAdapterType `yaml:"adapter_type"`
}

type Qemu struct {
	Name         string `yaml:"name"`
	NoGraphic    bool   `yaml:"no_graphic"`
	DebuggerPort int    `yaml:"debugger_port"`
}
type Firecracker struct {
	Name string `yaml:"name"`

	Binary string `yaml:"binary"`
	Kernel string `yaml:"kernel"`
	// either empty, stdio, or xterm
	Console string `yaml:"console"`
}

type Ukvm struct {
	Name string `yaml:"name"`
	Tap  string `yaml:"tap_device"`
}

type Xen struct {
	Name       string `yaml:"name"`
	KernelPath string `yaml:"pv_kernel"`
	XenBridge  string `yaml:"xen_bridge"`
}

type Openstack struct {
	Name string `yaml:"name"`

	UserName   string `yaml:"username"`
	UserId     string `yaml:"userid"`
	Password   string `yaml:"password"`
	AuthUrl    string `yaml:"auth_url"`
	TenantId   string `yaml:"tenant_id"`
	TenantName string `yaml:"tenant_name"`
	DomainId   string `yaml:"domain_id"`
	DomainName string `yaml:"domain_name"`

	ProjectName string `yaml:"project_name"`
	RegionId    string `yaml:"region_id"`
	RegionName  string `yaml:"region_name"`

	NetworkUUID string `yaml:"network_uuid"`
}

type VirtualboxAdapterType string

const (
	BridgedAdapter  = VirtualboxAdapterType("bridged")
	HostOnlyAdapter = VirtualboxAdapterType("host_only")
)

type ClientConfig struct {
	Host string `yaml:"host"`
}

type HubConfig struct {
	URL      string `yaml:"url",json:"url"`
	Username string `yaml:"user",json:"user"`
	Password string `yaml:"pass",json:"pass"`
}
