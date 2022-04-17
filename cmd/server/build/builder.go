package build

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
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	WorkDir       string
	GoRoot        string
	Basedir       string
	BuildTest     bool
	KernelVersion string
	GoArgs        []string
}

type Builder struct {
	cfg     Config
	basedir string
}

func NewBuilder(cfg Config) *Builder {
	return &Builder{
		cfg: cfg,
	}
}

func (b *Builder) Build() error {
	if b.cfg.Basedir == "" {
		basedir, err := ioutil.TempDir("", "bhojpur-build")
		if err != nil {
			return err
		}
		b.basedir = basedir
		defer os.RemoveAll(basedir)
	} else {
		b.basedir = b.cfg.Basedir
	}

	err := b.buildPrepare()
	if err != nil {
		return err
	}

	return b.buildPkg()
}

func (b *Builder) gobin() string {
	if b.cfg.GoRoot == "" {
		return "go"
	}
	return filepath.Join(b.cfg.GoRoot, "bin", "go")
}

func (b *Builder) fixGoTags() bool {
	args := b.cfg.GoArgs
	for i, arg := range args {
		if arg == "-tags" {
			if i >= len(b.cfg.GoArgs)-1 {
				return false
			}
			idx := i + 1
			tags := args[idx]
			tags += " bhojpur"
			args[idx] = tags
			return true
		}
	}
	return false
}

func (b *Builder) buildPkg() error {
	var buildArgs []string
	ldflags := "-E github.com/bhojpur/kernel/kernel.rt0 -T 0x100000"
	if !b.cfg.BuildTest {
		buildArgs = append(buildArgs, "build")
	} else {
		buildArgs = append(buildArgs, "test", "-c")
	}
	hasGoTags := b.fixGoTags()
	if !hasGoTags {
		buildArgs = append(buildArgs, "-tags", "bhojpur")
	}
	buildArgs = append(buildArgs, "-ldflags", ldflags)
	buildArgs = append(buildArgs, "-overlay", b.overlayFile())
	buildArgs = append(buildArgs, b.cfg.GoArgs...)

	env := append([]string{}, os.Environ()...)
	env = append(env, []string{
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	}...)

	cmd := exec.Command(b.gobin(), buildArgs...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if b.cfg.WorkDir != "" {
		cmd.Dir = b.cfg.WorkDir
	}
	err := cmd.Run()
	return err
}
