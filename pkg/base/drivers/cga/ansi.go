package cga

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
	"bytes"
	"errors"
)

const (
	stateBegin = iota
	stateESC
	stateLeft
	stateParam
	stateDone
)

const (
	_ESC = 0x1b
)

var (
	errInvalidChar = errors.New("invalid char")
	errNormalChar  = errors.New("normal char")
	errCSIDone     = errors.New("done")
)

type ansiParser struct {
	state int

	action   byte
	parambuf []byte
	params   []string
}

func (p *ansiParser) step(ch byte) error {
	switch p.state {
	case stateBegin:
		if ch != _ESC {
			return errNormalChar
		}
		p.state = stateESC
	case stateESC:
		if ch != '[' {
			return errInvalidChar
		}
		p.state = stateLeft
	case stateLeft:
		switch {
		case ch >= 0x30 && ch <= 0x3f:
			p.state = stateParam
			p.parambuf = append(p.parambuf, ch)
		case ch >= 0x40 && ch <= 0x7f:
			p.action = ch
			p.state = stateDone
		default:
			return errInvalidChar
		}
	case stateParam:
		switch {
		case ch >= 0x30 && ch <= 0x3f:
			p.parambuf = append(p.parambuf, ch)
		case ch >= 0x40 && ch <= 0x7f:
			p.action = ch
			p.state = stateDone
		default:
			return errInvalidChar
		}
	}
	if p.state == stateDone {
		p.decodePram()
		return errCSIDone
	}
	return nil
}

func (p *ansiParser) decodePram() {
	if len(p.parambuf) == 0 {
		return
	}
	params := bytes.Split(p.parambuf, []byte(";"))
	for _, param := range params {
		p.params = append(p.params, string(param))
	}
}

func (p *ansiParser) Action() byte {
	return p.action
}

func (p *ansiParser) Params() []string {
	return p.params
}

func (p *ansiParser) Reset() {
	p.state = stateBegin
	p.action = 0
	p.parambuf = p.parambuf[:0]
	p.params = p.params[:0]
}
