package os

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
	"io"
	"os"
	"os/exec"
	"path"

	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func RunLogCommand(name string, args ...string) error {
	log.WithFields(log.Fields{"cmd": name, "args": args}).Debug("running " + name)

	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"out": string(out)}).Error(name + " failed")

	}
	return err
}

func GetDirSize(dir string) (int64, error) {

	stat, err := os.Stat(dir)
	if err != nil {
		return 0, err
	}

	if !stat.IsDir() {
		return stat.Size(), nil
	} else {
		entries, err := listDir(dir)
		if err != nil {
			return 0, err
		}
		var sum int64 = 0
		for _, obj := range entries {
			curSize, err := GetDirSize(path.Join(dir, obj.Name()))
			if err != nil {
				return 0, err
			}
			sum += curSize
		}
		return sum, nil

	}

}

func listDir(path string) ([]os.FileInfo, error) {
	directory, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer directory.Close()
	return directory.Readdir(-1)
}

func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {

		sourcefilepointer := path.Join(source, obj.Name())

		destinationfilepointer := path.Join(dest, obj.Name())

		sfi, err := os.Stat(sourcefilepointer)
		if err != nil {
			return err
		}
		if sfi.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		}
	}

	return
}

/// http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang/21067803#21067803

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return nil
		}
	}
	if err = os.Link(src, dst); err == nil {
		return nil
	}
	err = copyFileContents(src, dst)
	return err
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		log.Debugf("total size %v after adding file %s", (int64(size)>>20)+10, info.Name())
		return err
	})
	return size, err
}
