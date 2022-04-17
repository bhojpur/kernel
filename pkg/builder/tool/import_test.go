package tool

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
	"testing"
)

func TestBuilderImportsList(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		List:   true,
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := `
Targets:
  buildSubdir        Builds stuff.
  ns:deploy          deploys stuff.
  root               
  zz:buildSubdir2    Builds stuff.
  zz:ns:deploy2*     deploys stuff.

* default target
`[1:]

	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}

func TestBuilderImportsHelp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Help:   true,
		Args:   []string{"buildSubdir"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := `
BuildSubdir Builds stuff.

Usage:

	builder buildsubdir

`[1:]

	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}

func TestBuilderImportsHelpNamed(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Help:   true,
		Args:   []string{"zz:buildSubdir2"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := `
BuildSubdir2 Builds stuff.

Usage:

	builder zz:buildsubdir2

`[1:]

	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}

func TestBuilderImportsHelpNamedNS(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Help:   true,
		Args:   []string{"zz:ns:deploy2"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := `
Deploy2 deploys stuff.

Usage:

	builder zz:ns:deploy2

Aliases: nsd2

`[1:]

	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}

func TestBuilderImportsRoot(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"root"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "root\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsNamedNS(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"zz:nS:deploy2"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "deploy2\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsNamedRoot(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"zz:buildSubdir2"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "buildsubdir2\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
	if stderr := stderr.String(); stderr != "" {
		t.Fatal("unexpected output to stderr: ", stderr)
	}
}

func TestBuilderImportsRootImportNS(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"nS:deploy"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "deploy\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsRootImport(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"buildSubdir"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "buildsubdir\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsAliasToNS(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"nsd2"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "deploy2\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsOneLine(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport/oneline",
		Stdout: stdout,
		Stderr: stderr,
		Args:   []string{"build"},
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := "build\n"
	if actual != expected {
		t.Fatalf("expected: %q got: %q", expected, actual)
	}
}

func TestBuilderImportsTaggedPackage(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport/tagged",
		Stdout: stdout,
		Stderr: stderr,
		List:   true,
	}

	code := Invoke(inv)
	if code != 1 {
		t.Fatalf("expected to exit with code 1, but got %v, stdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}

	actual := stderr.String()
	// Match a shorter version of the error message, since the output from go list differs between versions
	expected := `
Error parsing builderfiles: error running "go list -f {{.Dir}}||{{.Name}} github.com/bhojpur/kernel/pkg/builder/tool/testdata/builderimport/tagged/pkg": exit status 1`[1:]
	actualShortened := actual[:len(expected)]
	if actualShortened != expected {
		t.Logf("expected: %q", expected)
		t.Logf("actual: %q", actualShortened)
		t.Fatalf("expected:\n%s\n\ngot:\n%s", expected, actualShortened)
	}
}

func TestBuilderImportsSameNamespaceUniqueTargets(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport/samenamespace/uniquetargets",
		Stdout: stdout,
		Stderr: stderr,
		List:   true,
	}

	code := Invoke(inv)
	if code != 0 {
		t.Fatalf("expected to exit with code 0, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stdout.String()
	expected := `
Targets:
  samenamespace:build1    
  samenamespace:build2    
`[1:]

	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}

func TestBuilderImportsSameNamespaceDupTargets(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/builderimport/samenamespace/duptargets",
		Stdout: stdout,
		Stderr: stderr,
		List:   true,
	}

	code := Invoke(inv)
	if code != 1 {
		t.Fatalf("expected to exit with code 1, but got %v, stderr:\n%s", code, stderr)
	}
	actual := stderr.String()
	expected := `
Error parsing builderfiles: "samenamespace:build" target has multiple definitions: github.com/bhojpur/kernel/pkg/builder/tool/testdata/builderimport/samenamespace/duptargets/package1.Build, github.com/bhojpur/kernel/pkg/builder/tool/testdata/builderimport/samenamespace/duptargets/package2.Build

`[1:]
	if actual != expected {
		t.Logf("expected: %q", expected)
		t.Logf("  actual: %q", actual)
		t.Fatalf("expected:\n%v\n\ngot:\n%v", expected, actual)
	}
}
