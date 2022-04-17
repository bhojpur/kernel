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
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	bhojpurModulePath = "github.com/bhojpur/kernel/pkg/base"
	bhojpurImportFile = "import_bhojpur.go"
	overlayFile       = "overlay.json"
)

var (
	bhojpurImportTpl = template.Must(template.New("bhojpur").Parse(`
	//+build bhojpur

	package {{.name}}
	import _ "github.com/bhojpur/kernel/pkg/base"
	`))
)

type gomodule struct {
	Module struct {
		Path string `json:"Path"`
	} `json:"Module"`
	Go      string `json:"Go"`
	Require []struct {
		Path    string `json:"Path"`
		Version string `json:"Version"`
	} `json:"Require"`
	Exclude interface{} `json:"Exclude"`
	Replace interface{} `json:"Replace"`
	Retract interface{} `json:"Retract"`
}

type buildOverlay struct {
	Replace map[string]string
}

func (b *Builder) bhojpurImportFile() string {
	return filepath.Join(b.basedir, bhojpurImportFile)
}

func (b *Builder) overlayFile() string {
	return filepath.Join(b.basedir, overlayFile)
}

func (b *Builder) buildPrepare() error {
	var err error

	if !b.modHasBhojpur() {
		log.Printf("bhojpur kernel not found in go.mod")
		err = b.editGoMod()
		if err != nil {
			return err
		}
	}

	err = b.writeImportFile(b.bhojpurImportFile())
	if err != nil {
		return err
	}

	err = writeOverlayFile(b.overlayFile(), bhojpurImportFile, b.bhojpurImportFile())
	if err != nil {
		return err
	}
	return nil
}

func writeOverlayFile(overlayFile, dest, source string) error {
	overlay := buildOverlay{
		Replace: map[string]string{
			dest: source,
		},
	}
	buf, _ := json.Marshal(overlay)
	return os.WriteFile(overlayFile, buf, 0644)
}

func (b *Builder) readGomodule() (*gomodule, error) {
	var buf bytes.Buffer
	cmd := exec.Command(b.gobin(), "mod", "edit", "-json")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var mod gomodule
	err = json.Unmarshal(buf.Bytes(), &mod)
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

func (b *Builder) modHasBhojpur() bool {
	if b.currentModulePath() == bhojpurModulePath {
		return true
	}

	mods, err := b.readGomodule()
	if err != nil {
		panic(err)
	}
	for _, mod := range mods.Require {
		if mod.Path == bhojpurModulePath {
			return true
		}
	}
	return false
}

func (b *Builder) editGoMod() error {
	getPath := bhojpurModulePath
	if b.cfg.KernelVersion != "" {
		getPath = getPath + "@" + b.cfg.KernelVersion
	}
	log.Printf("go get %s", getPath)
	env := []string{
		"GOOS=linux",
		"GOARCH=amd64",
	}
	env = append(env, os.Environ()...)
	cmd := exec.Command(b.gobin(), "get", getPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (b *Builder) currentPkgName() string {
	out, err := exec.Command(b.gobin(), "list", "-f", `{{.Name}}`).CombinedOutput()
	if err != nil {
		log.Panicf("get current package name:%s", out)
	}
	return strings.TrimSpace(string(out))
}

func (b *Builder) currentModulePath() string {
	out, err := exec.Command(b.gobin(), "list", "-f", `{{.Module.Path}}`).CombinedOutput()
	if err != nil {
		log.Panicf("get current module path:%s", out)
	}
	return strings.TrimSpace(string(out))

}

func (b *Builder) writeImportFile(fname string) error {
	pkgname := b.currentPkgName()
	var rawFile bytes.Buffer
	err := bhojpurImportTpl.Execute(&rawFile, map[string]interface{}{
		"name": pkgname,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(fname, rawFile.Bytes(), 0644)
	return err
}
