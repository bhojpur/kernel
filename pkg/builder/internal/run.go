package internal

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
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var debug *log.Logger = log.New(ioutil.Discard, "", 0)

func SetDebug(l *log.Logger) {
	debug = l
}

func RunDebug(cmd string, args ...string) error {
	env, err := EnvWithCurrentGOOS()
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}
	debug.Println("running", cmd, strings.Join(args, " "))
	c := exec.Command(cmd, args...)
	c.Env = env
	c.Stderr = errbuf
	c.Stdout = buf
	if err := c.Run(); err != nil {
		debug.Print("error running '", cmd, strings.Join(args, " "), "': ", err, ": ", errbuf)
		return err
	}
	debug.Println(buf)
	return nil
}

func OutputDebug(cmd string, args ...string) (string, error) {
	env, err := EnvWithCurrentGOOS()
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}
	debug.Println("running", cmd, strings.Join(args, " "))
	c := exec.Command(cmd, args...)
	c.Env = env
	c.Stderr = errbuf
	c.Stdout = buf
	if err := c.Run(); err != nil {
		errMsg := strings.TrimSpace(errbuf.String())
		debug.Print("error running '", cmd, strings.Join(args, " "), "': ", err, ": ", errMsg)
		return "", fmt.Errorf("error running \"%s %s\": %s\n%s", cmd, strings.Join(args, " "), err, errMsg)
	}
	return strings.TrimSpace(buf.String()), nil
}

// splitEnv takes the results from os.Environ() (a []string of foo=bar values)
// and makes a map[string]string out of it.
func splitEnv(env []string) (map[string]string, error) {
	out := map[string]string{}

	for _, s := range env {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("badly formatted environment variable: %v", s)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}

// joinEnv converts the given map into a list of foo=bar environment variables,
// such as that outputted by os.Environ().
func joinEnv(env map[string]string) []string {
	vals := make([]string, 0, len(env))
	for k, v := range env {
		vals = append(vals, k+"="+v)
	}
	return vals
}

// EnvWithCurrentGOOS returns a copy of os.Environ with the GOOS and GOARCH set
// to runtime.GOOS and runtime.GOARCH.
func EnvWithCurrentGOOS() ([]string, error) {
	vals, err := splitEnv(os.Environ())
	if err != nil {
		return nil, err
	}
	vals["GOOS"] = runtime.GOOS
	vals["GOARCH"] = runtime.GOARCH
	return joinEnv(vals), nil
}

// EnvWithGOOS retuns the os.Environ() values with GOOS and/or GOARCH either set
// to their runtime value, or the given value if non-empty.
func EnvWithGOOS(goos, goarch string) ([]string, error) {
	env, err := splitEnv(os.Environ())
	if err != nil {
		return nil, err
	}
	if goos == "" {
		env["GOOS"] = runtime.GOOS
	} else {
		env["GOOS"] = goos
	}
	if goarch == "" {
		env["GOARCH"] = runtime.GOARCH
	} else {
		env["GOARCH"] = goarch
	}
	return joinEnv(env), nil
}
