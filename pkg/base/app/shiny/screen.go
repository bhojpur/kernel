package shiny

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
	"image"

	"golang.org/x/exp/shiny/screen"
)

var (
	_ screen.Buffer = (*bufferImpl)(nil)
	_ screen.Window = (*windowImpl)(nil)

	defaultScreen screenImpl
)

type screenImpl struct {
}

// NewBuffer returns a new Buffer for this screen.
func (s *screenImpl) NewBuffer(size image.Point) (screen.Buffer, error) {
	m := image.NewRGBA(image.Rectangle{Max: size})
	return &bufferImpl{
		buf:  m.Pix,
		rgba: *m,
		size: size,
	}, nil
}

// NewTexture returns a new Texture for this screen.
func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	panic("not implemented") // TODO: Implement
}

// NewWindow returns a new Window for this screen.
//
// A nil opts is valid and means to use the default option values.
func (s *screenImpl) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return newWindow(), nil
}
