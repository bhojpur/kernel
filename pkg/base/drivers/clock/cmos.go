package clock

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
	"time"

	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
)

type CmosTime struct {
	Second int
	Minute int
	Hour   int
	Day    int
	Month  int
	Year   int
}

func ReadCmosTime() CmosTime {
	var t CmosTime
	for {
		readCmosTime(&t)
		if bcdDecode(readCmosSecond()) == t.Second {
			break
		}
	}
	return t
}

func (c *CmosTime) Time() time.Time {
	return time.Date(c.Year, time.Month(c.Month), c.Day, c.Hour, c.Minute, c.Second, 0, time.UTC)
}

// https://wiki.osdev.org/CMOS
func readCmosTime(t *CmosTime) {
	t.Year = bcdDecode(readCmosReg(0x09)) + bcdDecode(readCmosReg(0x32))*100
	t.Month = bcdDecode(readCmosReg(0x08))
	t.Day = bcdDecode(readCmosReg(0x07))
	t.Hour = bcdDecode(readCmosReg(0x04))
	t.Minute = bcdDecode(readCmosReg(0x02))
	t.Second = bcdDecode(readCmosReg(0x00))
}

func readCmosSecond() int {
	return readCmosReg(0x00)
}

// decode bcd format
func bcdDecode(v int) int {
	return v&0x0F + v/16*10
}

//go:nosplit
func readCmosReg(reg uint16) int {
	sys.Outb(0x70, 0x80|byte(reg))
	return int(sys.Inb(0x71))
}
