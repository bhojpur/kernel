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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

// uses rump docker conter container
// the container expectes code in /opt/code and will produce program.bin in the same folder.
// we need to take the program bin and combine with json config produce an image

type RumpGoCompiler struct {
	RumCompilerBase
	BootstrapType string //ec2 | udp | nostub | gcloud
}

func (r *RumpGoCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir
	godepsFile := filepath.Join(sourcesDir, "Godeps", "Godeps.json")
	_, err := os.Stat(godepsFile)
	if err != nil {
		return nil, errors.New("the Go compiler requires Godeps file in the root of your project. see https://github.com/tools/godep", nil)
	}
	data, err := ioutil.ReadFile(godepsFile)
	if err != nil {
		return nil, errors.New("could not read godeps file", err)
	}
	var g godeps
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, errors.New("invalid json in godeps file", err)
	}
	containerEnv := []string{
		fmt.Sprintf("ROOT_PATH=%s", g.ImportPath),
		fmt.Sprintf("BOOTSTRAP_TYPE=%s", r.BootstrapType),
	}

	if err := r.runContainer(sourcesDir, containerEnv); err != nil {
		return nil, err
	}

	// now we should program.bin
	resultFile := path.Join(sourcesDir, "program.bin")
	logrus.Debugf("finished kernel binary at %s", resultFile)
	img, err := r.CreateImage(resultFile, params.Args, params.MntPoints, nil, params.NoCleanup)
	if err != nil {
		return nil, errors.New("creating boot volume from kernel binary", err)
	}
	return img, nil
}

func (r *RumpGoCompiler) Usage() *compilers.CompilerUsage {
	return nil
}

type godeps struct {
	ImportPath   string   `json:"ImportPath"`
	GoVersion    string   `json:"GoVersion"`
	GodepVersion string   `json:"GodepVersion"`
	Packages     []string `json:"Packages"`
	Deps         []struct {
		ImportPath string `json:"ImportPath"`
		Rev        string `json:"Rev"`
		Comment    string `json:"Comment,omitempty"`
	} `json:"Deps"`
}
