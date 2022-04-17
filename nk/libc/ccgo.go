package libc

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

// Copyright 2020 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"sync/atomic"
	"unsafe"
)

func AtomicStoreNInt32(ptr uintptr, val int32, memorder int32) {
	atomic.StoreInt32((*int32)(unsafe.Pointer(ptr)), val)
}

func AtomicStoreNInt64(ptr uintptr, val int64, memorder int32) {
	atomic.StoreInt64((*int64)(unsafe.Pointer(ptr)), val)
}

func AtomicStoreNUint32(ptr uintptr, val uint32, memorder int32) {
	atomic.StoreUint32((*uint32)(unsafe.Pointer(ptr)), val)
}

func AtomicStoreNUint64(ptr uintptr, val uint64, memorder int32) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(ptr)), val)
}

func AtomicStoreNUintptr(ptr uintptr, val uintptr, memorder int32) {
	atomic.StoreUintptr((*uintptr)(unsafe.Pointer(ptr)), val)
}

func AtomicLoadNInt32(ptr uintptr, memorder int32) int32 {
	return atomic.LoadInt32((*int32)(unsafe.Pointer(ptr)))
}

func AtomicLoadNInt64(ptr uintptr, memorder int32) int64 {
	return atomic.LoadInt64((*int64)(unsafe.Pointer(ptr)))
}

func AtomicLoadNUint32(ptr uintptr, memorder int32) uint32 {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(ptr)))
}

func AtomicLoadNUint64(ptr uintptr, memorder int32) uint64 {
	return atomic.LoadUint64((*uint64)(unsafe.Pointer(ptr)))
}

func AtomicLoadNUintptr(ptr uintptr, memorder int32) uintptr {
	return atomic.LoadUintptr((*uintptr)(unsafe.Pointer(ptr)))
}
func AssignInt8(p *int8, v int8) int8             { *p = v; return v }
func AssignInt16(p *int16, v int16) int16         { *p = v; return v }
func AssignInt32(p *int32, v int32) int32         { *p = v; return v }
func AssignInt64(p *int64, v int64) int64         { *p = v; return v }
func AssignUint8(p *uint8, v uint8) uint8         { *p = v; return v }
func AssignUint16(p *uint16, v uint16) uint16     { *p = v; return v }
func AssignUint32(p *uint32, v uint32) uint32     { *p = v; return v }
func AssignUint64(p *uint64, v uint64) uint64     { *p = v; return v }
func AssignFloat32(p *float32, v float32) float32 { *p = v; return v }
func AssignFloat64(p *float64, v float64) float64 { *p = v; return v }
func AssignUintptr(p *uintptr, v uintptr) uintptr { *p = v; return v }

func AssignPtrInt8(p uintptr, v int8) int8          { *(*int8)(unsafe.Pointer(p)) = v; return v }
func AssignPtrInt16(p uintptr, v int16) int16       { *(*int16)(unsafe.Pointer(p)) = v; return v }
func AssignPtrInt32(p uintptr, v int32) int32       { *(*int32)(unsafe.Pointer(p)) = v; return v }
func AssignPtrInt64(p uintptr, v int64) int64       { *(*int64)(unsafe.Pointer(p)) = v; return v }
func AssignPtrUint8(p uintptr, v uint8) uint8       { *(*uint8)(unsafe.Pointer(p)) = v; return v }
func AssignPtrUint16(p uintptr, v uint16) uint16    { *(*uint16)(unsafe.Pointer(p)) = v; return v }
func AssignPtrUint32(p uintptr, v uint32) uint32    { *(*uint32)(unsafe.Pointer(p)) = v; return v }
func AssignPtrUint64(p uintptr, v uint64) uint64    { *(*uint64)(unsafe.Pointer(p)) = v; return v }
func AssignPtrFloat32(p uintptr, v float32) float32 { *(*float32)(unsafe.Pointer(p)) = v; return v }
func AssignPtrFloat64(p uintptr, v float64) float64 { *(*float64)(unsafe.Pointer(p)) = v; return v }
func AssignPtrUintptr(p uintptr, v uintptr) uintptr { *(*uintptr)(unsafe.Pointer(p)) = v; return v }

func AssignMulInt8(p *int8, v int8) int8             { *p *= v; return *p }
func AssignMulInt16(p *int16, v int16) int16         { *p *= v; return *p }
func AssignMulInt32(p *int32, v int32) int32         { *p *= v; return *p }
func AssignMulInt64(p *int64, v int64) int64         { *p *= v; return *p }
func AssignMulUint8(p *uint8, v uint8) uint8         { *p *= v; return *p }
func AssignMulUint16(p *uint16, v uint16) uint16     { *p *= v; return *p }
func AssignMulUint32(p *uint32, v uint32) uint32     { *p *= v; return *p }
func AssignMulUint64(p *uint64, v uint64) uint64     { *p *= v; return *p }
func AssignMulFloat32(p *float32, v float32) float32 { *p *= v; return *p }
func AssignMulFloat64(p *float64, v float64) float64 { *p *= v; return *p }
func AssignMulUintptr(p *uintptr, v uintptr) uintptr { *p *= v; return *p }

func AssignDivInt8(p *int8, v int8) int8             { *p /= v; return *p }
func AssignDivInt16(p *int16, v int16) int16         { *p /= v; return *p }
func AssignDivInt32(p *int32, v int32) int32         { *p /= v; return *p }
func AssignDivInt64(p *int64, v int64) int64         { *p /= v; return *p }
func AssignDivUint8(p *uint8, v uint8) uint8         { *p /= v; return *p }
func AssignDivUint16(p *uint16, v uint16) uint16     { *p /= v; return *p }
func AssignDivUint32(p *uint32, v uint32) uint32     { *p /= v; return *p }
func AssignDivUint64(p *uint64, v uint64) uint64     { *p /= v; return *p }
func AssignDivFloat32(p *float32, v float32) float32 { *p /= v; return *p }
func AssignDivFloat64(p *float64, v float64) float64 { *p /= v; return *p }
func AssignDivUintptr(p *uintptr, v uintptr) uintptr { *p /= v; return *p }

