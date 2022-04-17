package log

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
	"os"

	"github.com/bhojpur/kernel/pkg/base/console"
	"github.com/bhojpur/kernel/pkg/base/drivers/uart"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
)

type LogLevel int8

const (
	LoglvlDebug LogLevel = iota
	LoglvlInfo
	LoglvlWarn
	LoglvlError
	LoglvlNone
)

const (
	loglvlEnv      = "BHOJPUR_KERNEL_LOG_LEVEL"
	loglvlEnvDebug = "debug"
	loglvlEnvInfo  = "info"
	loglvlEnvWarn  = "warn"
	loglvlEnvError = "error"
	loglvlEnvNone  = "none"

	defaultLoglvl = LoglvlError
)

var (
	Level LogLevel

	ErrInvalidLogLevel = fmt.Errorf("invalid log level")
)

func init() {
	lvl := os.Getenv("BHOJPUR_KERNEL_LOG_LEVEL")
	switch lvl {
	case loglvlEnvDebug:
		Level = LoglvlDebug
	case loglvlEnvInfo:
		Level = LoglvlInfo
	case loglvlEnvWarn:
		Level = LoglvlWarn
	case loglvlEnvError:
		Level = LoglvlError
	default:
		Level = defaultLoglvl
	}
}

func SetLevel(l LogLevel) error {
	if l < LoglvlDebug || l > LoglvlNone {
		return ErrInvalidLogLevel
	}

	Level = l

	return nil
}

func logf(lvl LogLevel, fmtstr string, args ...interface{}) {
	if lvl < Level {
		return
	}

	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, fmtstr, args...)
	buf.WriteByte('\n')
	buf.WriteTo(console.Console())
}

func Debugf(fmtstr string, args ...interface{}) {
	logf(LoglvlDebug, fmtstr, args...)
}

func Infof(fmtstr string, args ...interface{}) {
	logf(LoglvlInfo, fmtstr, args...)
}

func Warnf(fmtstr string, args ...interface{}) {
	logf(LoglvlWarn, fmtstr, args...)
}

func Errorf(fmtstr string, args ...interface{}) {
	logf(LoglvlError, fmtstr, args...)
}

//go:nosplit
func PrintStr(s string) {
	uart.WriteString(s)
}

const hextab = "0123456789abcdef"

//go:nosplit
func PrintHex(n uintptr) {
	shift := sys.PtrSize*8 - 4
	for ; shift > 0; shift = shift - 4 {
		v := (n >> shift) & 0x0F
		ch := hextab[v]
		uart.WriteByte(ch)
	}
	uart.WriteByte(hextab[n&0x0F])
}
