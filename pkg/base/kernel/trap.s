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

#define m_fpstate 32

TEXT alltraps(SB), NOSPLIT, $0
	PUSHQ R15
	PUSHQ R14
	PUSHQ R13
	PUSHQ R12
	PUSHQ R11
	PUSHQ R10
	PUSHQ R9
	PUSHQ R8
	PUSHQ DI
	PUSHQ SI
	PUSHQ BP
	PUSHQ DX
	PUSHQ CX
	PUSHQ BX
	PUSHQ AX

	// CX store mythread
	MOVQ   0(GS), CX
	MOVQ   m_fpstate(CX), DX
	FXSAVE (DX)

	// make top stack frame
	XORQ  BP, BP
	PUSHQ SP
	CALL  ·dotrap(SB)
	ADDQ  $8, SP
	JMP   ·trapret(SB)

TEXT ·trapret(SB), NOSPLIT, $0
	// CX store mythread
	MOVQ 0(GS), CX

	// restore FPU
	MOVQ    m_fpstate(CX), DX
	FXRSTOR (DX)

	POPQ AX
	POPQ BX
	POPQ CX
	POPQ DX
	POPQ BP
	POPQ SI
	POPQ DI
	POPQ R8
	POPQ R9
	POPQ R10
	POPQ R11
	POPQ R12
	POPQ R13
	POPQ R14
	POPQ R15

	ADDQ $16, SP // skip trapno and errcode

	IRETQ
