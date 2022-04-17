package console

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
	"io"
	"sync"
	"syscall"
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/drivers/cga"
	"github.com/bhojpur/kernel/pkg/base/drivers/kbd"
	"github.com/bhojpur/kernel/pkg/base/drivers/uart"
)

const (
	CON_BUFLEN = 128
)

type console struct {
	rawch chan byte

	buf     [CON_BUFLEN]byte
	r, w, e uint

	tios syscall.Termios

	mutex  sync.Mutex
	notify *sync.Cond

	wmutex sync.Mutex
}

var (
	con *console
)

func newConsole() *console {
	c := &console{
		tios: syscall.Termios{
			Lflag: syscall.ICANON | syscall.ECHO,
		},
	}
	c.notify = sync.NewCond(&c.mutex)
	return c
}

func ctrl(c byte) byte {
	return c - '@'
}

//go:nosplit
func (c *console) intr(ch byte) {
	c.handleInput(ch)
}

func (c *console) rawmode() bool {
	return c.tios.Lflag&syscall.ICANON == 0
}

func (c *console) handleRaw(ch byte) {
	if c.e-c.r >= CON_BUFLEN {
		return
	}
	idx := c.e % CON_BUFLEN
	c.e++
	c.buf[idx] = byte(ch)
	c.w = c.e
	c.notify.Broadcast()
}

func (c *console) handleInput(ch byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.rawmode() {
		c.handleRaw(ch)
		return
	}

	switch ch {
	case 0x7f, ctrl('H'):
		if c.e > c.w {
			c.e--
			c.putc(0x7f)
		}
		return
	}

	if c.e-c.r >= CON_BUFLEN {
		return
	}
	if ch == '\r' {
		ch = '\n'
	}
	idx := c.e % CON_BUFLEN
	c.e++
	c.buf[idx] = byte(ch)
	c.putc(ch)
	if ch == '\n' || c.e == c.r+CON_BUFLEN {
		c.w = c.e
		c.notify.Broadcast()
	}
}

func (c *console) loop() {
	for ch := range c.rawch {
		c.handleInput(ch)
	}
}

func (c *console) putc(ch byte) {
	uart.WriteByte(ch)
	cga.WriteByte(ch)
}

func (c *console) read(p []byte) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i := 0
	for i < len(p) {
		for c.r == c.w {
			c.notify.Wait()
		}
		idx := c.r
		c.r++
		ch := c.buf[idx%CON_BUFLEN]
		p[i] = byte(ch)
		i++
		if ch == '\n' || c.rawmode() {
			break
		}
	}
	return i
}

func (c *console) Read(p []byte) (int, error) {
	return c.read(p), nil
}

func (c *console) Write(p []byte) (int, error) {
	c.wmutex.Lock()
	defer c.wmutex.Unlock()
	for _, ch := range p {
		c.putc(ch)
	}
	return len(p), nil
}

func (c *console) Ioctl(op, arg uintptr) error {
	switch op {
	case syscall.TIOCGWINSZ:
		w := (*winSize)(unsafe.Pointer(arg))
		w.row = 25
		w.col = 80
		return nil
	case syscall.TCGETS:
		tios := (*syscall.Termios)(unsafe.Pointer(arg))
		*tios = c.tios
		return nil
	case syscall.TCSETS:
		tios := (*syscall.Termios)(unsafe.Pointer(arg))
		c.tios = *tios
		return nil

	default:
		return syscall.EINVAL
	}
}

type winSize struct {
	row, col       uint16
	xpixel, ypixel uint16
}

var cononce sync.Once

func Console() io.ReadWriter {
	cononce.Do(func() {
		con = newConsole()
	})
	return con
}

func Init() {
	cononce.Do(func() {
		con = newConsole()
	})
	uart.OnInput(con.intr)
	kbd.OnInput(con.intr)
}
