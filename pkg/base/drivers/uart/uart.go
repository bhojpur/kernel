package uart

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
)

const (
	com1      = uint16(0x3f8)
	_IRQ_COM1 = pic.IRQ_BASE + pic.LINE_COM1
)

var (
	inputCallback func(byte)
)

//go:nosplit
func ReadByte() int {
	if sys.Inb(com1+5)&0x01 == 0 {
		return -1
	}
	return int(sys.Inb(com1 + 0))
}

//go:nosplit
func WriteByte(ch byte) {
	const lstatus = uint16(5)
	for {
		ret := sys.Inb(com1 + lstatus)
		if ret&0x20 != 0 {
			break
		}
	}
	sys.Outb(com1, uint8(ch))
}

//go:nosplit
func Write(s []byte) (int, error) {
	for _, c := range s {
		WriteByte(c)
	}
	return len(s), nil
}

//go:nosplit
func WriteString(s string) (int, error) {
	for i := 0; i < len(s); i++ {
		WriteByte(s[i])
	}
	return len(s), nil
}

//go:nosplit
func intr() {
	if inputCallback == nil {
		return
	}
	for {
		ch := ReadByte()
		if ch == -1 {
			break
		}
		inputCallback(byte(ch))
	}
	pic.EOI(_IRQ_COM1)
}

//go:nosplit
func PreInit() {
	sys.Outb(com1+3, 0x80) // unlock divisor
	sys.Outb(com1+0, 115200/9600)
	sys.Outb(com1+1, 0)

	sys.Outb(com1+3, 0x03) // lock divisor
	// disable fifo
	sys.Outb(com1+2, 0)

	// enable receive interrupt
	sys.Outb(com1+4, 0x00)
	sys.Outb(com1+1, 0x01)
}

func OnInput(callback func(byte)) {
	inputCallback = callback
}

func Init() {
	trap.Register(_IRQ_COM1, intr)
	pic.EnableIRQ(pic.LINE_COM1)
}
