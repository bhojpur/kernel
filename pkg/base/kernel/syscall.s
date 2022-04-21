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

#define tls_my 0
#define tls_ax 8
#define m_kstack 40
#define ucode_idx  3
#define udata_idx  4
#define rpl_user   3

#define SYS_clockgettime 228

TEXT ·syscallEntry(SB), NOSPLIT, $0
	// save AX
	MOVQ AX, tls_ax(GS)

	// AX == pointer of current thread
	MOVQ tls_my(GS), AX

	// AX == kernel stack
	MOVQ m_kstack(AX), AX

	// push regs like INT 0x80
	SUBQ $40, AX

	// CX store IP
	MOVQ CX, 0(AX)

	// save CS
	MOVQ $ucode_idx<<3|rpl_user, 8(AX)

	// R11 store FLAGS
	MOVQ R11, 16(AX)

	// save SP
	MOVQ SP, 24(AX)

	// save SS
	MOVQ $udata_idx<<3|rpl_user, 32(AX)

	// change SP
	MOVQ AX, SP

	// restore AX
	MOVQ tls_ax(GS), AX

	// jmp INT 0x80
	JMP ·trap128(SB)

TEXT ·vdsoGettimeofday(SB), NOSPLIT, $0
	MOVQ $SYS_clockgettime, AX

	// DI store *TimeSpec, but clockgettime need SI
	MOVQ DI, SI
	MOVQ $0, DI
	INT  $0x80
	RET
