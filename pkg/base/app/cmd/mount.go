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
	"errors"
	"net/url"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/bhojpur/kernel/pkg/base/fs"
	"github.com/bhojpur/kernel/pkg/base/fs/smb"
	"github.com/bhojpur/kernel/pkg/base/fs/stripprefix"
)

func mountmain(ctx *app.Context) error {
	if len(ctx.Args) < 3 {
		return errors.New("usage: mount $uri target")
	}
	uristr, target := ctx.Args[1], ctx.Args[2]
	uri, err := url.Parse(uristr)
	if err != nil {
		return err
	}
	switch uri.Scheme {
	case "smb":
		return mountsmb(uri, target)
	default:
		return errors.New("unsupported scheme " + uri.Scheme)
	}
}

func mountsmb(uri *url.URL, target string) error {
	passwd, _ := uri.User.Password()
	smbfs, err := smb.New(&smb.Config{
		Host:     uri.Host,
		User:     uri.User.Username(),
		Password: passwd,
		Mount:    uri.Path[1:],
	})
	if err != nil {
		return err
	}
	return fs.Mount(target, stripprefix.New("/", smbfs))
}

func init() {
	app.Register("mount", mountmain)
}
