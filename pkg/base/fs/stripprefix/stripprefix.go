package stripprefix

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

// The stripprefix strip the leading prefix of file name on access fs methods
// if file name is not an abs path, stripprefix do nothing

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// assert that stripprefix.fs implements afero.Fs.
var _ afero.Fs = (*fs)(nil)

type fs struct {
	prefix  string
	backend afero.Fs
}

func New(prefix string, backend afero.Fs) afero.Fs {
	return &fs{
		prefix:  prefix,
		backend: backend,
	}
}

func (f *fs) strip(name string) (string, error) {
	if !filepath.IsAbs(name) {
		return name, nil
	}
	p := strings.TrimPrefix(name, f.prefix)
	if len(p) == len(name) {
		return "", os.ErrNotExist
	}
	if p == "" {
		return ".", nil
	}
	return p, nil
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func (f *fs) Create(name string) (afero.File, error) {
	p, err := f.strip(name)
	if err != nil {
		return nil, err
	}
	return f.backend.Create(p)
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func (f *fs) Mkdir(name string, perm os.FileMode) error {
	p, err := f.strip(name)
	if err != nil {
		return err
	}
	return f.backend.Mkdir(p, perm)
}

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func (f *fs) MkdirAll(path string, perm os.FileMode) error {
	p, err := f.strip(path)
	if err != nil {
		return err
	}
	return f.backend.MkdirAll(p, perm)
}

// Open opens a file, returning it or an error, if any happens.
func (f *fs) Open(name string) (afero.File, error) {
	p, err := f.strip(name)
	if err != nil {
		return nil, err
	}
	return f.backend.Open(p)
}

// OpenFile opens a file using the given flags and the given mode.
func (f *fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	p, err := f.strip(name)
	if err != nil {
		return nil, err
	}
	return f.backend.OpenFile(p, flag, perm)
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func (f *fs) Remove(name string) error {
	p, err := f.strip(name)
	if err != nil {
		return err
	}
	return f.backend.Remove(p)
}

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func (f *fs) RemoveAll(path string) error {
	p, err := f.strip(path)
	if err != nil {
		return err
	}
	return f.backend.RemoveAll(p)
}

// Rename renames a file.
func (f *fs) Rename(oldname string, newname string) error {
	p1, err := f.strip(oldname)
	if err != nil {
		return err
	}
	p2, err := f.strip(newname)
	if err != nil {
		return err
	}
	return f.backend.Rename(p1, p2)
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func (f *fs) Stat(name string) (os.FileInfo, error) {
	p, err := f.strip(name)
	if err != nil {
		return nil, err
	}
	return f.backend.Stat(p)
}

// The name of this FileSystem
func (f *fs) Name() string {
	return "stripprefix"
}

//Chmod changes the mode of the named file to mode.
func (f *fs) Chmod(name string, mode os.FileMode) error {
	p, err := f.strip(name)
	if err != nil {
		return err
	}
	return f.backend.Chmod(p, mode)
}

// Chown changes the uid and gid of the named file.
func (f *fs) Chown(name string, uid, gid int) error {
	p, err := f.strip(name)
	if err != nil {
		return err
	}
	return f.backend.Chown(p, uid, gid)
}

//Chtimes changes the access and modification times of the named file
func (f *fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	p, err := f.strip(name)
	if err != nil {
		return err
	}
	return f.backend.Chtimes(p, atime, mtime)
}
