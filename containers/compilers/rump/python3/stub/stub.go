package main

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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// to exit the program with error use fatal method only!!

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func getEnvAmazon() (map[string]string, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: dialTimeout,
		},
	}
	resp, err := client.Get("http://169.254.169.254/latest/user-data")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env map[string]string
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return env, nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("Bhojpur Kernel v0.0 boostrapping beginning...")
	envChan := make(chan map[string]string)

	closeChan := make(chan struct{})

	go func() {
		listenerIp, err := getListenerIp(closeChan)
		if err != nil {
			log.Printf("err getting listener ip: %v", err)
			return
		}
		if env, err := registerWithListener(listenerIp); err != nil {
			log.Printf("err registering with listener: %v", err)
			return
		} else {
			envChan <- env
		}
	}()

	errChan := make(chan error)
	go func() {
		env, err := getEnvAmazon()
		envChan <- env
		errChan <- err
		close(closeChan)
	}()

envLoop:
	for {
		log.Printf("waiting for Bhojpur Kernel bootstrap")
		select {
		case env := <-envChan:
			if env != nil {
				log.Printf("env was set: %v", env)
				setEnv(env)
				break envLoop
			}
		case err := <-errChan:
			log.Printf("error: %v", err)
		}
	}
	log.Printf("continuing to main\n")
}

func setEnv(env map[string]string) error {
	for key, val := range env {
		os.Setenv(key, val)
	}
	return nil
}

func registerWithListener(listenerIp string) (map[string]string, error) {
	//get MAC Addr
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("retrieving network interfaces" + err.Error())
	}
	macAddress := ""
	for _, iface := range ifaces {
		log.Printf("found an interface: %v\n", iface)
		if len(iface.HardwareAddr) > 0 {
			macAddress = iface.HardwareAddr.String()
			break
		}
	}
	if macAddress == "" {
		return nil, errors.New("could not find mac address")
	}

	resp, err := http.Post("http://"+listenerIp+":3000/register?mac_address="+macAddress, "", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var env map[string]string
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return env, nil
}

func getListenerIp(closeChan <-chan struct{}) (string, error) {
	log.Printf("listening for udp heartbeat...")
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9967,
	})
	if err != nil {
		return "", err
	}
	for {
		log.Printf("listening...")
		data := make([]byte, 4096)
		_, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			return "", err
		}
		log.Printf("received an ip from %s with data: %s", remoteAddr.IP.String(), string(data))
		if strings.Contains(string(data), "kernctl") {
			data = bytes.Trim(data, "\x00")
			return strings.Split(string(data), ":")[1], nil
		}
		select {
		case <-closeChan:
			return "", nil //registered with ec2
		default:
			continue
		}
	}
}
