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

// lgdt(gdtptr uint64) - Load Global Descriptor Table Register.
TEXT ·lgdt(SB), NOSPLIT, $0-8
	MOVQ gdtptr+0(FP), AX
	LGDT (AX)
	RET

// lidt(idtptr uint64) - Load Interrupt Descriptor Table Register.
TEXT ·lidt(SB), NOSPLIT, $0-8
	MOVQ idtptr+0(FP), AX
	LIDT (AX)
	RET

// ltr(sel uint64) - Load Task Register.
TEXT ·ltr(SB), NOSPLIT, $0-8
	MOVQ sel+0(FP), AX
	LTR  AX
	RET

// reloadCS returns from the current interrupt handler.
TEXT ·reloadCS(SB), NOSPLIT, $0
	// save ip
	MOVQ 0(SP), AX

	// save sp
	MOVQ SP, BX
	ADDQ $8, BX

	// rerange the stack, as in an interrupt stack
	PUSHQ $0x10 // SS
	PUSHQ BX
	PUSHFQ
	PUSHQ $8
	PUSHQ AX

	// IRET
	IRETQ
