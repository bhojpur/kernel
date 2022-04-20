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
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// uses rump docker conter container
// the container expectes code in /opt/code and will produce program.bin in the same folder.
// we need to take the program bin and combine with json config produce an image

const (
	BootstrapTypeEC2    = "ec2"
	BootstrapTypeUDP    = "udp"
	BootstrapTypeGCLOUD = "gcloud"
	BootstrapTypeNoStub = "nostub"
)

//compiler for building images from interpreted/scripting languages (python, javascript)
type RumpScriptCompiler struct {
	RumCompilerBase

	BootstrapType string //ec2 vs udp
	RunScriptArgs string
	ScriptEnv     []string
}

type scriptProjectConfig struct {
	MainFile    string `yaml:"main_file"`
	RuntimeArgs string `yaml:"runtime_args"`
}

func (r *RumpScriptCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {
	sourcesDir := params.SourcesDir
	var config scriptProjectConfig
	data, err := ioutil.ReadFile(filepath.Join(sourcesDir, "manifest.yaml"))
	if err != nil {
		return nil, errors.New("failed to read manifest.yaml file", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("failed to parse yaml manifest.yaml file", err)
	}

	if _, err := os.Stat(filepath.Join(sourcesDir, config.MainFile)); err != nil || config.MainFile == "" {
		return nil, errors.New("invalid main file specified", err)
	}

	logrus.Debugf("using main file %s", config.MainFile)

	containerEnv := []string{
		fmt.Sprintf("MAIN_FILE=%s", config.MainFile),
		fmt.Sprintf("BOOTSTRAP_TYPE=%s", r.BootstrapType),
	}

	if err := r.runContainer(sourcesDir, containerEnv); err != nil {
		return nil, err
	}

	resultFile := path.Join(sourcesDir, "program.bin")

	//build args string
	args := r.RunScriptArgs
	if config.RuntimeArgs != "" {
		args = config.RuntimeArgs + " " + args
	}
	if params.Args != "" {
		args = args + " " + params.Args
	}

	return r.CreateImage(resultFile, args, params.MntPoints, append(r.ScriptEnv, fmt.Sprintf("MAIN_FILE=%s", config.MainFile), fmt.Sprintf("BOOTSTRAP_TYPE=%s", r.BootstrapType)), params.NoCleanup)
}

func (r *RumpScriptCompiler) Usage() *compilers.CompilerUsage {
	return nil
}

func NewRumpPythonCompiler(dockerImage string, createImage func(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error), bootStrapType string) *RumpScriptCompiler {
	return &RumpScriptCompiler{
		RumCompilerBase: RumCompilerBase{
			DockerImage: dockerImage,
			CreateImage: createImage,
		},
		BootstrapType: bootStrapType,
		RunScriptArgs: "/bootpart/python-wrapper.py",
		ScriptEnv: []string{
			"PYTHONHOME=/bootpart/python",
			"PYTHONPATH=/bootpart/lib/python3.5/site-packages/:/bootpart/bin/",
		},
	}
}

func NewRumpJavaCompiler(dockerImage string, createImage func(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error), bootStrapType string) *RumpScriptCompiler {
	return &RumpScriptCompiler{
		RumCompilerBase: RumCompilerBase{
			DockerImage: dockerImage,
			CreateImage: createImage,
		},
		BootstrapType: bootStrapType,
		RunScriptArgs: "-jar /bootpart/program.jar",
		ScriptEnv: []string{
			"CLASSPATH=/bootpart/jetty:/bootpart/jdk/jre/lib",
			"JAVA_HOME=/bootpart/jdk/",
		},
	}
}
