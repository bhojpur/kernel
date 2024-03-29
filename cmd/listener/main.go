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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const statefile = "/data/statefile.json"

type state struct {
	MacIpMap  map[string]string            `json:"Ips"`
	MacEnvMap map[string]map[string]string `json:"Envs"`
}

func main() {
	args := os.Args
	for i, arg := range args {
		log.Printf("arg %v: %s", i, arg)
	}
	dataPrefix := flag.String("prefix", "bhojpur_", "prefix for data sent via udp (for identification purposes")
	enablePersistence := flag.Bool("enablePersistence", true, "assume a persistent volume is mounted to /data")
	flag.Parse()
	for i, arg := range flag.Args() {

		log.Printf("flagarg %v: %s", i, arg)
	}
	if *dataPrefix == "bhojpur_" {
		log.Printf("ERROR: must provide -prefix")
		return
	}
	if *dataPrefix == "" {
		log.Printf("ERROR: -prefix cannot be \"\"")
		return
	}
	if !*enablePersistence {
		os.MkdirAll("/data", 0755)
	}
	ipMapLock := sync.RWMutex{}
	envMapLock := sync.RWMutex{}
	saveLock := sync.Mutex{}
	var s state
	s.MacIpMap = make(map[string]string)
	s.MacEnvMap = make(map[string]map[string]string)

	data, err := ioutil.ReadFile(statefile)
	if err != nil {
		log.Printf("could not read statefile, maybe this is first boot: " + err.Error())
	} else {
		if err := json.Unmarshal(data, &s); err != nil {
			log.Printf("failed to parse state json: " + err.Error())
		}
	}

	listenerIp, listenerIpMask, err := getLocalIp()
	if err != nil {
		log.Printf("ERROR: failed to get local IP: %v", err)
		return
	}

	log.Printf("Starting Bhojpur Kernel discovery (udp heartbeat broadcast) with IP %s", listenerIp.String())
	info := []byte(*dataPrefix + ":" + listenerIp.String())
	BROADCAST_IPv4 := reverseMask(listenerIp, listenerIpMask)
	if listenerIpMask == nil {
		log.Printf("ERROR: listener-ip: %v; listener-ip-mask: %v; could not calculate broadcast address", listenerIp, listenerIpMask)
		return
	}
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 9967,
	})
	if err != nil {
		log.Printf(fmt.Sprintf("ERROR: broadcast-ip: %v; failed to dial udp broadcast connection", BROADCAST_IPv4))
		return
	}
	go func() {
		log.Printf("broadcasting...")
		for {
			_, err = socket.Write(info)
			if err != nil {
				log.Printf("ERROR: broadcast-ip: %v; failed writing to broadcast udp socket: "+err.Error(), BROADCAST_IPv4)
				return
			}
			time.Sleep(5000 * time.Millisecond)
		}
	}()
	m := http.NewServeMux()
	m.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		splitAddr := strings.Split(req.RemoteAddr, ":")
		if len(splitAddr) < 1 {
			log.Printf("req.RemoteAddr: %v, could not parse remote addr into ip/port combination", req.RemoteAddr)
			return
		}
		instanceIp := splitAddr[0]
		macAddress := req.URL.Query().Get("mac_address")
		log.Printf("Instance registered")
		log.Printf("ip: %v", instanceIp)
		log.Printf("ip: %v", macAddress)
		//mac address = the instance id in vsphere/vbox
		go func() {
			ipMapLock.Lock()
			defer ipMapLock.Unlock()
			s.MacIpMap[macAddress] = instanceIp
			go save(s, saveLock)
		}()
		envMapLock.RLock()
		defer envMapLock.RUnlock()
		env, ok := s.MacEnvMap[macAddress]
		if !ok {
			env = make(map[string]string)
			log.Printf("mac: %v", macAddress)
			log.Printf("env: %v", s.MacEnvMap)
			log.Printf("no env set for instance, replying with empty map")
		}
		data, err := json.Marshal(env)
		if err != nil {
			log.Printf("could not marshal env to json: " + err.Error())
			return
		}
		log.Printf("responding with data: %s", data)
		fmt.Fprintf(res, "%s", data)
	})
	m.HandleFunc("/set_instance_env", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		macAddress := req.URL.Query().Get("mac_address")
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		defer req.Body.Close()
		var env map[string]string
		if err := json.Unmarshal(data, &env); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		log.Printf("Env set for instance")
		log.Printf("mac: %v", macAddress)
		log.Printf("env: %v", env)
		envMapLock.Lock()
		defer envMapLock.Unlock()
		s.MacEnvMap[macAddress] = env
		go save(s, saveLock)
		res.WriteHeader(http.StatusAccepted)
	})
	m.HandleFunc("/instances", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		ipMapLock.RLock()
		defer ipMapLock.RUnlock()
		data, err := json.Marshal(s.MacIpMap)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
		}
		res.Write(data)
	})
	log.Printf("Bhojpur Kernel - Instance Listener serving on port 3000")
	http.ListenAndServe(":3000", m)
}

func getLocalIp() (net.IP, net.IPMask, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return net.IP{}, net.IPMask{}, errors.New("retrieving network interfaces" + err.Error())
	}
	for _, iface := range ifaces {
		log.Printf("found an interface: %v\n", iface)
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			log.Printf("inspecting address: %v", addr)
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() && v.IP.IsGlobalUnicast() && v.IP.To4() != nil && v.Mask != nil {
					return v.IP.To4(), v.Mask, nil
				}
			}
		}
	}
	return net.IP{}, net.IPMask{}, errors.New("failed to find IP on interfaces: " + fmt.Sprintf("%v", ifaces))
}

// ReverseMask returns the result of masking the IP address ip with mask.
func reverseMask(ip net.IP, mask net.IPMask) net.IP {
	n := len(ip)
	if n != len(mask) {
		return nil
	}
	out := make(net.IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] | (^mask[i])
	}
	return out
}

func save(s state, l sync.Mutex) {
	if err := func() error {
		l.Lock()
		defer l.Unlock()
		data, err := json.Marshal(s)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(statefile, data, 0644); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		log.Printf("failed to save state file %s", statefile)
	}
}
