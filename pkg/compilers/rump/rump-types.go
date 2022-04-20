package rump

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

type blk struct {
	Source     string `json:"source"`
	Path       string `json:"path"`
	FSType     string `json:"fstype"`
	MountPoint string `json:"mountpoint,omitempty"`
	DiskFile   string `json:"diskfile,omitempty"`
}

type Method string

const (
	Static Method = "static"
	DHCP   Method = "dhcp"
)

type net struct {
	If     string `json:"if,omitempty"`
	Type   string `json:"type,omitempty"`
	Method Method `json:"method,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Mask   string `json:"mask,omitempty"`
	Cloner string `json:"cloner,omitempty"`
}

type commandLine struct {
	Bin     string   `json:"bin"`
	Argv    []string `json:"argv"`
	Runmode *string  `json:"runmode,omitempty"`
}

type rumpConfig struct {
	Rc   []commandLine     `json:"rc"`
	Net  *net              `json:"net,omitempty"`
	Net1 *net              `json:"net1,omitempty"`
	Blk  []blk             `json:"blk,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
}
