package sh_test

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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bhojpur/kernel/pkg/builder/sh"
)

// compareFiles checks that two files are identical for testing purposes. That means they have the same length,
// the same contents, and the same permissions. It does NOT mean they have the same timestamp, as that is expected
// to change in normal Builder sh.Copy operation.
func compareFiles(file1 string, file2 string) error {
	s1, err := os.Stat(file1)
	if err != nil {
		return fmt.Errorf("can't stat %s: %v", file1, err)
	}
	s2, err := os.Stat(file2)
	if err != nil {
		return fmt.Errorf("can't stat %s: %v", file2, err)
	}
	if s1.Size() != s2.Size() {
		return fmt.Errorf("files %s and %s have different sizes: %d vs %d", file1, file2, s1.Size(), s2.Size())
	}
	if s1.Mode() != s2.Mode() {
		return fmt.Errorf("files %s and %s have different permissions: %#4o vs %#4o", file1, file2, s1.Mode(), s2.Mode())
	}
	f1bytes, err := ioutil.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", file1, err)
	}
	f2bytes, err := ioutil.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("can't read %s: %v", file2, err)
	}
	if !bytes.Equal(f1bytes, f2bytes) {
		return fmt.Errorf("files %s and %s have different contents", file1, file2)
	}
	return nil
}

func TestHelpers(t *testing.T) {

	mytmpdir, err := ioutil.TempDir("", "builder")
	if err != nil {
		t.Fatalf("can't create test directory: %v", err)
	}
	defer func() {
		derr := os.RemoveAll(mytmpdir)
		if derr != nil {
			fmt.Printf("error cleaning up after TestHelpers: %v", derr)
		}
	}()
	srcname := filepath.Join(mytmpdir, "test1.txt")
	err = ioutil.WriteFile(srcname, []byte("All work and no play makes Jack a dull boy."), 0644)
	if err != nil {
		t.Fatalf("can't create test file %s: %v", srcname, err)
	}
	destname := filepath.Join(mytmpdir, "test2.txt")

	t.Run("sh/copy", func(t *testing.T) {
		cerr := sh.Copy(destname, srcname)
		if cerr != nil {
			t.Errorf("test file copy from %s to %s failed: %v", srcname, destname, cerr)
		}
		cerr = compareFiles(srcname, destname)
		if cerr != nil {
			t.Errorf("test file copy verification failed: %v", cerr)
		}
	})

	// While we've got a temporary directory, test how forgiving sh.Rm is
	t.Run("sh/rm/ne", func(t *testing.T) {
		nef := filepath.Join(mytmpdir, "file_not_exist.txt")
		rerr := sh.Rm(nef)
		if rerr != nil {
			t.Errorf("sh.Rm complained when removing nonexistent file %s: %v", nef, rerr)
		}
	})

	t.Run("sh/copy/ne", func(t *testing.T) {
		nef := filepath.Join(mytmpdir, "file_not_exist.txt")
		nedf := filepath.Join(mytmpdir, "file_not_exist2.txt")
		cerr := sh.Copy(nedf, nef)
		if cerr == nil {
			t.Errorf("sh.Copy succeeded copying nonexistent file %s", nef)
		}
	})

	// We test sh.Rm by clearing up our own test files and directories
	t.Run("sh/rm", func(t *testing.T) {
		rerr := sh.Rm(destname)
		if rerr != nil {
			t.Errorf("failed to remove file %s: %v", destname, rerr)
		}
		rerr = sh.Rm(srcname)
		if rerr != nil {
			t.Errorf("failed to remove file %s: %v", srcname, rerr)
		}
		rerr = sh.Rm(mytmpdir)
		if rerr != nil {
			t.Errorf("failed to remove dir %s: %v", mytmpdir, rerr)
		}
		_, rerr = os.Stat(mytmpdir)
		if rerr == nil {
			t.Errorf("removed dir %s but it's still there?", mytmpdir)
		}
	})

	t.Run("sh/rm/nedir", func(t *testing.T) {
		rerr := sh.Rm(mytmpdir)
		if rerr != nil {
			t.Errorf("sh.Rm complained removing nonexistent dir %s", mytmpdir)
		}
	})

}
