package includeos

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
	goerrors "errors"
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

type IncludeosVirtualboxCompiler struct{}

func (i *IncludeosVirtualboxCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir
	env := make(map[string]string)
	if err := kutil.NewContainer("compilers-includeos-cpp-hw").WithVolume(sourcesDir, "/opt/code").WithEnvs(env).Run(); err != nil {
		return nil, err
	}

	res := &types.RawImage{}
	localImageFile, err := i.findFirstImageFile(sourcesDir)
	if err != nil {
		return nil, errors.New("error getting local image file name", err)
	}
	res.LocalImagePath = path.Join(sourcesDir, localImageFile)
	res.StageSpec.ImageFormat = types.ImageFormat_RAW
	res.RunSpec.StorageDriver = types.StorageDriver_IDE
	res.RunSpec.DefaultInstanceMemory = 256
	return res, nil
}

func (i *IncludeosVirtualboxCompiler) findFirstImageFile(directory string) (string, error) {
	dir, err := os.Open(directory)
	if err != nil {
		return "", errors.New("could not open dir", err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return "", errors.New("could not read dir", err)
	}
	for _, file := range files {
		logrus.Debugf("searching for .img file: %v", file.Name())
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".img" {
				return file.Name(), nil
			}
		}
	}
	return "", errors.New("no image file found", goerrors.New("end of dir"))
}

func (r *IncludeosVirtualboxCompiler) Usage() *compilers.CompilerUsage {
	return nil
}
