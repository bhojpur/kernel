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

#include "textflag.h"

// pageEnable enables translation from virtual address (linear address) to
// physical address, based on the page directory set in the CR3 register.
TEXT ·pageEnable(SB), NOSPLIT, $0-0
	// enable PAE
	MOVQ CR4, AX
	BTSQ $5, AX
	MOVQ AX, CR4

	// enable page
	MOVQ CR0, AX
	BTSQ $31, AX
	MOVQ AX, CR0
	RET

// lcr3(topPage uint64) sets the CR3 register.
TEXT ·lcr3(SB), NOSPLIT, $0-8
	// setup page dir
	MOVQ topPage+0(FP), AX
	MOVQ AX, CR3
	RET
