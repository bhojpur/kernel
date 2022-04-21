package mm

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

import "syscall"

const (
	// sync with kernel
	_SYS_FIXED_MMAP = 502
)

// SysMmap like Mmap but can run in user mode
// wraper of syscall.Mmap
func SysMmap(vaddr, size uintptr) uintptr {
	mem, _, err := syscall.Syscall6(syscall.SYS_MMAP, uintptr(vaddr), size, syscall.PROT_READ|syscall.PROT_WRITE, 0, 0, 0)
	if err != 0 {
		panic(err.Error())
	}
	return mem
}

// SysFixedMmap map the same physical address to the virtual address
// run in user mode
func SysFixedMmap(vaddr, paddr, size uintptr) {
	_, _, err := syscall.Syscall(_SYS_FIXED_MMAP, vaddr, paddr, size)
	if err != 0 {
		panic(err.Error())
	}
}
