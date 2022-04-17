package cga

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

	"github.com/bhojpur/kernel/pkg/base/drivers/cga/fbcga"
	"github.com/bhojpur/kernel/pkg/base/drivers/vbe"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
)

type Backend interface {
	GetPos() int
	SetPos(pos int)
	// WritePos write char at given pos but not update pos
	WritePos(pos int, char byte)
	// WriteByte write char and advance pos
	WriteByte(ch byte)
}

const (
	CRTPORT = 0x3d4
	bs      = '\b'
	del     = 0x7f
)

var (
	crt = (*[25 * 80]uint16)(unsafe.Pointer(uintptr(0xb8000)))
)

type cgabackend struct {
}

func (c *cgabackend) SetPos(pos int) {
	sys.Outb(CRTPORT, 14)
	sys.Outb(CRTPORT+1, byte(pos>>8))
	sys.Outb(CRTPORT, 15)
	sys.Outb(CRTPORT+1, byte(pos))
}

func (c *cgabackend) GetPos() int {
	var pos int

	// Cursor position: col + 80*row.
	sys.Outb(CRTPORT, 14)
	pos = int(sys.Inb(CRTPORT+1)) << 8
	sys.Outb(CRTPORT, 15)
	pos |= int(sys.Inb(CRTPORT + 1))
	return pos
}

func (c *cgabackend) WritePos(pos int, ch byte) {
	crt[pos] = uint16(ch) | 0x0700
}

func (c *cgabackend) WriteByte(ch byte) {
	var pos int

	// Cursor position: col + 80*row.
	sys.Outb(CRTPORT, 14)
	pos = int(sys.Inb(CRTPORT+1)) << 8
	sys.Outb(CRTPORT, 15)
	pos |= int(sys.Inb(CRTPORT + 1))

	switch ch {
	case '\n':
		pos += 80 - pos%80
	case bs, del:
		if pos > 0 {
			pos--
		}
	default:
		// black on white
		crt[pos] = uint16(ch&0xff) | 0x0700
		pos++
	}

	// Scroll up
	if pos/80 >= 25 {
		copy(crt[:], crt[80:25*80])
		pos -= 80
		s := crt[pos : 25*80]
		for i := range s {
			// mac qemu-M accel=hvf,`failed to decode instruction f 7f`
			// memclrNoHeapPointer
			if false {
			}
			s[i] = 0
		}
	}
	c.SetPos(pos)
	crt[pos] = ' ' | 0x0700
}

func getbackend() Backend {
	if vbe.IsEnable() {
		return &fbcga.Backend
	}
	return &backend
}
