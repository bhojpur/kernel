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

import "math"

func X__builtin_isnan(t *TLS, x float64) int32    { return Bool32(math.IsNaN(x)) }
func Xacos(t *TLS, x float64) float64             { return math.Acos(x) }
func Xacosh(t *TLS, x float64) float64            { return math.Acosh(x) }
func Xasin(t *TLS, x float64) float64             { return math.Asin(x) }
func Xasinh(t *TLS, x float64) float64            { return math.Asinh(x) }
func Xatan(t *TLS, x float64) float64             { return math.Atan(x) }
func Xatan2(t *TLS, x, y float64) float64         { return math.Atan2(x, y) }
func Xatanh(t *TLS, x float64) float64            { return math.Atanh(x) }
func Xceil(t *TLS, x float64) float64             { return math.Ceil(x) }
func Xceilf(t *TLS, x float32) float32            { return float32(math.Ceil(float64(x))) }
func Xcopysign(t *TLS, x, y float64) float64      { return math.Copysign(x, y) }
func Xcopysignf(t *TLS, x, y float32) float32     { return float32(math.Copysign(float64(x), float64(y))) }
func Xcos(t *TLS, x float64) float64              { return math.Cos(x) }
func Xcosf(t *TLS, x float32) float32             { return float32(math.Cos(float64(x))) }
func Xcosh(t *TLS, x float64) float64             { return math.Cosh(x) }
func Xexp(t *TLS, x float64) float64              { return math.Exp(x) }
func Xfabs(t *TLS, x float64) float64             { return math.Abs(x) }
func Xfabsf(t *TLS, x float32) float32            { return float32(math.Abs(float64(x))) }
func Xfloor(t *TLS, x float64) float64            { return math.Floor(x) }
func Xfmod(t *TLS, x, y float64) float64          { return math.Mod(x, y) }
func Xhypot(t *TLS, x, y float64) float64         { return math.Hypot(x, y) }
func Xisnan(t *TLS, x float64) int32              { return X__builtin_isnan(t, x) }
func Xisnanf(t *TLS, x float32) int32             { return Bool32(math.IsNaN(float64(x))) }
func Xisnanl(t *TLS, x float64) int32             { return Bool32(math.IsNaN(x)) } // ccgo has to handle long double as double as Go does not support long double.
func Xldexp(t *TLS, x float64, exp int32) float64 { return math.Ldexp(x, int(exp)) }
func Xlog(t *TLS, x float64) float64              { return math.Log(x) }
func Xlog10(t *TLS, x float64) float64            { return math.Log10(x) }
func Xround(t *TLS, x float64) float64            { return math.Round(x) }
func Xsin(t *TLS, x float64) float64              { return math.Sin(x) }
func Xsinf(t *TLS, x float32) float32             { return float32(math.Sin(float64(x))) }
func Xsinh(t *TLS, x float64) float64             { return math.Sinh(x) }
func Xsqrt(t *TLS, x float64) float64             { return math.Sqrt(x) }
func Xtan(t *TLS, x float64) float64              { return math.Tan(x) }
func Xtanh(t *TLS, x float64) float64             { return math.Tanh(x) }
func Xtrunc(t *TLS, x float64) float64            { return math.Trunc(x) }

func Xpow(t *TLS, x, y float64) float64 {
	r := math.Pow(x, y)
	if x > 0 && r == 1 && y >= -1.0000000000000000715e-18 && y < -1e-30 {
		r = 0.9999999999999999
	}
	return r
}
