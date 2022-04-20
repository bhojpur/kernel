package compilers

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
	"strings"

	"github.com/bhojpur/kernel/pkg/types"
)

type Compiler interface {
	CompileRawImage(params types.CompileImageParams) (*types.RawImage, error)

	// Usage describes how to prepare project to run it with Bhojpur Kernel
	// The returned text should describe what configuration files to
	// prepare and how.
	Usage() *CompilerUsage
}

type CompilerUsage struct {
	// PrepareApplication section briefly describes how user should
	// prepare her application PRIOR composing unikernel with Bhojpur Kernel
	PrepareApplication string

	// ConfigurationFiles lists configuration files needed by Bhojpur Kernel.
	// It is a map filename:content_description.
	ConfigurationFiles map[string]string

	// Other is arbitrary content that will be printed at the end.
	Other string
}

func (c *CompilerUsage) ToString() string {
	prepApp := strings.TrimSpace(c.PrepareApplication)
	other := strings.TrimSpace(c.Other)

	configFiles := ""
	for k := range c.ConfigurationFiles {
		configFiles += fmt.Sprintf("------ %s ------\n", k)
		configFiles += strings.TrimSpace(c.ConfigurationFiles[k])
		configFiles += "\n\n"
	}
	configFiles = strings.TrimSpace(configFiles)

	description := fmt.Sprintf(`
HOW TO PREPARE APPLICATION	
%s

CONFIGURATION FILES
%s
`, prepApp, configFiles)

	if other != "" {
		description += fmt.Sprintf("\nOTHER\n%s", other)
	}

	return description
}
