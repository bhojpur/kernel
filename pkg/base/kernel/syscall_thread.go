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
	"runtime"
	"syscall"
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/kernel/isyscall"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
	"github.com/bhojpur/kernel/pkg/base/log"
)

var (
	syscalltask threadptr

	// pendingCall is the address of pending forward syscall
	pendingCall uintptr
)

//go:nosplit
func forwardCall(call *isyscall.Request) {
	// wait syscall task fetch pendingCall
	for pendingCall != 0 {
		sleepon(&pendingCall)
	}
	pendingCall = uintptr(unsafe.Pointer(call))
	// tell syscall task pendingCall is avaiable
	// we can't only wakeup only one thread here
	wakeup(&pendingCall, -1)

	// wait on syscall task handle request
	sleepon(&call.Lock)

	// for debug purpose
	if -call.Ret() == uintptr(isyscall.EPANIC) {
		preparePanic(Mythread().tf)
	}
}

//go:nosplit
func fetchPendingCall() uintptr {
	// waiting someone call forward syscall
	for pendingCall == 0 {
		sleepon(&pendingCall)
	}
	ret := pendingCall
	pendingCall = 0
	// wakeup one thread, pendingCall is avaiable
	wakeup(&pendingCall, 1)
	return ret
}

// runSyscallThread run in normal go code space
func runSyscallThread() {
	runtime.LockOSThread()
	my := Mythread()
	syscalltask = (threadptr)(unsafe.Pointer(my))
	log.Infof("[syscall] tid:%d", my.id)
	for {
		callptr, _, err := syscall.Syscall(SYS_WAIT_SYSCALL, 0, 0, 0)
		if err != 0 {
			throw("bad SYS_WAIT_SYSCALL return")
		}
		call := (*isyscall.Request)(unsafe.Pointer(callptr))

		no := call.NO()
		handler := isyscall.GetHandler(no)
		if handler == nil {
			log.Errorf("[syscall] unhandled %s(%d)(0x%x, 0x%x, 0x%x, 0x%x, 0x%x, 0x%x)",
				syscallName(int(no)), no,
				call.Arg(0), call.Arg(1), call.Arg(2), call.Arg(3),
				call.Arg(4), call.Arg(5))
			call.SetErrorNO(syscall.ENOSYS)
			call.Done()
			continue
		}
		go func() {
			handler(call)
			var iret interface{}
			ret := call.Ret()
			if hasErrno(ret) {
				iret = syscall.Errno(-ret)
			} else {
				iret = ret
			}
			log.Debugf("[syscall] %s(%d)(0x%x, 0x%x, 0x%x, 0x%x, 0x%x, 0x%x) = %v",
				syscallName(int(no)), no,
				call.Arg(0), call.Arg(1), call.Arg(2), call.Arg(3),
				call.Arg(4), call.Arg(5), iret,
			)
			call.Done()
		}()
	}
}

func hasErrno(n uintptr) bool {
	return 1<<(sys.PtrSize*8-1)&n != 0
}
