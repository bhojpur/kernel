package firecracker

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
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/sirupsen/logrus"
)

type FirecrackerCompiler struct{}

func (f *FirecrackerCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir

	// run dep ensure and go build
	if err := kutil.NewContainer("compilers-firecracker").Privileged(true).WithVolume(sourcesDir, "/opt/code").Run(); err != nil {
		return nil, err
	}
	res := &types.RawImage{}
	localImageFile, err := f.getImagefile(sourcesDir)
	if err != nil {
		logrus.Errorf("error getting local image file name")
	}
	res.LocalImagePath = localImageFile
	res.StageSpec.ImageFormat = types.ImageFormat_RAW
	res.RunSpec.DefaultInstanceMemory = 256
	return res, nil
}

func (f *FirecrackerCompiler) getImagefile(directory string) (string, error) {

	rootfs := filepath.Join(directory, "rootfs")

	_, err := os.Stat(rootfs)
	return rootfs, err
}

func (f *FirecrackerCompiler) Usage() *compilers.CompilerUsage {
	return nil
}
