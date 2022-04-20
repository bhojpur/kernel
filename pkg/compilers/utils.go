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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bhojpur/kernel/pkg/util/errors"

	kos "github.com/bhojpur/kernel/pkg/os"
	kutil "github.com/bhojpur/kernel/pkg/util"
)

func BuildBootableImage(kernel, cmdline string, usePartitionTables, noCleanup bool) (string, error) {
	directory, err := ioutil.TempDir("", "bootable-image-directory.")
	if err != nil {
		return "", errors.New("creating tmpdir", err)
	}
	if !noCleanup {
		defer os.RemoveAll(directory)
	}
	kernelBaseName := "program.bin"

	if err := kos.CopyDir(filepath.Dir(kernel), directory); err != nil {
		return "", errors.New("copying dir "+filepath.Dir(kernel)+" to "+directory, err)
	}

	if err := kos.CopyFile(kernel, path.Join(directory, kernelBaseName)); err != nil {
		return "", errors.New("copying kernel "+kernel+" to "+kernelBaseName, err)
	}

	tmpResultFile, err := ioutil.TempFile(directory, "boot-creator-result.img.")
	if err != nil {
		return "", err
	}
	tmpResultFile.Close()

	const contextDir = "/opt/vol/"
	cmds := []string{
		"-d", contextDir,
		"-p", kernelBaseName,
		"-a", cmdline,
		"-o", filepath.Base(tmpResultFile.Name()),
		fmt.Sprintf("-part=%v", usePartitionTables),
	}
	binds := map[string]string{directory: contextDir, "/dev/": "/dev/"}

	if err := kutil.NewContainer("boot-creator").Privileged(true).WithVolumes(binds).Run(cmds...); err != nil {
		return "", err
	}

	resultFile, err := ioutil.TempFile("", "boot-creator-result.img.")
	if err != nil {
		return "", err
	}
	resultFile.Close()

	if err := os.Rename(tmpResultFile.Name(), resultFile.Name()); err != nil {
		return "", errors.New("renaming "+tmpResultFile.Name()+" to "+resultFile.Name(), err)
	}
	return resultFile.Name(), nil
}
