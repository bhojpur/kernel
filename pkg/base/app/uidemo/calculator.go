package gui

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
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
)

type calcDemo struct {
	a, b    float64
	current *float64

	set      bool
	op, prev string

	editor nucular.TextEditor
}

var calcBtns = []string{
	"7", "8", "9", "+",
	"4", "5", "6", "-",
	"1", "2", "3", "*",
	"C", "0", "=", "/",
}

func (c *calcDemo) calculatorDemo(w *nucular.Window) {
	w.Row(35).Dynamic(1)
	c.editor.Flags = nucular.EditSimple
	c.editor.Filter = nucular.FilterFloat
	c.editor.Maxlen = 255
	c.editor.Buffer = []rune(fmt.Sprintf("%.2f", *c.current))
	c.editor.Edit(w)
	*c.current, _ = strconv.ParseFloat(string(c.editor.Buffer), 64)

	w.Row(35).Dynamic(4)
	solve := false
	for _, btn := range calcBtns {
		if w.Button(label.T(btn), false) {
			switch btn {
			case "+", "-", "*", "/":
				if !c.set {
					if c.current != &c.b {
						c.current = &c.b
					} else {
						c.prev = c.op
						solve = true
					}
				}
				c.op = btn
				c.set = true
			case "C":
				c.a = 0.0
				c.b = 0.0
				c.op = ""
				c.current = &c.a
				c.set = false
			case "=":
				solve = true
				c.prev = c.op
				c.op = ""
				c.set = false
			default:
				*c.current = *c.current*10 + float64(btn[0]-'0')
			}
		}
	}
	if solve {
		switch c.prev {
		case "+":
			c.a = c.a + c.b
		case "-":
			c.a = c.a - c.b
		case "*":
			c.a = c.a * c.b
		case "/":
			c.a = c.a / c.b
		}
		c.current = &c.a
		if c.set {
			c.current = &c.b
		}
		c.b = 0
		c.set = false
	}

}
