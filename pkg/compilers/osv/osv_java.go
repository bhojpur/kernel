package osv

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
	"io/ioutil"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type OSvJavaCompiler struct {
	ImageFinisher ImageFinisher
}

// javaProjectConfig defines available inputs
type javaProjectConfig struct {
	MainFile    string `yaml:"main_file"`
	RuntimeArgs string `yaml:"runtime_args"`
	BuildCmd    string `yaml:"build_command"`
}

func (r *OSvJavaCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir

	var config javaProjectConfig
	data, err := ioutil.ReadFile(filepath.Join(sourcesDir, "manifest.yaml"))
	if err != nil {
		return nil, errors.New("failed to read manifest.yaml file", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("failed to parse yaml manifest.yaml file", err)
	}

	container := kutil.NewContainer("compilers-osv-java").WithVolume("/dev", "/dev").WithVolume(sourcesDir+"/", "/project_directory")
	var args []string
	if r.ImageFinisher.UseEc2() {
		args = append(args, "-ec2")
	}

	args = append(args, "-main_file", config.MainFile)
	args = append(args, "-args", params.Args)
	if config.BuildCmd != "" {
		args = append(args, "-buildCmd", config.BuildCmd)
	}
	if len(config.RuntimeArgs) > 0 {
		args = append(args, "-runtime", config.RuntimeArgs)
	}

	logrus.WithFields(logrus.Fields{
		"args": args,
	}).Debugf("running compilers-osv-java container")

	if err := container.Run(args...); err != nil {
		return nil, errors.New("failed running compilers-osv-java on "+sourcesDir, err)
	}

	// And finally bootstrap.
	convertParams := FinishParams{
		CompileParams:    params,
		CapstanImagePath: filepath.Join(sourcesDir, "boot.qcow2"),
	}
	return r.ImageFinisher.FinishImage(convertParams)
}

func (r *OSvJavaCompiler) Usage() *compilers.CompilerUsage {
	return &compilers.CompilerUsage{
		PrepareApplication: "Compile your Java application into a fat jar.",
		ConfigurationFiles: map[string]string{
			"/manifest.yaml": "main_file: <relative-path-to-your-fat-jar>",
		},
	}
}