func AssignRemInt8(p *int8, v int8) int8             { *p %= v; return *p }
func AssignRemInt16(p *int16, v int16) int16         { *p %= v; return *p }
func AssignRemInt32(p *int32, v int32) int32         { *p %= v; return *p }
func AssignRemInt64(p *int64, v int64) int64         { *p %= v; return *p }
func AssignRemUint8(p *uint8, v uint8) uint8         { *p %= v; return *p }
func AssignRemUint16(p *uint16, v uint16) uint16     { *p %= v; return *p }
func AssignRemUint32(p *uint32, v uint32) uint32     { *p %= v; return *p }
func AssignRemUint64(p *uint64, v uint64) uint64     { *p %= v; return *p }
func AssignRemUintptr(p *uintptr, v uintptr) uintptr { *p %= v; return *p }

func AssignAddInt8(p *int8, v int8) int8             { *p += v; return *p }
func AssignAddInt16(p *int16, v int16) int16         { *p += v; return *p }
func AssignAddInt32(p *int32, v int32) int32         { *p += v; return *p }
func AssignAddInt64(p *int64, v int64) int64         { *p += v; return *p }
func AssignAddUint8(p *uint8, v uint8) uint8         { *p += v; return *p }
func AssignAddUint16(p *uint16, v uint16) uint16     { *p += v; return *p }
func AssignAddUint32(p *uint32, v uint32) uint32     { *p += v; return *p }
func AssignAddUint64(p *uint64, v uint64) uint64     { *p += v; return *p }
func AssignAddFloat32(p *float32, v float32) float32 { *p += v; return *p }
func AssignAddFloat64(p *float64, v float64) float64 { *p += v; return *p }
func AssignAddUintptr(p *uintptr, v uintptr) uintptr { *p += v; return *p }

func AssignSubInt8(p *int8, v int8) int8             { *p -= v; return *p }
func AssignSubInt16(p *int16, v int16) int16         { *p -= v; return *p }
func AssignSubInt32(p *int32, v int32) int32         { *p -= v; return *p }
func AssignSubInt64(p *int64, v int64) int64         { *p -= v; return *p }
func AssignSubUint8(p *uint8, v uint8) uint8         { *p -= v; return *p }
func AssignSubUint16(p *uint16, v uint16) uint16     { *p -= v; return *p }
func AssignSubUint32(p *uint32, v uint32) uint32     { *p -= v; return *p }
func AssignSubUint64(p *uint64, v uint64) uint64     { *p -= v; return *p }
func AssignSubFloat32(p *float32, v float32) float32 { *p -= v; return *p }
func AssignSubFloat64(p *float64, v float64) float64 { *p -= v; return *p }
func AssignSubUintptr(p *uintptr, v uintptr) uintptr { *p -= v; return *p }

func AssignAndInt8(p *int8, v int8) int8             { *p &= v; return *p }
func AssignAndInt16(p *int16, v int16) int16         { *p &= v; return *p }
func AssignAndInt32(p *int32, v int32) int32         { *p &= v; return *p }
func AssignAndInt64(p *int64, v int64) int64         { *p &= v; return *p }
func AssignAndUint8(p *uint8, v uint8) uint8         { *p &= v; return *p }
func AssignAndUint16(p *uint16, v uint16) uint16     { *p &= v; return *p }
func AssignAndUint32(p *uint32, v uint32) uint32     { *p &= v; return *p }
func AssignAndUint64(p *uint64, v uint64) uint64     { *p &= v; return *p }
func AssignAndUintptr(p *uintptr, v uintptr) uintptr { *p &= v; return *p }

func AssignXorInt8(p *int8, v int8) int8             { *p ^= v; return *p }
func AssignXorInt16(p *int16, v int16) int16         { *p ^= v; return *p }
func AssignXorInt32(p *int32, v int32) int32         { *p ^= v; return *p }
func AssignXorInt64(p *int64, v int64) int64         { *p ^= v; return *p }
func AssignXorUint8(p *uint8, v uint8) uint8         { *p ^= v; return *p }
func AssignXorUint16(p *uint16, v uint16) uint16     { *p ^= v; return *p }
func AssignXorUint32(p *uint32, v uint32) uint32     { *p ^= v; return *p }
func AssignXorUint64(p *uint64, v uint64) uint64     { *p ^= v; return *p }
func AssignXorUintptr(p *uintptr, v uintptr) uintptr { *p ^= v; return *p }

func AssignOrInt8(p *int8, v int8) int8             { *p |= v; return *p }
func AssignOrInt16(p *int16, v int16) int16         { *p |= v; return *p }
func AssignOrInt32(p *int32, v int32) int32         { *p |= v; return *p }
func AssignOrInt64(p *int64, v int64) int64         { *p |= v; return *p }
func AssignOrUint8(p *uint8, v uint8) uint8         { *p |= v; return *p }
func AssignOrUint16(p *uint16, v uint16) uint16     { *p |= v; return *p }
func AssignOrUint32(p *uint32, v uint32) uint32     { *p |= v; return *p }
func AssignOrUint64(p *uint64, v uint64) uint64     { *p |= v; return *p }
func AssignOrUintptr(p *uintptr, v uintptr) uintptr { *p |= v; return *p }

func AssignMulPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) *= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignMulPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) *= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignMulPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) *= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignMulPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) *= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignMulPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) *= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignMulPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) *= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignMulPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) *= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignMulPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) *= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignMulPtrFloat32(p uintptr, v float32) float32 {
	*(*float32)(unsafe.Pointer(p)) *= v
	return *(*float32)(unsafe.Pointer(p))
}

func AssignMulPtrFloat64(p uintptr, v float64) float64 {
	*(*float64)(unsafe.Pointer(p)) *= v
	return *(*float64)(unsafe.Pointer(p))
}

func AssignMulPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) *= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignDivPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) /= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignDivPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) /= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignDivPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) /= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignDivPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) /= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignDivPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) /= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignDivPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) /= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignDivPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) /= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignDivPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) /= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignDivPtrFloat32(p uintptr, v float32) float32 {
	*(*float32)(unsafe.Pointer(p)) /= v
	return *(*float32)(unsafe.Pointer(p))
}

func AssignDivPtrFloat64(p uintptr, v float64) float64 {
	*(*float64)(unsafe.Pointer(p)) /= v
	return *(*float64)(unsafe.Pointer(p))
}

func AssignDivPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) /= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignRemPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) %= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignRemPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) %= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignRemPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) %= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignRemPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) %= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignRemPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) %= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignRemPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) %= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignRemPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) %= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignRemPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) %= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignRemPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) %= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignAddPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) += v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignAddPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) += v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignAddPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) += v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignAddPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) += v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignAddPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) += v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignAddPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) += v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignAddPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) += v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignAddPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) += v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignAddPtrFloat32(p uintptr, v float32) float32 {
	*(*float32)(unsafe.Pointer(p)) += v
	return *(*float32)(unsafe.Pointer(p))
}

func AssignAddPtrFloat64(p uintptr, v float64) float64 {
	*(*float64)(unsafe.Pointer(p)) += v
	return *(*float64)(unsafe.Pointer(p))
}

func AssignAddPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) += v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignSubPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) -= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignSubPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) -= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignSubPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) -= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignSubPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) -= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignSubPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) -= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignSubPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) -= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignSubPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) -= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignSubPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) -= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignSubPtrFloat32(p uintptr, v float32) float32 {
	*(*float32)(unsafe.Pointer(p)) -= v
	return *(*float32)(unsafe.Pointer(p))
}

func AssignSubPtrFloat64(p uintptr, v float64) float64 {
	*(*float64)(unsafe.Pointer(p)) -= v
	return *(*float64)(unsafe.Pointer(p))
}

func AssignSubPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) -= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignAndPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) &= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignAndPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) &= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignAndPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) &= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignAndPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) &= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignAndPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) &= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignAndPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) &= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignAndPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) &= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignAndPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) &= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignAndPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) &= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignXorPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) ^= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignXorPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) ^= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignXorPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) ^= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignXorPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) ^= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignXorPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) ^= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignXorPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) ^= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignXorPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) ^= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignXorPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) ^= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignXorPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) ^= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignOrPtrInt8(p uintptr, v int8) int8 {
	*(*int8)(unsafe.Pointer(p)) |= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignOrPtrInt16(p uintptr, v int16) int16 {
	*(*int16)(unsafe.Pointer(p)) |= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignOrPtrInt32(p uintptr, v int32) int32 {
	*(*int32)(unsafe.Pointer(p)) |= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignOrPtrInt64(p uintptr, v int64) int64 {
	*(*int64)(unsafe.Pointer(p)) |= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignOrPtrUint8(p uintptr, v uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) |= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignOrPtrUint16(p uintptr, v uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) |= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignOrPtrUint32(p uintptr, v uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) |= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignOrPtrUint64(p uintptr, v uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) |= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignOrPtrUintptr(p uintptr, v uintptr) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) |= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignShlPtrInt8(p uintptr, v int) int8 {
	*(*int8)(unsafe.Pointer(p)) <<= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignShlPtrInt16(p uintptr, v int) int16 {
	*(*int16)(unsafe.Pointer(p)) <<= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignShlPtrInt32(p uintptr, v int) int32 {
	*(*int32)(unsafe.Pointer(p)) <<= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignShlPtrInt64(p uintptr, v int) int64 {
	*(*int64)(unsafe.Pointer(p)) <<= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignShlPtrUint8(p uintptr, v int) uint8 {
	*(*uint8)(unsafe.Pointer(p)) <<= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignShlPtrUint16(p uintptr, v int) uint16 {
	*(*uint16)(unsafe.Pointer(p)) <<= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignShlPtrUint32(p uintptr, v int) uint32 {
	*(*uint32)(unsafe.Pointer(p)) <<= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignShlPtrUint64(p uintptr, v int) uint64 {
	*(*uint64)(unsafe.Pointer(p)) <<= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignShlPtrUintptr(p uintptr, v int) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) <<= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignShrPtrInt8(p uintptr, v int) int8 {
	*(*int8)(unsafe.Pointer(p)) >>= v
	return *(*int8)(unsafe.Pointer(p))
}

func AssignShrPtrInt16(p uintptr, v int) int16 {
	*(*int16)(unsafe.Pointer(p)) >>= v
	return *(*int16)(unsafe.Pointer(p))
}

func AssignShrPtrInt32(p uintptr, v int) int32 {
	*(*int32)(unsafe.Pointer(p)) >>= v
	return *(*int32)(unsafe.Pointer(p))
}

func AssignShrPtrInt64(p uintptr, v int) int64 {
	*(*int64)(unsafe.Pointer(p)) >>= v
	return *(*int64)(unsafe.Pointer(p))
}

func AssignShrPtrUint8(p uintptr, v int) uint8 {
	*(*uint8)(unsafe.Pointer(p)) >>= v
	return *(*uint8)(unsafe.Pointer(p))
}

func AssignShrPtrUint16(p uintptr, v int) uint16 {
	*(*uint16)(unsafe.Pointer(p)) >>= v
	return *(*uint16)(unsafe.Pointer(p))
}

func AssignShrPtrUint32(p uintptr, v int) uint32 {
	*(*uint32)(unsafe.Pointer(p)) >>= v
	return *(*uint32)(unsafe.Pointer(p))
}

func AssignShrPtrUint64(p uintptr, v int) uint64 {
	*(*uint64)(unsafe.Pointer(p)) >>= v
	return *(*uint64)(unsafe.Pointer(p))
}

func AssignShrPtrUintptr(p uintptr, v int) uintptr {
	*(*uintptr)(unsafe.Pointer(p)) >>= v
	return *(*uintptr)(unsafe.Pointer(p))
}

func AssignShlInt8(p *int8, v int) int8 { *p <<= v; return *p }

func AssignShlInt16(p *int16, v int) int16 { *p <<= v; return *p }

func AssignShlInt32(p *int32, v int) int32 { *p <<= v; return *p }

func AssignShlInt64(p *int64, v int) int64 { *p <<= v; return *p }

func AssignShlUint8(p *uint8, v int) uint8 { *p <<= v; return *p }

func AssignShlUint16(p *uint16, v int) uint16 { *p <<= v; return *p }

func AssignShlUint32(p *uint32, v int) uint32 { *p <<= v; return *p }

func AssignShlUint64(p *uint64, v int) uint64 { *p <<= v; return *p }

func AssignShlUintptr(p *uintptr, v int) uintptr { *p <<= v; return *p }

func AssignShrInt8(p *int8, v int) int8 { *p >>= v; return *p }

func AssignShrInt16(p *int16, v int) int16 { *p >>= v; return *p }

func AssignShrInt32(p *int32, v int) int32 { *p >>= v; return *p }

func AssignShrInt64(p *int64, v int) int64 { *p >>= v; return *p }

func AssignShrUint8(p *uint8, v int) uint8 { *p >>= v; return *p }

func AssignShrUint16(p *uint16, v int) uint16 { *p >>= v; return *p }

func AssignShrUint32(p *uint32, v int) uint32 { *p >>= v; return *p }

func AssignShrUint64(p *uint64, v int) uint64 { *p >>= v; return *p }

func AssignShrUintptr(p *uintptr, v int) uintptr { *p >>= v; return *p }

func PreIncInt8(p *int8, d int8) int8             { *p += d; return *p }
func PreIncInt16(p *int16, d int16) int16         { *p += d; return *p }
func PreIncInt32(p *int32, d int32) int32         { *p += d; return *p }
func PreIncInt64(p *int64, d int64) int64         { *p += d; return *p }
func PreIncUint8(p *uint8, d uint8) uint8         { *p += d; return *p }
func PreIncUint16(p *uint16, d uint16) uint16     { *p += d; return *p }
func PreIncUint32(p *uint32, d uint32) uint32     { *p += d; return *p }
func PreIncUint64(p *uint64, d uint64) uint64     { *p += d; return *p }
func PreIncFloat32(p *float32, d float32) float32 { *p += d; return *p }
func PreIncFloat64(p *float64, d float64) float64 { *p += d; return *p }
func PreIncUintptr(p *uintptr, d uintptr) uintptr { *p += d; return *p }

func PreIncAtomicInt32(p *int32, d int32) int32         { return atomic.AddInt32(p, d) }
func PreIncAtomicInt64(p *int64, d int64) int64         { return atomic.AddInt64(p, d) }
func PreIncAtomicUint32(p *uint32, d uint32) uint32     { return atomic.AddUint32(p, d) }
func PreIncAtomicUint64(p *uint64, d uint64) uint64     { return atomic.AddUint64(p, d) }
func PreIncAtomicUintptr(p *uintptr, d uintptr) uintptr { return atomic.AddUintptr(p, d) }

func PreDecInt8(p *int8, d int8) int8             { *p -= d; return *p }
func PreDecInt16(p *int16, d int16) int16         { *p -= d; return *p }
func PreDecInt32(p *int32, d int32) int32         { *p -= d; return *p }
func PreDecInt64(p *int64, d int64) int64         { *p -= d; return *p }
func PreDecUint8(p *uint8, d uint8) uint8         { *p -= d; return *p }
func PreDecUint16(p *uint16, d uint16) uint16     { *p -= d; return *p }
func PreDecUint32(p *uint32, d uint32) uint32     { *p -= d; return *p }
func PreDecUint64(p *uint64, d uint64) uint64     { *p -= d; return *p }
func PreDecFloat32(p *float32, d float32) float32 { *p -= d; return *p }
func PreDecFloat64(p *float64, d float64) float64 { *p -= d; return *p }
func PreDecUintptr(p *uintptr, d uintptr) uintptr { *p -= d; return *p }

func PreDecAtomicInt32(p *int32, d int32) int32         { return atomic.AddInt32(p, -d) }
func PreDecAtomicInt64(p *int64, d int64) int64         { return atomic.AddInt64(p, -d) }
func PreDecAtomicUint32(p *uint32, d uint32) uint32     { return atomic.AddUint32(p, -d) }
func PreDecAtomicUint64(p *uint64, d uint64) uint64     { return atomic.AddUint64(p, -d) }
func PreDecAtomicUintptr(p *uintptr, d uintptr) uintptr { return atomic.AddUintptr(p, -d) }

func PostIncInt8(p *int8, d int8) int8             { r := *p; *p += d; return r }
func PostIncInt16(p *int16, d int16) int16         { r := *p; *p += d; return r }
func PostIncInt32(p *int32, d int32) int32         { r := *p; *p += d; return r }
func PostIncInt64(p *int64, d int64) int64         { r := *p; *p += d; return r }
func PostIncUint8(p *uint8, d uint8) uint8         { r := *p; *p += d; return r }
func PostIncUint16(p *uint16, d uint16) uint16     { r := *p; *p += d; return r }
func PostIncUint32(p *uint32, d uint32) uint32     { r := *p; *p += d; return r }
func PostIncUint64(p *uint64, d uint64) uint64     { r := *p; *p += d; return r }
func PostIncFloat32(p *float32, d float32) float32 { r := *p; *p += d; return r }
func PostIncFloat64(p *float64, d float64) float64 { r := *p; *p += d; return r }
func PostIncUintptr(p *uintptr, d uintptr) uintptr { r := *p; *p += d; return r }

func PostIncAtomicInt32(p *int32, d int32) int32         { return atomic.AddInt32(p, d) - d }
func PostIncAtomicInt64(p *int64, d int64) int64         { return atomic.AddInt64(p, d) - d }
func PostIncAtomicUint32(p *uint32, d uint32) uint32     { return atomic.AddUint32(p, d) - d }
func PostIncAtomicUint64(p *uint64, d uint64) uint64     { return atomic.AddUint64(p, d) - d }
func PostIncAtomicUintptr(p *uintptr, d uintptr) uintptr { return atomic.AddUintptr(p, d) - d }

func PostDecInt8(p *int8, d int8) int8             { r := *p; *p -= d; return r }
func PostDecInt16(p *int16, d int16) int16         { r := *p; *p -= d; return r }
func PostDecInt32(p *int32, d int32) int32         { r := *p; *p -= d; return r }
func PostDecInt64(p *int64, d int64) int64         { r := *p; *p -= d; return r }
func PostDecUint8(p *uint8, d uint8) uint8         { r := *p; *p -= d; return r }
func PostDecUint16(p *uint16, d uint16) uint16     { r := *p; *p -= d; return r }
func PostDecUint32(p *uint32, d uint32) uint32     { r := *p; *p -= d; return r }
func PostDecUint64(p *uint64, d uint64) uint64     { r := *p; *p -= d; return r }
func PostDecFloat32(p *float32, d float32) float32 { r := *p; *p -= d; return r }
func PostDecFloat64(p *float64, d float64) float64 { r := *p; *p -= d; return r }
func PostDecUintptr(p *uintptr, d uintptr) uintptr { r := *p; *p -= d; return r }

func PostDecAtomicInt32(p *int32, d int32) int32         { return atomic.AddInt32(p, -d) + d }
func PostDecAtomicInt64(p *int64, d int64) int64         { return atomic.AddInt64(p, -d) + d }
func PostDecAtomicUint32(p *uint32, d uint32) uint32     { return atomic.AddUint32(p, -d) + d }
func PostDecAtomicUint64(p *uint64, d uint64) uint64     { return atomic.AddUint64(p, -d) + d }
func PostDecAtomicUintptr(p *uintptr, d uintptr) uintptr { return atomic.AddUintptr(p, -d) + d }

func Int8FromInt8(n int8) int8             { return int8(n) }
func Int8FromInt16(n int16) int8           { return int8(n) }
func Int8FromInt32(n int32) int8           { return int8(n) }
func Int8FromInt64(n int64) int8           { return int8(n) }
func Int8FromUint8(n uint8) int8           { return int8(n) }
func Int8FromUint16(n uint16) int8         { return int8(n) }
func Int8FromUint32(n uint32) int8         { return int8(n) }
func Int8FromUint64(n uint64) int8         { return int8(n) }
func Int8FromFloat32(n float32) int8       { return int8(n) }
func Int8FromFloat64(n float64) int8       { return int8(n) }
func Int8FromUintptr(n uintptr) int8       { return int8(n) }
func Int16FromInt8(n int8) int16           { return int16(n) }
func Int16FromInt16(n int16) int16         { return int16(n) }
func Int16FromInt32(n int32) int16         { return int16(n) }
func Int16FromInt64(n int64) int16         { return int16(n) }
func Int16FromUint8(n uint8) int16         { return int16(n) }
func Int16FromUint16(n uint16) int16       { return int16(n) }
func Int16FromUint32(n uint32) int16       { return int16(n) }
func Int16FromUint64(n uint64) int16       { return int16(n) }
func Int16FromFloat32(n float32) int16     { return int16(n) }
func Int16FromFloat64(n float64) int16     { return int16(n) }
func Int16FromUintptr(n uintptr) int16     { return int16(n) }
func Int32FromInt8(n int8) int32           { return int32(n) }
func Int32FromInt16(n int16) int32         { return int32(n) }
func Int32FromInt32(n int32) int32         { return int32(n) }
func Int32FromInt64(n int64) int32         { return int32(n) }
func Int32FromUint8(n uint8) int32         { return int32(n) }
func Int32FromUint16(n uint16) int32       { return int32(n) }
func Int32FromUint32(n uint32) int32       { return int32(n) }
func Int32FromUint64(n uint64) int32       { return int32(n) }
func Int32FromFloat32(n float32) int32     { return int32(n) }
func Int32FromFloat64(n float64) int32     { return int32(n) }
func Int32FromUintptr(n uintptr) int32     { return int32(n) }
func Int64FromInt8(n int8) int64           { return int64(n) }
func Int64FromInt16(n int16) int64         { return int64(n) }
func Int64FromInt32(n int32) int64         { return int64(n) }
func Int64FromInt64(n int64) int64         { return int64(n) }
func Int64FromUint8(n uint8) int64         { return int64(n) }
func Int64FromUint16(n uint16) int64       { return int64(n) }
func Int64FromUint32(n uint32) int64       { return int64(n) }
func Int64FromUint64(n uint64) int64       { return int64(n) }
func Int64FromFloat32(n float32) int64     { return int64(n) }
func Int64FromFloat64(n float64) int64     { return int64(n) }
func Int64FromUintptr(n uintptr) int64     { return int64(n) }
func Uint8FromInt8(n int8) uint8           { return uint8(n) }
func Uint8FromInt16(n int16) uint8         { return uint8(n) }
func Uint8FromInt32(n int32) uint8         { return uint8(n) }
func Uint8FromInt64(n int64) uint8         { return uint8(n) }
func Uint8FromUint8(n uint8) uint8         { return uint8(n) }
func Uint8FromUint16(n uint16) uint8       { return uint8(n) }
func Uint8FromUint32(n uint32) uint8       { return uint8(n) }
func Uint8FromUint64(n uint64) uint8       { return uint8(n) }
func Uint8FromFloat32(n float32) uint8     { return uint8(n) }
func Uint8FromFloat64(n float64) uint8     { return uint8(n) }
func Uint8FromUintptr(n uintptr) uint8     { return uint8(n) }
func Uint16FromInt8(n int8) uint16         { return uint16(n) }
func Uint16FromInt16(n int16) uint16       { return uint16(n) }
func Uint16FromInt32(n int32) uint16       { return uint16(n) }
func Uint16FromInt64(n int64) uint16       { return uint16(n) }
func Uint16FromUint8(n uint8) uint16       { return uint16(n) }
func Uint16FromUint16(n uint16) uint16     { return uint16(n) }
func Uint16FromUint32(n uint32) uint16     { return uint16(n) }
func Uint16FromUint64(n uint64) uint16     { return uint16(n) }
func Uint16FromFloat32(n float32) uint16   { return uint16(n) }
func Uint16FromFloat64(n float64) uint16   { return uint16(n) }
func Uint16FromUintptr(n uintptr) uint16   { return uint16(n) }
func Uint32FromInt8(n int8) uint32         { return uint32(n) }
func Uint32FromInt16(n int16) uint32       { return uint32(n) }
func Uint32FromInt32(n int32) uint32       { return uint32(n) }
func Uint32FromInt64(n int64) uint32       { return uint32(n) }
func Uint32FromUint8(n uint8) uint32       { return uint32(n) }
func Uint32FromUint16(n uint16) uint32     { return uint32(n) }
func Uint32FromUint32(n uint32) uint32     { return uint32(n) }
func Uint32FromUint64(n uint64) uint32     { return uint32(n) }
func Uint32FromFloat32(n float32) uint32   { return uint32(n) }
func Uint32FromFloat64(n float64) uint32   { return uint32(n) }
func Uint32FromUintptr(n uintptr) uint32   { return uint32(n) }
func Uint64FromInt8(n int8) uint64         { return uint64(n) }
func Uint64FromInt16(n int16) uint64       { return uint64(n) }
func Uint64FromInt32(n int32) uint64       { return uint64(n) }
func Uint64FromInt64(n int64) uint64       { return uint64(n) }
func Uint64FromUint8(n uint8) uint64       { return uint64(n) }
func Uint64FromUint16(n uint16) uint64     { return uint64(n) }
func Uint64FromUint32(n uint32) uint64     { return uint64(n) }
func Uint64FromUint64(n uint64) uint64     { return uint64(n) }
func Uint64FromFloat32(n float32) uint64   { return uint64(n) }
func Uint64FromFloat64(n float64) uint64   { return uint64(n) }
func Uint64FromUintptr(n uintptr) uint64   { return uint64(n) }
func Float32FromInt8(n int8) float32       { return float32(n) }
func Float32FromInt16(n int16) float32     { return float32(n) }
func Float32FromInt32(n int32) float32     { return float32(n) }
func Float32FromInt64(n int64) float32     { return float32(n) }
func Float32FromUint8(n uint8) float32     { return float32(n) }
func Float32FromUint16(n uint16) float32   { return float32(n) }
func Float32FromUint32(n uint32) float32   { return float32(n) }
func Float32FromUint64(n uint64) float32   { return float32(n) }
func Float32FromFloat32(n float32) float32 { return float32(n) }
func Float32FromFloat64(n float64) float32 { return float32(n) }
func Float32FromUintptr(n uintptr) float32 { return float32(n) }
func Float64FromInt8(n int8) float64       { return float64(n) }
func Float64FromInt16(n int16) float64     { return float64(n) }
func Float64FromInt32(n int32) float64     { return float64(n) }
func Float64FromInt64(n int64) float64     { return float64(n) }
func Float64FromUint8(n uint8) float64     { return float64(n) }
func Float64FromUint16(n uint16) float64   { return float64(n) }
func Float64FromUint32(n uint32) float64   { return float64(n) }
func Float64FromUint64(n uint64) float64   { return float64(n) }
func Float64FromFloat32(n float32) float64 { return float64(n) }
func Float64FromFloat64(n float64) float64 { return float64(n) }
func Float64FromUintptr(n uintptr) float64 { return float64(n) }
func UintptrFromInt8(n int8) uintptr       { return uintptr(n) }
func UintptrFromInt16(n int16) uintptr     { return uintptr(n) }
func UintptrFromInt32(n int32) uintptr     { return uintptr(n) }
func UintptrFromInt64(n int64) uintptr     { return uintptr(n) }
func UintptrFromUint8(n uint8) uintptr     { return uintptr(n) }
func UintptrFromUint16(n uint16) uintptr   { return uintptr(n) }
func UintptrFromUint32(n uint32) uintptr   { return uintptr(n) }
func UintptrFromUint64(n uint64) uintptr   { return uintptr(n) }
func UintptrFromFloat32(n float32) uintptr { return uintptr(n) }
func UintptrFromFloat64(n float64) uintptr { return uintptr(n) }
func UintptrFromUintptr(n uintptr) uintptr { return uintptr(n) }

func Int8(n int8) int8          { return n }
func Int16(n int16) int16       { return n }
func Int32(n int32) int32       { return n }
func Int64(n int64) int64       { return n }
func Uint8(n uint8) uint8       { return n }
func Uint16(n uint16) uint16    { return n }
func Uint32(n uint32) uint32    { return n }
func Uint64(n uint64) uint64    { return n }
func Float32(n float32) float32 { return n }
func Float64(n float64) float64 { return n }
func Uintptr(n uintptr) uintptr { return n }

func NegInt8(n int8) int8          { return -n }
func NegInt16(n int16) int16       { return -n }
func NegInt32(n int32) int32       { return -n }
func NegInt64(n int64) int64       { return -n }
func NegUint8(n uint8) uint8       { return -n }
func NegUint16(n uint16) uint16    { return -n }
func NegUint32(n uint32) uint32    { return -n }
func NegUint64(n uint64) uint64    { return -n }
func NegUintptr(n uintptr) uintptr { return -n }

func CplInt8(n int8) int8          { return ^n }
func CplInt16(n int16) int16       { return ^n }
func CplInt32(n int32) int32       { return ^n }
func CplInt64(n int64) int64       { return ^n }
func CplUint8(n uint8) uint8       { return ^n }
func CplUint16(n uint16) uint16    { return ^n }
func CplUint32(n uint32) uint32    { return ^n }
func CplUint64(n uint64) uint64    { return ^n }
func CplUintptr(n uintptr) uintptr { return ^n }

func BoolInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func BoolInt16(b bool) int16 {
	if b {
		return 1
	}
	return 0
}

func BoolInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func BoolInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func BoolUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func BoolUint16(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

func BoolUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func BoolUint64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func SetBitFieldPtr8Int8(p uintptr, v int8, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Int16(p uintptr, v int16, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Int32(p uintptr, v int32, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Int64(p uintptr, v int64, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Uint8(p uintptr, v uint8, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Uint16(p uintptr, v uint16, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Uint32(p uintptr, v uint32, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr8Uint64(p uintptr, v uint64, off int, mask uint8) {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
}

func SetBitFieldPtr16Int8(p uintptr, v int8, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Int16(p uintptr, v int16, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Int32(p uintptr, v int32, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Int64(p uintptr, v int64, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Uint8(p uintptr, v uint8, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Uint16(p uintptr, v uint16, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Uint32(p uintptr, v uint32, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr16Uint64(p uintptr, v uint64, off int, mask uint16) {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
}

func SetBitFieldPtr32Int8(p uintptr, v int8, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Int16(p uintptr, v int16, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Int32(p uintptr, v int32, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Int64(p uintptr, v int64, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Uint8(p uintptr, v uint8, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Uint16(p uintptr, v uint16, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Uint32(p uintptr, v uint32, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr32Uint64(p uintptr, v uint64, off int, mask uint32) {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
}

func SetBitFieldPtr64Int8(p uintptr, v int8, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Int16(p uintptr, v int16, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Int32(p uintptr, v int32, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Int64(p uintptr, v int64, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Uint8(p uintptr, v uint8, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Uint16(p uintptr, v uint16, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Uint32(p uintptr, v uint32, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func SetBitFieldPtr64Uint64(p uintptr, v uint64, off int, mask uint64) {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
}

func AssignBitFieldPtr8Int8(p uintptr, v int8, w, off int, mask uint8) int8 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	s := 8 - w
	return v << s >> s
}

func AssignBitFieldPtr8Int16(p uintptr, v int16, w, off int, mask uint8) int16 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	s := 16 - w
	return v << s >> s
}

func AssignBitFieldPtr8Int32(p uintptr, v int32, w, off int, mask uint8) int32 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	s := 32 - w
	return v << s >> s
}

func AssignBitFieldPtr8Int64(p uintptr, v int64, w, off int, mask uint8) int64 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	s := 64 - w
	return v << s >> s
}

func AssignBitFieldPtr16Int8(p uintptr, v int8, w, off int, mask uint16) int8 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	s := 8 - w
	return v << s >> s
}

func AssignBitFieldPtr16Int16(p uintptr, v int16, w, off int, mask uint16) int16 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	s := 16 - w
	return v << s >> s
}

func AssignBitFieldPtr16Int32(p uintptr, v int32, w, off int, mask uint16) int32 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	s := 32 - w
	return v << s >> s
}

func AssignBitFieldPtr16Int64(p uintptr, v int64, w, off int, mask uint16) int64 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	s := 64 - w
	return v << s >> s
}

func AssignBitFieldPtr32Int8(p uintptr, v int8, w, off int, mask uint32) int8 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	s := 8 - w
	return v << s >> s
}

func AssignBitFieldPtr32Int16(p uintptr, v int16, w, off int, mask uint32) int16 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	s := 16 - w
	return v << s >> s
}

func AssignBitFieldPtr32Int32(p uintptr, v int32, w, off int, mask uint32) int32 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	s := 32 - w
	return v << s >> s
}

func AssignBitFieldPtr32Int64(p uintptr, v int64, w, off int, mask uint32) int64 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	s := 64 - w
	return v << s >> s
}

func AssignBitFieldPtr64Int8(p uintptr, v int8, w, off int, mask uint64) int8 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	s := 8 - w
	return v << s >> s
}

func AssignBitFieldPtr64Int16(p uintptr, v int16, w, off int, mask uint64) int16 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	s := 16 - w
	return v << s >> s
}

func AssignBitFieldPtr64Int32(p uintptr, v int32, w, off int, mask uint64) int32 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	s := 32 - w
	return v << s >> s
}

func AssignBitFieldPtr64Int64(p uintptr, v int64, w, off int, mask uint64) int64 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	s := 64 - w
	return v << s >> s
}

func AssignBitFieldPtr8Uint8(p uintptr, v uint8, w, off int, mask uint8) uint8 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	return v & uint8(mask>>off)
}

func AssignBitFieldPtr8Uint16(p uintptr, v uint16, w, off int, mask uint8) uint16 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	return v & uint16(mask>>off)
}

func AssignBitFieldPtr8Uint32(p uintptr, v uint32, w, off int, mask uint8) uint32 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	return v & uint32(mask>>off)
}

func AssignBitFieldPtr8Uint64(p uintptr, v uint64, w, off int, mask uint8) uint64 {
	*(*uint8)(unsafe.Pointer(p)) = *(*uint8)(unsafe.Pointer(p))&^uint8(mask) | uint8(v<<off)&mask
	return v & uint64(mask>>off)
}

func AssignBitFieldPtr16Uint8(p uintptr, v uint8, w, off int, mask uint16) uint8 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	return v & uint8(mask>>off)
}

func AssignBitFieldPtr16Uint16(p uintptr, v uint16, w, off int, mask uint16) uint16 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	return v & uint16(mask>>off)
}

func AssignBitFieldPtr16Uint32(p uintptr, v uint32, w, off int, mask uint16) uint32 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	return v & uint32(mask>>off)
}

func AssignBitFieldPtr16Uint64(p uintptr, v uint64, w, off int, mask uint16) uint64 {
	*(*uint16)(unsafe.Pointer(p)) = *(*uint16)(unsafe.Pointer(p))&^uint16(mask) | uint16(v<<off)&mask
	return v & uint64(mask>>off)
}

func AssignBitFieldPtr32Uint8(p uintptr, v uint8, w, off int, mask uint32) uint8 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	return v & uint8(mask>>off)
}

func AssignBitFieldPtr32Uint16(p uintptr, v uint16, w, off int, mask uint32) uint16 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	return v & uint16(mask>>off)
}

func AssignBitFieldPtr32Uint32(p uintptr, v uint32, w, off int, mask uint32) uint32 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	return v & uint32(mask>>off)
}

func AssignBitFieldPtr32Uint64(p uintptr, v uint64, w, off int, mask uint32) uint64 {
	*(*uint32)(unsafe.Pointer(p)) = *(*uint32)(unsafe.Pointer(p))&^uint32(mask) | uint32(v<<off)&mask
	return v & uint64(mask>>off)
}

func AssignBitFieldPtr64Uint8(p uintptr, v uint8, w, off int, mask uint64) uint8 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	return v & uint8(mask>>off)
}

func AssignBitFieldPtr64Uint16(p uintptr, v uint16, w, off int, mask uint64) uint16 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	return v & uint16(mask>>off)
}

func AssignBitFieldPtr64Uint32(p uintptr, v uint32, w, off int, mask uint64) uint32 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	return v & uint32(mask>>off)
}

func AssignBitFieldPtr64Uint64(p uintptr, v uint64, w, off int, mask uint64) uint64 {
	*(*uint64)(unsafe.Pointer(p)) = *(*uint64)(unsafe.Pointer(p))&^uint64(mask) | uint64(v<<off)&mask
	return v & uint64(mask>>off)
}

func PostDecBitFieldPtr8Int8(p uintptr, d int8, w, off int, mask uint8) (r int8) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Int16(p uintptr, d int16, w, off int, mask uint8) (r int16) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Int32(p uintptr, d int32, w, off int, mask uint8) (r int32) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Int64(p uintptr, d int64, w, off int, mask uint8) (r int64) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Int8(p uintptr, d int8, w, off int, mask uint16) (r int8) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Int16(p uintptr, d int16, w, off int, mask uint16) (r int16) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Int32(p uintptr, d int32, w, off int, mask uint16) (r int32) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Int64(p uintptr, d int64, w, off int, mask uint16) (r int64) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Int8(p uintptr, d int8, w, off int, mask uint32) (r int8) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Int16(p uintptr, d int16, w, off int, mask uint32) (r int16) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Int32(p uintptr, d int32, w, off int, mask uint32) (r int32) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Int64(p uintptr, d int64, w, off int, mask uint32) (r int64) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Int8(p uintptr, d int8, w, off int, mask uint64) (r int8) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Int16(p uintptr, d int16, w, off int, mask uint64) (r int16) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Int32(p uintptr, d int32, w, off int, mask uint64) (r int32) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Int64(p uintptr, d int64, w, off int, mask uint64) (r int64) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Uint8(p uintptr, d uint8, w, off int, mask uint8) (r uint8) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Uint16(p uintptr, d uint16, w, off int, mask uint8) (r uint16) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Uint32(p uintptr, d uint32, w, off int, mask uint8) (r uint32) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr8Uint64(p uintptr, d uint64, w, off int, mask uint8) (r uint64) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Uint8(p uintptr, d uint8, w, off int, mask uint16) (r uint8) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Uint16(p uintptr, d uint16, w, off int, mask uint16) (r uint16) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Uint32(p uintptr, d uint32, w, off int, mask uint16) (r uint32) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr16Uint64(p uintptr, d uint64, w, off int, mask uint16) (r uint64) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Uint8(p uintptr, d uint8, w, off int, mask uint32) (r uint8) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Uint16(p uintptr, d uint16, w, off int, mask uint32) (r uint16) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Uint32(p uintptr, d uint32, w, off int, mask uint32) (r uint32) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr32Uint64(p uintptr, d uint64, w, off int, mask uint32) (r uint64) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Uint8(p uintptr, d uint8, w, off int, mask uint64) (r uint8) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Uint16(p uintptr, d uint16, w, off int, mask uint64) (r uint16) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Uint32(p uintptr, d uint32, w, off int, mask uint64) (r uint32) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostDecBitFieldPtr64Uint64(p uintptr, d uint64, w, off int, mask uint64) (r uint64) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r-d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Int8(p uintptr, d int8, w, off int, mask uint8) (r int8) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Int16(p uintptr, d int16, w, off int, mask uint8) (r int16) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Int32(p uintptr, d int32, w, off int, mask uint8) (r int32) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Int64(p uintptr, d int64, w, off int, mask uint8) (r int64) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Int8(p uintptr, d int8, w, off int, mask uint16) (r int8) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Int16(p uintptr, d int16, w, off int, mask uint16) (r int16) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Int32(p uintptr, d int32, w, off int, mask uint16) (r int32) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Int64(p uintptr, d int64, w, off int, mask uint16) (r int64) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Int8(p uintptr, d int8, w, off int, mask uint32) (r int8) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Int16(p uintptr, d int16, w, off int, mask uint32) (r int16) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Int32(p uintptr, d int32, w, off int, mask uint32) (r int32) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Int64(p uintptr, d int64, w, off int, mask uint32) (r int64) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Int8(p uintptr, d int8, w, off int, mask uint64) (r int8) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 8 - w
	r = int8(x0) & int8(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Int16(p uintptr, d int16, w, off int, mask uint64) (r int16) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 16 - w
	r = int16(x0) & int16(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Int32(p uintptr, d int32, w, off int, mask uint64) (r int32) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 32 - w
	r = int32(x0) & int32(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Int64(p uintptr, d int64, w, off int, mask uint64) (r int64) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	s := 64 - w
	r = int64(x0) & int64(mask) << s >> s
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Uint8(p uintptr, d uint8, w, off int, mask uint8) (r uint8) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Uint16(p uintptr, d uint16, w, off int, mask uint8) (r uint16) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Uint32(p uintptr, d uint32, w, off int, mask uint8) (r uint32) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr8Uint64(p uintptr, d uint64, w, off int, mask uint8) (r uint64) {
	x0 := *(*uint8)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint8)(unsafe.Pointer(p)) = x0&^uint8(mask) | uint8(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Uint8(p uintptr, d uint8, w, off int, mask uint16) (r uint8) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Uint16(p uintptr, d uint16, w, off int, mask uint16) (r uint16) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Uint32(p uintptr, d uint32, w, off int, mask uint16) (r uint32) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr16Uint64(p uintptr, d uint64, w, off int, mask uint16) (r uint64) {
	x0 := *(*uint16)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint16)(unsafe.Pointer(p)) = x0&^uint16(mask) | uint16(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Uint8(p uintptr, d uint8, w, off int, mask uint32) (r uint8) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Uint16(p uintptr, d uint16, w, off int, mask uint32) (r uint16) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Uint32(p uintptr, d uint32, w, off int, mask uint32) (r uint32) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr32Uint64(p uintptr, d uint64, w, off int, mask uint32) (r uint64) {
	x0 := *(*uint32)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint32)(unsafe.Pointer(p)) = x0&^uint32(mask) | uint32(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Uint8(p uintptr, d uint8, w, off int, mask uint64) (r uint8) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint8(x0) & uint8(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Uint16(p uintptr, d uint16, w, off int, mask uint64) (r uint16) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint16(x0) & uint16(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Uint32(p uintptr, d uint32, w, off int, mask uint64) (r uint32) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint32(x0) & uint32(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}

func PostIncBitFieldPtr64Uint64(p uintptr, d uint64, w, off int, mask uint64) (r uint64) {
	x0 := *(*uint64)(unsafe.Pointer(p))
	r = uint64(x0) & uint64(mask) >> off
	*(*uint64)(unsafe.Pointer(p)) = x0&^uint64(mask) | uint64(r+d)<<off&mask
	return r
}
