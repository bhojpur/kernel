package utils_test

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

	"github.com/bhojpur/kernel/pkg/builder/utils"
)

func Example() {
	// Deps will run each dependency exactly once, and will run leaf-dependencies before those
	// functions that depend on them (if you put utils.Deps first in the function).

	// Normal (non-serial) Deps runs all dependencies in goroutines, so which one finishes first is
	// non-deterministic. Here we use SerialDeps here to ensure the example always produces the same
	// output.

	utils.SerialDeps(utils.F(Say, "hi"), Bark)
	// output:
	// hi
	// woof
}

func Say(something string) {
	fmt.Println(something)
}

func Bark() {
	utils.Deps(utils.F(Say, "woof"))
}
