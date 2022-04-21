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
	"github.com/bhojpur/kernel/pkg/base/drivers/pic"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
	"github.com/bhojpur/kernel/pkg/base/kernel/trap"
	"gvisor.dev/gvisor/pkg/abi/linux"
)

const (
	_PIT_HZ = 1193180
	_HZ     = 100

	_IRQ_TIMER = pic.IRQ_BASE + pic.LINE_TIMER
)

const (
	ns     = 1
	ms     = 1000000 * ns
	second = 1000 * ms
)

var (
	// the counter of sched clock
	counter int64 = 1

	// the unix time of cmos read time
	baseUnixTime int64
	// the counter of cmos read time
	clockBaseCounter int64

	sleeplock uintptr
)

// pitCounter return the current counter of 8259a
//go:nosplit
func pitCounter() int32 {
	const div = (_PIT_HZ / _HZ)
	// Send the latch command to channel 0
	sys.Outb(0x43, 0)
	lo := sys.Inb(0x40)
	hi := sys.Inb(0x40)
	ax := (int32(hi)<<8 | int32(lo))
	return div - ax
}

//go:nosplit
func nanosecond() int64 {
	var t int64 = counter * (second / _HZ)
	elapse := int64(pitCounter()) * (second / _PIT_HZ)
	t += elapse
	return t
}

//go:nosplit
func clocktime() linux.Timespec {
	var ts linux.Timespec
	n := counter - clockBaseCounter
	ts.Sec = n/_HZ + baseUnixTime
	ts.Nsec = n % _HZ * (second / _HZ)
	ts.Nsec += int64(pitCounter()) * (second / _PIT_HZ)
	return ts
}

//go:nosplit
func nanosleep(tc *linux.Timespec) {
	deadline := nanosecond() + int64(tc.Nsec+tc.Sec*second)
	now := nanosecond()
	for now < deadline {
		sleepon(&sleeplock)
		now = nanosecond()
	}
}

//go:nosplit
func timerIntr() {
	counter++
	wakeup(&sleeplock, -1)
	pic.EOI(_IRQ_TIMER)
	Yield()
}

//go:nosplit
func timerInit() {
	div := int(_PIT_HZ / _HZ)
	sys.Outb(0x43, 0x36)
	sys.Outb(0x40, byte(div&0xff))
	sys.Outb(0x40, byte((div>>8)&0xff))
	trap.Register(_IRQ_TIMER, timerIntr)
	pic.EnableIRQ(pic.LINE_TIMER)
}
