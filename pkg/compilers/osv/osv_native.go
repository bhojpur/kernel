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
	"fmt"

	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
)

type OSvNativeCompiler struct {
	ImageFinisher ImageFinisher
}

func (r *OSvNativeCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {

	// Prepare meta/run.yaml for node runtime.
	if err := addRuntimeStanzaToMetaRun(params.SourcesDir, "native"); err != nil {
		return nil, err
	}

	// Create meta/package.yaml if not exist.
	if err := assureMetaPackage(params.SourcesDir); err != nil {
		return nil, err
	}

	// Parse image size from manifest.yaml.
	params.SizeMB = int(readImageSizeFromManifest(params.SourcesDir))

	// Compose image inside Docker container.
	imagePath, err := CreateImageDynamic(params, r.ImageFinisher.UseEc2())
	if err != nil {
		return nil, err
	}

	// And finalize it.
	convertParams := FinishParams{
		CompileParams:    params,
		CapstanImagePath: imagePath,
	}
	return r.ImageFinisher.FinishImage(convertParams)
}

func (r *OSvNativeCompiler) Usage() *compilers.CompilerUsage {
	return &compilers.CompilerUsage{
		PrepareApplication: `
(this is only needed if you want to run your own C/C++ application)
Compile your application into relocatable shared-object (a file normally
given a ".so" extension) that is PIC (position independent code).
`,
		ConfigurationFiles: map[string]string{
			"/meta/run.yaml": `
config_set:
   conf1:
      bootcmd: <boot-command-that-starts-application>    
config_set_default: conf1
`,
			"/meta/package.yaml": `
title: <your-unikernel-title>
name: <your-unikernel-name>
author: <your-name>
require:
  - <first-required-package-title>
  - <second-required-package-title>
  # ...
`,
			"/manifest.yaml": `
image_size: "10GB"  # logical image size
`,
		},
		Other: fmt.Sprintf(`
Below please find a list of packages in remote repository:
%s
`, listOfPackages()),
	}
}

func listOfPackages() string {
	return `
eu.mikelangelo-project.app.hadoop-hdfs
eu.mikelangelo-project.app.mysql-5.6.21
eu.mikelangelo-project.erlang
eu.mikelangelo-project.ompi
eu.mikelangelo-project.openfoam.core
eu.mikelangelo-project.openfoam.pimplefoam
eu.mikelangelo-project.openfoam.pisofoam
eu.mikelangelo-project.openfoam.poroussimplefoam
eu.mikelangelo-project.openfoam.potentialfoam
eu.mikelangelo-project.openfoam.rhoporoussimplefoam
eu.mikelangelo-project.openfoam.rhosimplefoam
eu.mikelangelo-project.openfoam.simplefoam
eu.mikelangelo-project.osv.cli
eu.mikelangelo-project.osv.cloud-init
eu.mikelangelo-project.osv.httpserver
eu.mikelangelo-project.osv.nfs
`
}
