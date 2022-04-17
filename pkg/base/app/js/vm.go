package js

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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
	"github.com/robertkrimen/otto"
)

func NewVM() *otto.Otto {
	vm := otto.New()
	addBuiltins(vm)
	return vm
}

func addBuiltins(vm *otto.Otto) {
	vm.Set("http", map[string]interface{}{
		"Get": func(url string) string {
			resp, err := http.Get(url)
			if err != nil {
				Throw(err)
			}
			defer resp.Body.Close()
			buf, _ := ioutil.ReadAll(resp.Body)
			return string(buf)
		},
	})
	vm.Set("sys", map[string]interface{}{
		"in8": func(port uint16) byte {
			return sys.Inb(port)
		},
		"out8": func(port uint16, data byte) {
			sys.Outb(port, data)
		},
	})
	vm.Set("printf", func(fmtstr string, args ...interface{}) int {
		n, _ := fmt.Printf(fmtstr, args...)
		return n
	})
}

// Throw throw go error in js vm as an Exception
func Throw(err error) {
	v, _ := otto.ToValue("Exception: " + err.Error())
	panic(v)
}

// Throws throw go string in js vm as an Exception
func Throws(msg string) {
	Throw(errors.New(msg))
}
