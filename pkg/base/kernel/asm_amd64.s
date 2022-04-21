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

// copy from go_tls.h
#define get_tls(r)	MOVL TLS, r
#define g(r)	0(r)(TLS*1)

// rt0 is the entry point of the kernel, which invokes kernel.preinit.
TEXT ·rt0(SB), NOSPLIT, $0-0
	// switch to new stack
	MOVQ $0x80000, SP
	XORQ BP, BP

	// DI and SI store multiboot magic and info passed by bootloader
	SUBQ $0x10, SP
	MOVQ DI, 0(SP)
	MOVQ SI, 8(SP)
	CALL ·preinit(SB)
	INT  $3

	// never return

// go_entry invokes _rt0_amd64_linux of the Go runtime.
TEXT ·go_entry(SB), NOSPLIT, $0
	SUBQ  $256, SP
	PUSHQ SP
	CALL  ·prepareArgs(SB)
	ADDQ  $8, SP
	JMP   _rt0_amd64_linux(SB)

// sseInit initializes the SSE instruction set.
TEXT ·sseInit(SB), NOSPLIT, $0
	MOVL CR0, AX
	ANDW $0xFFFB, AX
	ORW  $0x2, AX
	MOVL AX, CR0
	MOVL CR4, AX
	ORW  $3<<9, AX
	MOVL AX, CR4
	RET

// avxInit initializes the AVX instruction set.
TEXT ·avxInit(SB), NOSPLIT, $0
	// enable XGETBV and XSETBV
	MOVL CR4, AX
	ORL  $1<<18, AX
	MOVL AX, CR4

	// enable avx
	XORQ CX, CX
	XGETBV

	ORQ  $7, AX
	XORQ CX, CX
	XSETBV
	RET

// rdmsr(reg uint32) (ax, dx uint32) - Read from Model Specific Register.
TEXT ·rdmsr(SB), NOSPLIT, $0-16
	MOVL reg+0(FP), CX
	RDMSR
	MOVL AX, lo+8(FP)
	MOVL DX, hi+12(FP)
	RET

// wrmsr(reg, lo, hi uint32) - Write to Model Specific Register.
TEXT ·wrmsr(SB), NOSPLIT, $0-16
	MOVL reg+0(FP), CX
	MOVL lo+8(FP), AX
	MOVL hi+12(FP), DX
	WRMSR
	RET

// getg() uint64 - returns G from thread local storage of the current thread.
TEXT ·getg(SB), NOSPLIT, $0-8
	get_tls(CX)
	MOVQ g(CX), BX
	MOVQ BX, ret+0(FP)
	RET

// cpuid(fn, cx uint32) (ax, bx, cx, dx uint32) - CPU Identification.
TEXT ·cpuid(SB), NOSPLIT, $0-24
	MOVL fn+0(FP), AX
	MOVL cx+4(FP), CX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET
