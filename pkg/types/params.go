package types

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

import "github.com/bhojpur/kernel/pkg/config"

type RunInstanceParams struct {
	Name                 string
	ImageId              string
	MntPointsToVolumeIds map[string]string
	Env                  map[string]string
	InstanceMemory       int
	NoCleanup            bool
	DebugMode            bool
}

type StageImageParams struct {
	Name      string
	RawImage  *RawImage
	Force     bool
	NoCleanup bool
}

type CreateVolumeParams struct {
	Name      string
	ImagePath string
	NoCleanup bool
}

type CompileImageParams struct {
	SourcesDir string
	Args       string
	MntPoints  []string
	NoCleanup  bool
	SizeMB     int
}

type PullImagePararms struct {
	Config    config.HubConfig
	ImageName string
	Force     bool
}

type PushImagePararms struct {
	Config    config.HubConfig
	ImageName string
}

type RemoteDeleteImagePararms struct {
	Config    config.HubConfig
	ImageName string
}
