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
	"github.com/bhojpur/kernel/pkg/compilers"
	"github.com/bhojpur/kernel/pkg/types"
)

type OSvNodeCompiler struct {
	ImageFinisher ImageFinisher
}

func (r *OSvNodeCompiler) CompileRawImage(params types.CompileImageParams) (*types.RawImage, error) {

	// Prepare meta/run.yaml for node runtime.
	if err := addRuntimeStanzaToMetaRun(params.SourcesDir, "node"); err != nil {
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

func (r *OSvNodeCompiler) Usage() *compilers.CompilerUsage {
	return &compilers.CompilerUsage{
		PrepareApplication: "Install all libraries using `npm install`.",
		ConfigurationFiles: map[string]string{
			"/meta/run.yaml": `
config_set:
   conf1:
      main: <relative-path-to-your-entrypoint>   
config_set_default: conf1
`,
			"/manifest.yaml": `
image_size: "10GB"  # logical image size
`,
		},
	}
}
