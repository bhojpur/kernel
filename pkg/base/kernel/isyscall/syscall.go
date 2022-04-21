package isyscall

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
	"syscall"
	_ "unsafe"
)

const (
	EPANIC syscall.Errno = 0xfffff
)

var (
	handlers [512]Handler
)

//go:linkname wakeup github.com/bhojpur/kernel/pkg/base/kernel.wakeup
func wakeup(lock *uintptr, n int)

type Handler func(req *Request)

type Request struct {
	tf *trapFrame

	Lock uintptr
}

//go:nosplit
func (r *Request) NO() uintptr {
	return r.tf.NO()
}

//go:nosplit
func (r *Request) Arg(n int) uintptr {
	return r.tf.Arg(n)
}

//go:nosplit
func (r *Request) SetRet(v uintptr) {
	r.tf.SetRet(v)
}

//go:nosplit
func (r *Request) Ret() uintptr {
	return r.tf.Ret()
}

//go:nosplit
func (r *Request) SetErrorNO(errno syscall.Errno) {
	r.SetRet(Errno(errno))
}

//go:nosplit
func (r *Request) SetError(err error) {
	if err == nil {
		r.SetRet(0)
		return
	}
	r.SetRet(Error(err))
}

func (r *Request) Done() {
	wakeup(&r.Lock, 1)
}

func GetHandler(no uintptr) Handler {
	return handlers[no]
}

func Register(no uintptr, handler Handler) {
	handlers[no] = handler
}

func Errno(code syscall.Errno) uintptr {
	return uintptr(-code)
}

func Error(err error) uintptr {
	if err == nil {
		return 0
	}
	if code, ok := err.(syscall.Errno); ok {
		return Errno(code)
	}
	ret := uintptr(syscall.EINVAL)
	return -ret
}
