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
	"regexp"
	"strings"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/peterh/liner"
	"github.com/robertkrimen/otto"
)

var lastExpressionRegex = regexp.MustCompile(`[a-zA-Z0-9]([a-zA-Z0-9\.]*[a-zA-Z0-9])?\.?$`)

func setAutoComplete(r app.LineReader, vm *otto.Otto) {
	r.SetAutoComplete(jsAutocompleteWrapper(vm))
}

func jsAutocompleteWrapper(vm *otto.Otto) liner.Completer {
	return func(line string) []string {
		return jsAutocomplete(vm, line)
	}
}
func jsAutocomplete(vm *otto.Otto, line string) []string {
	lastExpression := lastExpressionRegex.FindString(line)

	bits := strings.Split(lastExpression, ".")

	first := bits[:len(bits)-1]
	last := bits[len(bits)-1]

	var l []string

	if len(first) == 0 {
		c := vm.Context()

		l = make([]string, len(c.Symbols))

		i := 0
		for k := range c.Symbols {
			l[i] = k
			i++
		}
	} else {
		r, err := vm.Eval(strings.Join(bits[:len(bits)-1], "."))
		if err != nil {
			return nil
		}

		if o := r.Object(); o != nil {
			for _, v := range o.KeysByParent() {
				l = append(l, v...)
			}
		}
	}

	var r []string
	for _, s := range l {
		if strings.HasPrefix(s, last) {
			r = append(r, line+strings.TrimPrefix(s, last))
		}
	}

	return r
}
