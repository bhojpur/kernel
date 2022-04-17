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
	"os"
)

// Path first expands environment variables like $FOO or ${FOO}, and then
// reports if any of the sources have been modified more recently than the
// destination. Path does not descend into directories, it literally just checks
// the modtime of each thing you pass to it. If the destination file doesn't
// exist, it always returns true and nil. It's an error if any of the sources
// don't exist.
func Path(dst string, sources ...string) (bool, error) {
	stat, err := os.Stat(os.ExpandEnv(dst))
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return PathNewer(stat.ModTime(), sources...)
}

// Glob expands each of the globs (file patterns) into individual sources and
// then calls Path on the result, reporting if any of the resulting sources have
// been modified more recently than the destination. Syntax for Glob patterns is
// the same as stdlib's filepath.Glob. Note that Glob does not expand
// environment variables before globbing -- env var expansion happens during
// the call to Path. It is an error for any glob to return an empty result.
func Glob(dst string, globs ...string) (bool, error) {
	stat, err := os.Stat(os.ExpandEnv(dst))
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return GlobNewer(stat.ModTime(), globs...)
}

// Dir reports whether any of the sources have been modified more recently
// than the destination. If a source or destination is a directory, this
// function returns true if a source has any file that has been modified more
// recently than the most recently modified file in dst. If the destination
// file doesn't exist, it always returns true and nil.  It's an error if any
// of the sources don't exist.
func Dir(dst string, sources ...string) (bool, error) {
	dst = os.ExpandEnv(dst)
	stat, err := os.Stat(dst)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	destTime := stat.ModTime()
	if stat.IsDir() {
		destTime, err = NewestModTime(dst)
		if err != nil {
			return false, err
		}
	}
	return DirNewer(destTime, sources...)
}
