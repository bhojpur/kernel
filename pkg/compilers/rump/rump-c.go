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

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"gopkg.in/yaml.v2"
)

//compiler for building images from interpreted/scripting languages (python, javascript)
type RumpCCompiler struct {
	RumCompilerBase
}

type cProjectConfig struct {
	BinaryName string `yaml:"binary_name"`
}

func (r *RumpCCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir
	var config cProjectConfig
	data, err := ioutil.ReadFile(filepath.Join(sourcesDir, "manifest.yaml"))
	if err != nil {
		return nil, errors.New("failed to read manifest.yaml file", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("failed to parse yaml manifest.yaml file", err)
	}

	containerEnv := []string{
		fmt.Sprintf("BINARY_NAME=%s", config.BinaryName),
	}

	if err := r.runContainer(sourcesDir, containerEnv); err != nil {
		return nil, err
	}

	resultFile := path.Join(sourcesDir, "program.bin")

	return r.CreateImage(resultFile, params.Args, params.MntPoints, nil, params.NoCleanup)
}

func (r *RumpCCompiler) Usage() *compilers.CompilerUsage {
	return nil
}

func NewRumpCCompiler(dockerImage string, createImage func(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error)) *RumpCCompiler {
	return &RumpCCompiler{
		RumCompilerBase: RumCompilerBase{
			DockerImage: dockerImage,
			CreateImage: createImage,
		},
	}
}
