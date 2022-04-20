package qemu

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
	"net"
	"path/filepath"

	kutil "github.com/bhojpur/kernel/pkg/util"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func startDebuggerListener(port int) error {
	addr := fmt.Sprintf(":%v", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("establishing tcp listener on "+addr, err)
	}
	logrus.Info("listening on " + addr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logrus.WithError(err).Warnf("failed to accept debugger connection")
				continue
			}
			go connectDebugger(conn)
		}
	}()
	return nil
}

func connectDebugger(conn net.Conn) {
	if debuggerTargetImageName == "" {
		logrus.Error("no debug instance is currently running")
		return
	}
	container := kutil.NewContainer("rump-debugger-qemu").
		WithNet("host").
		WithVolume(filepath.Dir(getKernelPath(debuggerTargetImageName)), "/opt/prog/").
		Interactive(true)

	cmd := container.BuildCmd(
		"/opt/gdb-7.11/gdb/gdb",
		"-ex", "target remote 192.168.99.1:1234",
		"/opt/prog/program.bin",
	)
	conn.Read([]byte("GET / HTTP/1.0\r\n\r\n"))
	logrus.WithField("command", cmd.Args).Info("running debug command")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	if err := cmd.Start(); err != nil {
		logrus.WithError(err).Error("error starting debugger container")
		return
	}
	defer func() {
		//reset debugger target
		debuggerTargetImageName = ""
		container.Stop()
	}()

	for {
		if _, err := conn.Write([]byte{0}); err != nil {
			logrus.Debug("debugger disconnected: %v", err)
			return
		}
	}
}
