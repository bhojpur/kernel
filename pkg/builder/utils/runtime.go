package utils

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
	"path/filepath"
	"runtime"
	"strconv"
)

// CacheEnv is the environment variable that users may set to change the
// location where builder stores its compiled binaries.
const CacheEnv = "BUILDERFILE_CACHE"

// VerboseEnv is the environment variable that indicates the user requested
// verbose mode when running a builderfile.
const VerboseEnv = "BUILDERFILE_VERBOSE"

// DebugEnv is the environment variable that indicates the user requested
// debug mode when running builder.
const DebugEnv = "BUILDERFILE_DEBUG"

// GoCmdEnv is the environment variable that indicates the go binary the user
// desires to utilize for builderfile compilation.
const GoCmdEnv = "BUILDERFILE_GOCMD"

// IgnoreDefaultEnv is the environment variable that indicates the user requested
// to ignore the default target specified in the builderfile.
const IgnoreDefaultEnv = "BUILDERFILE_IGNOREDEFAULT"

// HashFastEnv is the environment variable that indicates the user requested to
// use a quick hash of builderfiles to determine whether or not the builderfile binary
// needs to be rebuilt. This results in faster runtimes, but means that builder
// will fail to rebuild if a dependency has changed. To force a rebuild, run
// builder with the -f flag.
const HashFastEnv = "BUILDERFILE_HASHFAST"

// EnableColorEnv is the environment variable that indicates the user is using
// a terminal which supports a color output. The default is false for backwards
// compatibility. When the value is true and the detected terminal does support colors
// then the list of builder targets will be displayed in ANSI color. When the value
// is true but the detected terminal does not support colors, then the list of
// builder targets will be displayed in the default colors (e.g. black and white).
const EnableColorEnv = "BUILDERFILE_ENABLE_COLOR"

// TargetColorEnv is the environment variable that indicates which ANSI color
// should be used to colorize builder targets. This is only applicable when
// the BUILDERFILE_ENABLE_COLOR environment variable is true.
// The supported ANSI color names are any of these:
// - Black
// - Red
// - Green
// - Yellow
// - Blue
// - Magenta
// - Cyan
// - White
// - BrightBlack
// - BrightRed
// - BrightGreen
// - BrightYellow
// - BrightBlue
// - BrightMagenta
// - BrightCyan
// - BrightWhite
const TargetColorEnv = "BUILDERFILE_TARGET_COLOR"

// Verbose reports whether a builderfile was run with the verbose flag.
func Verbose() bool {
	b, _ := strconv.ParseBool(os.Getenv(VerboseEnv))
	return b
}

// Debug reports whether a builderfile was run with the debug flag.
func Debug() bool {
	b, _ := strconv.ParseBool(os.Getenv(DebugEnv))
	return b
}

// GoCmd reports the command that builder will use to build go code.  By default builder runs
// the "go" binary in the PATH.
func GoCmd() string {
	if cmd := os.Getenv(GoCmdEnv); cmd != "" {
		return cmd
	}
	return "go"
}

// HashFast reports whether the user has requested to use the fast hashing
// mechanism rather than rely on go's rebuilding mechanism.
func HashFast() bool {
	b, _ := strconv.ParseBool(os.Getenv(HashFastEnv))
	return b
}

// IgnoreDefault reports whether the user has requested to ignore the default target
// in the builderfile.
func IgnoreDefault() bool {
	b, _ := strconv.ParseBool(os.Getenv(IgnoreDefaultEnv))
	return b
}

// CacheDir returns the directory where builder caches compiled binaries.  It
// defaults to $HOME/.builderfile, but may be overridden by the BUILDERFILE_CACHE
// environment variable.
func CacheDir() string {
	d := os.Getenv(CacheEnv)
	if d != "" {
		return d
	}
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), "builderfile")
	default:
		return filepath.Join(os.Getenv("HOME"), ".builderfile")
	}
}

// EnableColor reports whether the user has requested to enable a color output.
func EnableColor() bool {
	b, _ := strconv.ParseBool(os.Getenv(EnableColorEnv))
	return b
}

// TargetColor returns the configured ANSI color name a color output.
func TargetColor() string {
	s, exists := os.LookupEnv(TargetColorEnv)
	if exists {
		if c, ok := getAnsiColor(s); ok {
			return c
		}
	}
	return DefaultTargetAnsiColor
}

// Namespace allows for the grouping of similar commands
type Namespace struct{}
