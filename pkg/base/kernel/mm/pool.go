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

import (
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
)

//go:notinheap
type memblk struct {
	next uintptr
}

// Pool used to manage fixed size memory block
//go:notinheap
type Pool struct {
	size uintptr
	head uintptr
}

// size will align ptr size
//go:nosplit
func PoolInit(p *Pool, size uintptr) {
	const align = sys.PtrSize - 1
	size = (size + align) &^ align
	p.size = size
}

//go:nosplit
func (p *Pool) grow() {
	start := kmm.alloc()
	end := start + PGSIZE
	for v := start; v+p.size <= end; v += p.size {
		p.Free(v)
	}
}

//go:nosplit
func (p *Pool) Alloc() uintptr {
	if p.head == 0 {
		p.grow()
	}
	ret := p.head
	h := (*memblk)(unsafe.Pointer(p.head))
	p.head = h.next
	sys.Memclr(ret, int(p.size))
	return ret
}

//go:nosplit
func (p *Pool) Free(ptr uintptr) {
	v := (*memblk)(unsafe.Pointer(ptr))
	v.next = p.head
	p.head = ptr
}
