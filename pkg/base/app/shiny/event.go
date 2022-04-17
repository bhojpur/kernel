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
	"github.com/bhojpur/kernel/pkg/base/console"
	imouse "github.com/bhojpur/kernel/pkg/base/drivers/ps2/mouse"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/sys/unix"
)

func (w *windowImpl) listenKeyboardEvent() {
	termios := unix.Termios{}
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	unix.IoctlSetTermios(0, unix.TCSETS, &termios)

	buf := make([]byte, 16)
	for {
		n, _ := console.Console().Read(buf)
		content := buf[:n]
		for _, ch := range content {
			var code key.Code
			var char rune
			if ch == '\b' {
				code = key.CodeDeleteBackspace
			} else {
				char = rune(ch)
			}
			event := key.Event{
				Code:      code,
				Rune:      char,
				Direction: key.DirPress,
			}
			w.eventch <- event
		}
	}
}

func (w *windowImpl) listenMouseEvent() {
	for e := range imouse.EventQueue() {
		w.updateCursor()
		w.sendMouseEvent(e)
		w.cursor = e
	}
}

func (w *windowImpl) sendMouseEvent(e imouse.Packet) {
	var btn mouse.Button
	var dir mouse.Direction

	if e.Left {
		if !w.cursor.Left {
			btn = mouse.ButtonLeft
			dir = mouse.DirPress
		}
	} else {
		if w.cursor.Left {
			btn = mouse.ButtonLeft
			dir = mouse.DirRelease
		}
	}
	if e.Right {
		if !w.cursor.Right {
			btn = mouse.ButtonRight
			dir = mouse.DirPress
		}
	} else {
		if w.cursor.Right {
			btn = mouse.ButtonRight
			dir = mouse.DirRelease
		}
	}

	event := mouse.Event{
		X:         float32(e.X),
		Y:         float32(e.Y),
		Button:    btn,
		Direction: dir,
	}
	w.eventch <- event
}
