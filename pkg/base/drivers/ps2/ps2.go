package ps2

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

import "github.com/bhojpur/kernel/pkg/base/kernel/sys"

const (
	_CMD_PORT  = 0x64
	_DATA_PORT = 0x60
)

func waitCanWrite() {
	timeout := 1000
	for timeout > 0 {
		timeout--
		x := sys.Inb(_CMD_PORT)
		// input buffer empty means we can write to controller
		if x&0x02 == 0 {
			return
		}
	}
}

func waitCanRead() {
	timeout := 1000
	for timeout > 0 {
		timeout--
		x := sys.Inb(_CMD_PORT)
		// output buffer full means we can read from controller
		if x&0x01 != 0 {
			return
		}
	}
}

func ReadDataNoWait() byte {
	return sys.Inb(_DATA_PORT)
}

func ReadData() byte {
	waitCanRead()
	return sys.Inb(_DATA_PORT)
}

func WriteData(x byte, needAck bool) {
	waitCanWrite()
	sys.Outb(_DATA_PORT, x)
	if needAck {
		ReadAck()
	}
}

func ReadAck() {
	x := ReadData()
	if x != 0xFA {
		panic("not a ps2 ack packet")
	}
}

func ReadCmd() byte {
	return sys.Inb(_CMD_PORT)
}

func WriteCmd(x byte) {
	waitCanWrite()
	sys.Outb(_CMD_PORT, x)
}

func WriteMouseData(x byte) {
	WriteCmd(0xD4)
	WriteData(x, true)
}

func ReadMouseData(x byte) byte {
	WriteCmd(0xD4)
	WriteData(x, true)
	return ReadData()
}
