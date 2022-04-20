package common

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
	"net"
	"strings"
	"time"

	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

var socket *net.UDPConn

const BROADCAST_PORT = 9967

func GetInstanceListenerIp(dataPrefix string, timeout time.Duration) (string, error) {
	errc := make(chan error)
	go func() {
		<-time.After(timeout)
		errc <- errors.New("getting instance listener ip timed out after "+timeout.String(), nil)
	}()
	logrus.Infof("listening for udp heartbeat...")
	var err error
	//only initialize socket once
	logrus.Debug("ARE WE LISTENING ON THE SOCKET YET?", socket)
	if socket == nil {
		socket, err = net.ListenUDP("udp4", &net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: BROADCAST_PORT,
		})
		logrus.Debug("socket was", socket, "err was", err)
		if err != nil {
			return "", errors.New("opening udp socket", err)
		}
	}
	resultc := make(chan string)
	var stopLoop bool
	go func() {
		logrus.Infof("UDP Server listening on %s:%v", "0.0.0.0", BROADCAST_PORT)
		for !stopLoop {
			data := make([]byte, 4096)
			_, remoteAddr, err := socket.ReadFromUDP(data)
			if err != nil {
				errc <- errors.New("reading udp data", err)
				return
			}
			logrus.Infof("received an ip from %s with data: %s", remoteAddr.IP.String(), string(data))
			if strings.Contains(string(data), dataPrefix) {
				data = bytes.Trim(data, "\x00")
				resultc <- strings.Split(string(data), ":")[1]
				return
			}
		}
	}()
	select {
	case result := <-resultc:
		return result, nil
	case err := <-errc:
		stopLoop = true
		return "", err
	}
}
