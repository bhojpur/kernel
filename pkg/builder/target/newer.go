package target

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
	"os"
	"path/filepath"
	"time"
)

var (
	// errNewer is an ugly sentinel error to cause filepath.Walk to abort
	// as soon as a newer file is encountered
	errNewer = fmt.Errorf("newer item encountered")
)

// DirNewer reports whether any item in sources is newer than the target time.
// Sources are searched recursively and searching stops as soon as any entry
// is newer than the target.
func DirNewer(target time.Time, sources ...string) (bool, error) {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.ModTime().After(target) {
			return errNewer
		}
		return nil
	}
	for _, source := range sources {
		source = os.ExpandEnv(source)
		err := filepath.Walk(source, walkFn)
		if err == nil {
			continue
		}
		if err == errNewer {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// GlobNewer performs glob expansion on each source and passes the results to
// PathNewer for inspection. It returns the first time PathNewer encounters a
// newer file
func GlobNewer(target time.Time, sources ...string) (bool, error) {
	for _, g := range sources {
		files, err := filepath.Glob(g)
		if err != nil {
			return false, err
		}
		if len(files) == 0 {
			return false, fmt.Errorf("glob didn't match any files: %s", g)
		}
		newer, err := PathNewer(target, files...)
		if err != nil {
			return false, err
		}
		if newer {
			return true, nil
		}
	}
	return false, nil
}

// PathNewer checks whether any of the sources are newer than the target time.
// It stops at the first newer file it encounters. Each source path is passed
// through os.ExpandEnv.
func PathNewer(target time.Time, sources ...string) (bool, error) {
	for _, source := range sources {
		source = os.ExpandEnv(source)
		stat, err := os.Stat(source)
		if err != nil {
			return false, err
		}
		if stat.ModTime().After(target) {
			return true, nil
		}
	}
	return false, nil
}

// OldestModTime recurses a list of target filesystem objects and finds the
// the oldest ModTime among them.
func OldestModTime(targets ...string) (time.Time, error) {
	t := time.Now().Add(time.Hour * 100000)
	for _, target := range targets {
		walkFn := func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			mTime := info.ModTime()
			if mTime.Before(t) {
				t = mTime
			}
			return nil
		}
		if err := filepath.Walk(target, walkFn); err != nil {
			return t, err
		}
	}
	return t, nil
}

// NewestModTime recurses a list of target filesystem objects and finds the
// the newest ModTime among them.
func NewestModTime(targets ...string) (time.Time, error) {
	t := time.Time{}
	for _, target := range targets {
		walkFn := func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			mTime := info.ModTime()
			if mTime.After(t) {
				t = mTime
			}
			return nil
		}
		if err := filepath.Walk(target, walkFn); err != nil {
			return t, err
		}
	}
	return t, nil
}
