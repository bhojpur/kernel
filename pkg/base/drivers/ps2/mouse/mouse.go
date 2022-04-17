package mouse

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
	"github.com/bhojpur/kernel/pkg/base/drivers/ps2"
	"github.com/bhojpur/kernel/pkg/base/kernel/trap"
)

const (
	_IRQ_MOUSE = pic.IRQ_BASE + pic.LINE_MOUSE
)

var (
	mouseCnt int

	packet     [3]byte
	status     byte
	xpos, ypos int

	eventch chan Packet
)

type Packet struct {
	X, Y        int
	Left, Right bool
}

func Cursor() (int, int) {
	return xpos, ypos
}

func LeftClick() bool {
	return status&0x01 != 0
}

func RightClick() bool {
	return status&0x02 != 0
}

func intr() {
	pic.EOI(_IRQ_MOUSE)
	for {
		st := ps2.ReadCmd()
		// log.Infof("status:%08b", st)
		if st&0x01 == 0 {
			break
		}
		x := ps2.ReadDataNoWait()
		// log.Infof("data:%08b", x)
		handlePacket(x)
	}
}

func handlePacket(v byte) {
	switch mouseCnt {
	case 0:
		packet[0] = v
		if v&0x08 == 0 {
			return
		}
		mouseCnt++
	case 1:
		packet[1] = v
		mouseCnt++
	case 2:
		packet[2] = v
		mouseCnt = 0
		// x overflow or y overflow, discard packet
		if packet[0]&0xC0 != 0 {
			return
		}
		status = packet[0]
		xpos += xrel(status, int(packet[1]))
		ypos -= yrel(status, int(packet[2]))

		p := Packet{
			X:     xpos,
			Y:     ypos,
			Left:  LeftClick(),
			Right: RightClick(),
		}
		select {
		case eventch <- p:
		default:
		}
	}
	// log.Infof("x:%d y:%d packet:%v status:%8b", xpos, ypos, packet, status)
}

func xrel(status byte, value int) int {
	var ret byte
	if status&0x10 != 0 {
		ret |= 0x80
	}
	ret |= byte(value)
	return int(int8(ret))
}

func yrel(status byte, value int) int {
	var ret byte
	if status&0x20 != 0 {
		ret |= 0x80
	}
	ret |= byte(value)
	return int(int8(ret))
}

func EventQueue() chan Packet {
	return eventch
}

func Init() {
	status := ps2.ReadCmd()
	// enable mouse IRQ and port clock
	status |= 0x22
	// enable keyboard IRQ and port clock
	status |= 0x11
	// enable keyboard translation
	status |= 0x40
	ps2.WriteCmd(0x60)
	ps2.WriteData(status, false)

	// enable ps2 mouse port
	ps2.WriteCmd(0xA8)

	// enable ps2 keyboard port
	ps2.WriteCmd(0xAE)

	// enable mouse send packet
	ps2.WriteMouseData(0xF4)

	trap.Register(_IRQ_MOUSE, intr)
	pic.EnableIRQ(pic.LINE_MOUSE)

	eventch = make(chan Packet, 10)
}
