package cmd

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
	"strings"
	"time"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/bhojpur/kernel/pkg/base/drivers/kbd"
	"github.com/bhojpur/kernel/pkg/base/kernel"
)

func printstat(ctx *app.Context) {
	var stat1, stat2 [20]int64
	kernel.ThreadStat(&stat1)
	time.Sleep(time.Second)
	kernel.ThreadStat(&stat2)

	var sum int64
	for i := range stat1 {
		sum += stat2[i] - stat1[i]
	}
	var tids []string
	var percents []string
	for i := range stat1 {
		if stat1[i] == 0 {
			continue
		}
		tids = append(tids, fmt.Sprintf("%3d", i))
		percent := int(float32(stat2[i]-stat1[i]) / float32(sum) * 100)
		percents = append(percents, fmt.Sprintf("%3d", percent))
	}
	fmt.Fprintf(ctx.Stdout, "%s\n", strings.Join(tids, " "))
	fmt.Fprintf(ctx.Stdout, "%s\n\n", strings.Join(percents, " "))
}

func topmain(ctx *app.Context) error {
	for !kbd.Pressed('q') {
		printstat(ctx)
	}
	return nil
}

func init() {
	app.Register("top", topmain)
}
