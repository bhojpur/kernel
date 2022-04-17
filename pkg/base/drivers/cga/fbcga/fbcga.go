package fbcga

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
	"image/color"
	"image/draw"

	"github.com/bhojpur/kernel/pkg/base/drivers/vbe"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

const (
	del = 0x7f
	bs  = '\b'

	twidth  = 80
	theight = 25
)

var (
	view   *vbe.View
	face   *basicfont.Face
	drawer *font.Drawer

	foreColor = color.RGBA{199, 199, 199, 0}
	backColor = color.Black

	Backend = fbbackend{}
)

type fbbackend struct {
	pos    int
	cursor int
	buffer [twidth * theight]byte
}

func (f *fbbackend) xy(pos int) (int, int) {
	x, y := pos%twidth*face.Advance, (pos/twidth+1)*face.Height
	if face.Left < 0 {
		x += -face.Left
	}
	return x, y
}

func (f *fbbackend) glyphPos(pos int) fixed.Point26_6 {
	x, y := f.xy(pos)
	return fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
}

func (f *fbbackend) drawChar(pos int, c byte, overide bool) {
	if c == 0 {
		return
	}
	fpos := f.glyphPos(pos)
	rect, _, _, _, _ := face.Glyph(fpos, rune(c))
	if overide {
		// clear old font
		draw.Draw(drawer.Dst, rect, image.NewUniform(backColor), image.Point{}, draw.Src)
	}
	drawer.Dot = fpos
	drawer.DrawBytes([]byte{c})
	view.CommitRect(rect)
}

func (f *fbbackend) setChar(pos int, c byte) {
	f.drawChar(pos, c, true)
	f.buffer[pos] = c
}

func (f *fbbackend) scrollup(pos int) int {
	copy(f.buffer[:], f.buffer[twidth:theight*twidth])
	pos -= twidth
	s := f.buffer[pos : theight*twidth]
	for i := range s {
		s[i] = 0
	}
	f.refresh()
	return pos
}

func (f *fbbackend) refresh() {
	rect := image.Rect(0, 0, twidth*face.Advance+face.Left, theight*face.Height+face.Descent)
	draw.Draw(drawer.Dst, rect, image.NewUniform(backColor), image.Point{}, draw.Src)

	for i := 0; i < theight; i++ {
		y := (i + 1) * face.Height
		drawer.Dot = fixed.Point26_6{X: fixed.I(0), Y: fixed.I(y)}
		drawer.DrawBytes(f.buffer[i*twidth : (i+1)*twidth])
	}
	view.CommitRect(rect)
}

func (f *fbbackend) updateCursor(n int) {
	old := f.cursor
	ch := f.buffer[old]
	f.drawChar(old, ch, true)

	f.drawChar(n, '_', false)
	f.cursor = n
}

func (f *fbbackend) SetPos(n int) {
	f.pos = n
	f.updateCursor(n)
}

func (f *fbbackend) GetPos() int {
	return f.pos
}

func (f *fbbackend) WritePos(n int, ch byte) {
	f.setChar(n, ch)
}

func (f *fbbackend) WriteByte(c byte) {
	pos := f.GetPos()
	switch c {
	case '\n', '\r':
		pos += twidth - pos%twidth
	case bs, del:
		if pos > 0 {
			f.setChar(pos, ' ')
			pos--
			f.setChar(pos, ' ')
		}
	default:
		f.setChar(pos, c)
		pos++
	}

	// Scroll up
	if pos/twidth >= theight {
		pos = f.scrollup(pos)
	}
	f.setChar(pos, ' ')
	f.SetPos(pos)
}

func Init() {
	if !vbe.IsEnable() {
		return
	}
	face = inconsolata.Regular8x16
	// face = inconsolata.Bold8x16
	view = vbe.DefaultView
	drawer = &font.Drawer{
		Dst:  view.Canvas(),
		Src:  image.NewUniform(foreColor),
		Face: face,
	}
}
