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
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/drivers/multiboot"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
	"gvisor.dev/gvisor/pkg/abi/linux"
)

//go:nosplit
func envput(pbuf *[]byte, v uintptr) uintptr {
	buf := *pbuf
	p := unsafe.Pointer(&buf[0])
	// *p = v
	*(*uintptr)(p) = v
	// advance buffer
	*pbuf = buf[unsafe.Sizeof(v):]
	// return p
	return uintptr(unsafe.Pointer(&buf[0]))
}

// envptr used to alloc an *uintptr
//go:nosplit
func envptr(pbuf *[]byte) *uintptr {
	return (*uintptr)(unsafe.Pointer(envput(pbuf, 0)))
}

//go:nosplit
func envdup(pbuf *[]byte, s string) uintptr {
	buf := *pbuf
	copy(buf, s)
	*pbuf = buf[len(s):]
	return uintptr(unsafe.Pointer(&buf[0]))
}

//go:nosplit
func prepareArgs(sp uintptr) {
	buf := sys.UnsafeBuffer(sp, 256)

	var argc uintptr = 1
	// put argc slot
	envput(&buf, argc)
	arg0 := envptr(&buf)
	// end of args
	envput(&buf, 0)

	envTerm := envptr(&buf)
	envGoDebug := envptr(&buf)
	putKernelArgs(&buf)
	// end of env
	envput(&buf, 0)

	// put auxillary vector
	envput(&buf, linux.AT_PAGESZ)
	envput(&buf, sys.PageSize)
	envput(&buf, linux.AT_NULL)
	envput(&buf, 0)

	*arg0 = envdup(&buf, "eggos\x00")
	*envTerm = envdup(&buf, "TERM=xterm\x00")
	*envGoDebug = envdup(&buf, "GODEBUG=asyncpreemptoff=1\x00")
}

//go:nosplit
func putKernelArgs(pbuf *[]byte) uintptr {
	var cnt uintptr
	info := multiboot.BootInfo
	var flag = info.Flags
	if flag&multiboot.FlagInfoCmdline == 0 {
		return 0
	}
	var pos uintptr = uintptr(info.Cmdline)
	if pos == 0 {
		return cnt
	}

	var arg uintptr
	for {
		arg = strtok(&pos)
		if arg == 0 {
			break
		}
		envput(pbuf, arg)
		cnt++
	}
	return cnt
}

//go:nosplit
func strtok(pos *uintptr) uintptr {
	addr := *pos

	// skip spaces
	for {
		ch := bytedef(addr)
		if ch == 0 {
			return 0
		}
		if ch != ' ' {
			break
		}
		addr++
	}
	ret := addr
	// scan util read space and \0
	for {
		ch := bytedef(addr)
		if ch == ' ' {
			*(*byte)(unsafe.Pointer(addr)) = 0
			addr++
			break
		}
		if ch == 0 {
			break
		}
		addr++
	}
	*pos = addr
	return ret
}

//go:nosplit
func bytedef(addr uintptr) byte {
	return *(*byte)(unsafe.Pointer(addr))
}
