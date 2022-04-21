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

// Outb(port uint16, data byte)
TEXT ·Outb(SB), NOSPLIT, $0-3
	MOVW port+0(FP), DX
	MOVB data+2(FP), AX
	OUTB
	RET

// byte Inb(port uint16)
TEXT ·Inb(SB), NOSPLIT, $0-5
	MOVW port+0(FP), DX
	XORW AX, AX
	INB
	MOVB AX, ret+4(FP)
	RET

// Outl(port uint16, data uint32)
TEXT ·Outl(SB), NOSPLIT, $0-8
	MOVW port+0(FP), DX
	MOVL data+4(FP), AX
	OUTL
	RET

// uint32 Inl(port uint16)
TEXT ·Inl(SB), NOSPLIT, $0-8
	MOVW port+0(FP), DX
	INL
	MOVL AX, ret+4(FP)
	RET

// SetAX(val uint32)
TEXT ·SetAX(SB), NOSPLIT, $0-4
	MOVL val+0(FP), AX
	RET

// uint32 Flags()
TEXT ·Flags(SB), NOSPLIT, $0-4
	PUSHFL
	POPL AX
	MOVL AX, ret+0(FP)
	RET

// uint32 Cr2()
TEXT ·Cr2(SB), NOSPLIT, $0-4
	MOVL CR2, AX
	MOVL AX, ret+0(FP)
	RET

// Fxsave(addr uint32)
TEXT ·Fxsave(SB), NOSPLIT, $0-4
	MOVL   addr+0(FP), AX
	FXSAVE (AX)
	RET
