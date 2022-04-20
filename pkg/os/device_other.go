//go:build !linux
// +build !linux

package os

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

func Mount(device BlockDevice) (mntpoint string, err error) {
	panic("Not supported")
}
func MountDevice(device string) (mntpoint string, err error) {
	panic("Not supported")
}

func Umount(point string) error {
	panic("Not supported")
}

type MsDosPartioner struct {
	Device string
}

func (m *MsDosPartioner) MakeTable() error {
	panic("Not supported")
}

func (m *MsDosPartioner) MakePart(partType string, start, size DiskSize) error {
	panic("Not supported")
	return nil
}
func (m *MsDosPartioner) MakePartTillEnd(partType string, start DiskSize) error {
	panic("Not supported")
	return nil
}

func (m *MsDosPartioner) Makebootable(partnum int) error {
	panic("Not supported")
	return nil
}

type DiskLabelPartioner struct {
	Device string
}

func (m *DiskLabelPartioner) MakeTable() error {
	panic("Not supported")
	return nil
}

func (m *DiskLabelPartioner) MakePart(partType string, start, size DiskSize) error {
	panic("Not supported")
	return nil
}

func ListParts(device BlockDevice) ([]Part, error) {
	panic("Not supported")
	return nil, nil
}

type PartedPart struct {
	Device BlockDevice
}

func (p *PartedPart) Size() DiskSize {
	panic("Not supported")
	return Bytes(0)
}
func (p *PartedPart) Offset() DiskSize {
	panic("Not supported")
	return Bytes(0)
}

func (p *PartedPart) Acquire() (BlockDevice, error) {

	panic("Not supported")
	return "", nil
}

func (p *PartedPart) Release() error {
	panic("Not supported")
	return nil
}

func (p *PartedPart) Get() BlockDevice {
	panic("Not supported")
	return ""
}

type DeviceMapperDevice struct {
	DeviceName string
}

func NewDevice(start, size Sectors, origDevice BlockDevice, deivceName string) Resource {
	panic("Not supported")
	return nil
}

func (p *DeviceMapperDevice) Size() DiskSize {
	panic("Not supported")
	return Bytes(0)
}
func (p *DeviceMapperDevice) Offset() DiskSize {
	panic("Not supported")
	return Bytes(0)
}

func (p *DeviceMapperDevice) Acquire() (BlockDevice, error) {

	panic("Not supported")
	return "", nil
}

func (p *DeviceMapperDevice) Release() error {

	panic("Not supported")
	return nil
}

func (p *DeviceMapperDevice) Get() BlockDevice {

	panic("Not supported")
	return ""
}

type LoDevice struct {
}

func NewLoDevice(device string) Resource {

	panic("Not supported")
	return nil
}

func (p *LoDevice) Acquire() (BlockDevice, error) {

	panic("Not supported")
	return "", nil
}

func (p *LoDevice) Release() error {

	panic("Not supported")
	return nil
}
