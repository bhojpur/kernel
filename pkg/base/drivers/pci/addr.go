package pci

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
	configAddrPort = 0xcf8
	configDataPort = 0xcfc
)

type Address struct {
	Bus, Device, Func uint8
}

func (a Address) ReadBAR(bar uint8) (addr, len uint32, prefetch, isMem bool) {
	if bar > 0x5 {
		panic("invalid BAR")
	}
	reg := 0x10 + bar*4
	addr0 := a.ReadPCIRegister(reg)
	if addr0&1 != 0 {
		// I/O address.
		return uint32(addr0 &^ 0b11), 0, false, false
	}
	// Mask off flags.
	addr = uint32(addr0 &^ 0xf)
	switch (addr0 >> 1) & 0b11 {
	case 0b01:
		// 16-bit address. Not used.
		return 0, 0, false, false
	case 0b10:
		// 64-bit address.
		return 0, 0, false, false
	case 0b00:
		a.WritePCIRegister(reg, 0xffffffff)
		len = ^(a.ReadPCIRegister(reg) & 0xfffffff0) + 1
		a.WritePCIRegister(reg, addr0)
	}
	prefetch = addr0&0b1000 != 0
	return addr, len, prefetch, true
}

func (a Address) ReadCapOffset() uint8 {
	return uint8(a.ReadPCIRegister(0x34)) &^ 0x3
}

func (a Address) ReadStatus() uint16 {
	return uint16(a.ReadPCIRegister(0x4) >> 16)
}

func (a Address) ReadDeviceID() uint16 {
	return uint16(a.ReadPCIRegister(0x0) >> 16)
}

func (a Address) ReadVendorID() uint16 {
	return uint16(a.ReadPCIRegister(0x0))
}

func (a Address) readHeaderType() uint8 {
	return uint8(a.ReadPCIRegister(0xc) >> 16)
}

func (a Address) ReadPCIClass() uint16 {
	return uint16(a.ReadPCIRegister(0x8) >> 16)
}

func (a Address) ReadIRQLine() uint8 {
	return uint8(a.ReadPCIRegister(0x3C) & 0xff)
}

func (a Address) EnableBusMaster() {
	a.WritePCIRegister(0x04, a.ReadPCIRegister(0x04)|(1<<2))
}

func (a Address) readSecondaryBus() uint8 {
	return uint8(a.ReadPCIRegister(0x18) >> 8)
}

func (a Address) ReadPCIRegister(reg uint8) uint32 {
	if reg&0x3 != 0 {
		panic("unaligned PCI register access")
	}
	addr := 0x80000000 | uint32(a.Bus)<<16 | uint32(a.Device)<<11 | uint32(a.Func)<<8 | uint32(reg)
	sys.Outl(configAddrPort, addr)
	return sys.Inl(configDataPort)
}

func (a Address) WritePCIRegister(reg uint8, val uint32) {
	if reg&0x3 != 0 {
		panic("unaligned PCI register access")
	}
	addr := 0x80000000 | uint32(a.Bus)<<16 | uint32(a.Device)<<11 | uint32(a.Func)<<8 | uint32(reg)
	sys.Outl(configAddrPort, addr)
	sys.Outl(configDataPort, val)
}
