package inet

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
	"context"
	"errors"
	"time"

	"github.com/bhojpur/kernel/pkg/base/inet/dhcp"
	"github.com/bhojpur/kernel/pkg/base/log"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/link/loopback"
	"gvisor.dev/gvisor/pkg/tcpip/network/arp"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"
)

const (
	defaultNIC  = 1
	loopbackNIC = 2
)

var (
	nstack *stack.Stack
)

func e(err tcpip.Error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.String())
}

func Init() {
	nstack = stack.New(stack.Options{
		NetworkProtocols:   []stack.NetworkProtocolFactory{arp.NewProtocol, ipv4.NewProtocol},
		TransportProtocols: []stack.TransportProtocolFactory{tcp.NewProtocol, udp.NewProtocol},
		HandleLocal:        true,
	})

	// add net card interface
	endpoint := New(&Options{})
	err := nstack.CreateNIC(defaultNIC, endpoint)
	if err != nil {
		panic(err)
	}
	err1 := dodhcp(endpoint.LinkAddress())
	if err1 != nil {
		panic(err)
	}

	// add loopback interface
	err = nstack.CreateNIC(loopbackNIC, loopback.New())
	if err != nil {
		panic(err)
	}
	addInterfaceAddr(nstack, loopbackNIC, tcpip.Address([]byte{127, 0, 0, 1}))
	return
}

func addInterfaceAddr(s *stack.Stack, nic tcpip.NICID, addr tcpip.Address) {
	s.AddAddress(nic, ipv4.ProtocolNumber, addr)
	// Add route for local network if it doesn't exist already.
	localRoute := tcpip.Route{
		Destination: addr.WithPrefix().Subnet(),
		Gateway:     "", // No gateway for local network.
		NIC:         nic,
	}

	for _, rt := range s.GetRouteTable() {
		if rt.Equal(localRoute) {
			return
		}
	}

	// Local route does not exist yet. Add it.
	s.AddRoute(localRoute)
}

func dodhcp(linkaddr tcpip.LinkAddress) error {
	dhcpclient := dhcp.NewClient(nstack, defaultNIC, linkaddr)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	log.Infof("[inet] begin dhcp")
	err1 := dhcpclient.Request(ctx, "")
	cancel()
	if err1 != nil {
		return err1
	}
	log.Infof("[inet] dhcp done")
	cfg := dhcpclient.Config()
	log.Infof("[inet] addr:%v", dhcpclient.Address())
	log.Infof("[inet] gateway:%v", cfg.Gateway)
	log.Infof("[inet] mask:%v", cfg.SubnetMask)
	log.Infof("[inet] dns:%v", cfg.DomainNameServer)

	addInterfaceAddr(nstack, defaultNIC, dhcpclient.Address())
	nstack.AddRoute(tcpip.Route{
		Destination: header.IPv4EmptySubnet,
		Gateway:     cfg.Gateway,
		NIC:         defaultNIC,
	})
	return nil
}
