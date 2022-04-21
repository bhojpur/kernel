package kernel

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
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/kernel/isyscall"
)

// Timer depends on ePoll and pipe.
// When timer is frequently used, it will be more efficient
// to implement it in the kernel

const (
	pipeReadFd  = epollFd + 1
	pipeWriteFd = epollFd + 2
)

var (
	// FIXME:
	// avoid dup create pipe
	epollPipeCreated bool
	// the bytes number in pipe
	pipeBufferBytes int
)

//go:nosplit
func sysPipe2(req *isyscall.Request) {
	if epollPipeCreated {
		req.SetErrorNO(syscall.EINVAL)
		return
	}
	epollPipeCreated = true
	fds := (*[2]int32)(unsafe.Pointer(req.Arg(0)))
	fds[0] = pipeReadFd
	fds[1] = pipeWriteFd
	req.SetRet(0)
}

//go:nosplit
func sysPipeRead(req *isyscall.Request) {
	fd := req.Arg(0)
	buffer := req.Arg(1)
	len := req.Arg(2)
	_ = fd
	_ = buffer

	var n int
	if int(len) < pipeBufferBytes {
		n = int(len)
	} else {
		n = pipeBufferBytes
	}

	pipeBufferBytes -= n
	epollNotify(pipeWriteFd, syscall.EPOLLOUT)
	req.SetRet(uintptr(n))
}

//go:nosplit
func sysPipeWrite(req *isyscall.Request) {
	fd := req.Arg(0)
	buffer := req.Arg(1)
	len := req.Arg(2)
	_ = fd
	_ = buffer

	pipeBufferBytes += int(len)
	epollNotify(pipeReadFd, syscall.EPOLLIN)
	req.SetRet(len)
}
