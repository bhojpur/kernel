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
	"github.com/bhojpur/kernel/pkg/base/kernel/mm"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/abi/linux/errno"
)

const (
	epollFd = 3

	maxFds = 1024
)

var (
	// source of fd events, set by netstack
	// cleared by epoll_wait
	fdevents [maxFds]uint32

	// to manage epoll event
	eventpool mm.Pool

	// header of registered epoll events
	epollEvents epollEvent

	// notify of epoll events
	epollNote note
)

//go:notinheap
type epollEvent struct {
	fd  uintptr
	sub linux.EpollEvent

	pre, next *epollEvent
}

//go:nosplit
func newEpollEvent() *epollEvent {
	ptr := eventpool.Alloc()
	e := (*epollEvent)(unsafe.Pointer(ptr))
	e.pre = &epollEvents
	e.next = epollEvents.next
	if epollEvents.next != nil {
		epollEvents.next.pre = e
	}
	epollEvents.next = e
	return e
}

//go:nosplit
func freeEpollEvent(e *epollEvent) {
	e.pre.next = e.next
	if e.next != nil {
		e.next.pre = e.pre
	}
	eventpool.Free(uintptr(unsafe.Pointer(e)))
}

//go:nosplit
func findEpollEvent(fd uintptr) *epollEvent {
	for e := epollEvents.next; e != nil; e = e.next {
		if e.fd == fd {
			return e
		}
	}
	return nil
}

//go:nosplit
func epollCtl(epfd, op, fd, desc uintptr) uintptr {
	euser := (*linux.EpollEvent)(unsafe.Pointer(desc))
	var e *epollEvent
	switch op {
	case syscall.EPOLL_CTL_ADD:
		e = findEpollEvent(fd)
		if e == nil {
			e = newEpollEvent()
		}
		e.fd = fd
		e.sub = *euser
		e.sub.Events |= syscall.EPOLLHUP
		return 0
	case syscall.EPOLL_CTL_MOD:
		e = findEpollEvent(fd)
		if e == nil {
			return isyscall.Errno(errno.EINVAL)
		}
		e.sub = *euser
		e.sub.Events |= syscall.EPOLLHUP
		return 0
	case syscall.EPOLL_CTL_DEL:
		e = findEpollEvent(fd)
		if e == nil {
			return isyscall.Errno(errno.EINVAL)
		}
		freeEpollEvent(e)
		return 0
	default:
		return isyscall.Errno(errno.EINVAL)
	}
}

//go:nosplit
func epollWait(epfd, eventptr, len, _ms uintptr) uintptr {
	if _ms != 0 {
		ts := linux.Timespec{
			Sec:  int64(_ms / 1000),
			Nsec: int64(_ms%1000) * ms,
		}
		// wait fd event
		epollNote.sleep(&ts)
		epollNote.clear()
	}

	events := (*[256]linux.EpollEvent)(unsafe.Pointer(eventptr))[:len]
	var cnt uintptr = 0
	for e := epollEvents.next; e != nil && cnt < len; e = e.next {
		event := fdevents[e.fd]
		if event == 0 {
			continue
		}
		ue := &events[cnt]
		ue.Data = e.sub.Data
		ue.Events = event & e.sub.Events
		// clear events
		// FIXME: only clear masked events?
		fdevents[e.fd] = 0
		cnt++
	}
	return cnt
}

//go:nosplit
func epollNotify(fd, events uintptr) {
	fdevents[fd] |= uint32(events)
	epollNote.wakeup()
}

//go:nosplit
func epollInit() {
	mm.PoolInit(&eventpool, unsafe.Sizeof(epollEvent{}))
}
