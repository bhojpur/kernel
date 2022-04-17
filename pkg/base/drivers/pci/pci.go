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

import (
	"github.com/bhojpur/kernel/pkg/base/drivers/pic"
	"github.com/bhojpur/kernel/pkg/base/kernel/trap"
	"github.com/bhojpur/kernel/pkg/base/log"
)

type Identity struct {
	Vendor uint16
	Device uint16
}

type Device struct {
	Ident Identity
	Addr  Address

	Class, SubClass uint8

	IRQLine uint8
	IRQNO   uint8
}

var devices []*Device

func Scan() []*Device {
	var devices []*Device
	for bus := int(0); bus < 256; bus++ {
		for dev := uint8(0); dev < 32; dev++ {
			for f := uint8(0); f < 8; f++ {
				addr := Address{
					Bus:    uint8(bus),
					Device: dev,
					Func:   f,
				}
				vendor := addr.ReadVendorID()
				if vendor == 0xffff {
					continue
				}
				devid := addr.ReadDeviceID()
				class := addr.ReadPCIClass()
				irqline := addr.ReadIRQLine()
				device := &Device{
					Ident: Identity{
						Vendor: vendor,
						Device: devid,
					},
					Addr:     addr,
					Class:    uint8((class >> 8) & 0xff),
					SubClass: uint8(class & 0xff),
					IRQLine:  irqline,
					IRQNO:    pic.IRQ_BASE + irqline,
				}
				devices = append(devices, device)
			}
		}
	}
	return devices
}

func findDev(idents []Identity) *Device {
	for _, ident := range idents {
		for _, dev := range devices {
			if dev.Ident == ident {
				return dev
			}
		}
	}
	return nil
}

func Init() {
	devices = Scan()
	for _, driver := range drivers {
		dev := findDev(driver.Idents())
		if dev == nil {
			log.Infof("[pci] no pci device found for %v\n", driver.Name())
			continue
		}
		log.Infof("[pci] found %x:%x for %s, irq:%d\n", dev.Ident.Vendor, dev.Ident.Device, driver.Name(), dev.IRQNO)
		driver.Init(dev)
		pic.EnableIRQ(uint16(dev.IRQLine))
		trap.Register(int(dev.IRQNO), driver.Intr)
	}
}
