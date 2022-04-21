package kernel

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

const (
	_CPUID_ECX_XSAVE = 1 << 26
	_CPUID_ECX_AVX   = 1 << 28
	_CPUID_EBX_AVX2  = 1 << 5

	_CPUID_FN_STD = 0x00000001
	_CPUID_FN_EXT = 0x80000001
)

//go:nosplit
func sseInit()

//go:nosplit
func avxInit()

//go:nosplit
func cpuid(fn, cx uint32) (eax, ebx, ecx, edx uint32)

//go:nosplit
func simdInit() {
	sseInit()

	// init for avx
	// first check avx function
	_, _, ecx, _ := cpuid(_CPUID_FN_STD, 0)
	if ecx&_CPUID_ECX_XSAVE == 0 {
		return
	}
	if ecx&_CPUID_ECX_AVX == 0 {
		return
	}
	_, ebx, _, _ := cpuid(0x0007, 0)
	if ebx&_CPUID_EBX_AVX2 == 0 {
		return
	}
	// all check passed, init avx
	avxInit()
}
