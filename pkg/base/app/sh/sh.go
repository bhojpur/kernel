package sh

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
	"log"
	"strings"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/bhojpur/kernel/pkg/base/console"

	"github.com/mattn/go-shellwords"
)

const prompt = "root@belaur# "

func main(ctx *app.Context) error {
	r := ctx.LineReader()
	defer r.Close()
	r.SetAutoComplete(autocompleteWrapper(ctx))
	for {
		line, err := r.Prompt(prompt)
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r.AppendHistory(line)
		err = doline(ctx, line)
		if err != nil {
			fmt.Fprintf(ctx.Stderr, "%s\n", err)
		}
	}
	fmt.Fprintf(ctx.Stdout, "exit\n")
	return nil
}

func doline(ctx *app.Context, line string) error {
	list, err := shellwords.Parse(line)
	if err != nil {
		return err
	}
	var bg bool
	if list[0] == "go" {
		bg = true
		list = list[1:]
	}
	name, args := list[0], list[1:]
	err = runApp(ctx, name, args, bg)
	if err != nil {
		return err
	}
	return nil
}

func runApp(ctx *app.Context, name string, args []string, bg bool) error {
	nctx := *ctx
	nctx.Args = append([]string{name}, args...)
	if bg {
		go func() {
			app.Run(name, &nctx)
			fmt.Fprintf(ctx.Stderr, "job %s done\n", name)
		}()
		return nil
	}
	return app.Run(name, &nctx)
}

func Bootstrap() {
	con := console.Console()
	log.SetOutput(con)
	ctx := &app.Context{
		Args:   []string{"sh"},
		Stdin:  con,
		Stdout: con,
		Stderr: con,
	}
	ctx.Init()
	main(ctx)
}

func init() {
	app.Register("sh", main)
}
