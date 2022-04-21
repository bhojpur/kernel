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

	"gvisor.dev/gvisor/pkg/abi/linux"
)

const (
	_FUTEX_WAIT         = 0
	_FUTEX_WAKE         = 1
	_FUTEX_PRIVATE_FLAG = 128
	_FUTEX_WAIT_PRIVATE = _FUTEX_WAIT | _FUTEX_PRIVATE_FLAG
	_FUTEX_WAKE_PRIVATE = _FUTEX_WAKE | _FUTEX_PRIVATE_FLAG
)

//go:nosplit
func futex(addr *uintptr, op, val uintptr, ts *linux.Timespec) {
	switch op {
	case _FUTEX_WAIT, _FUTEX_WAIT_PRIVATE:
		if ts != nil {
			sleeptimeout(addr, val, ts)
			return
		}
		for *addr == val {
			sleepon(addr)
		}
		return
	case _FUTEX_WAKE, _FUTEX_WAKE_PRIVATE:
		wakeup(addr, int(val))
	default:
		panic("futex: invalid op")
	}
}

//go:nosplit
func sleeptimeout(addr *uintptr, val uintptr, ts *linux.Timespec) {
	if ts == nil {
		panic("sleeptimeout: nil ts")
	}
	deadline := nanosecond() + int64(ts.Nsec) + int64(ts.Sec)*second
	// check on every timer intr
	now := nanosecond()
	t := Mythread()
	for now < deadline && *addr == val {
		t.timerKey = uintptr(unsafe.Pointer(&sleeplock))
		t.sleepKey = uintptr(unsafe.Pointer(addr))
		t.state = SLEEPING
		Sched()
		t.timerKey = 0
		t.sleepKey = 0
		now = nanosecond()
	}
}

//go:nosplit
func sleepon(lock *uintptr) {
	t := Mythread()
	t.sleepKey = uintptr(unsafe.Pointer(lock))
	t.state = SLEEPING
	Sched()
	t.sleepKey = 0
}

// wakeup thread sleep on lock, n == -1 means all threads
//go:nosplit
func wakeup(lock *uintptr, n int) {
	limit := uint(n)
	cnt := uint(0)
	lockKey := uintptr(unsafe.Pointer(lock))
	for i := 0; i < _NTHREDS; i++ {
		t := &threads[i]
		if (t.sleepKey == lockKey || t.timerKey == lockKey) && cnt < limit {
			cnt++
			t.state = RUNNABLE
		}
	}
}

type note uintptr

//go:nosplit
func (n *note) sleep(ts *linux.Timespec) {
	futex((*uintptr)(unsafe.Pointer(n)), _FUTEX_WAIT, 0, ts)
}

//go:nosplit
func (n *note) wakeup() {
	*n = 1
	futex((*uintptr)(unsafe.Pointer(n)), _FUTEX_WAKE, 1, nil)
}

//go:nosplit
func (n *note) clear() {
	*n = 0
}
