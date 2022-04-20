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
	"flag"
	"path"

	log "github.com/sirupsen/logrus"

	"io"
	"os"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/pborman/uuid"
)

const staticFileDir = "/tmp/staticfiles"

func main() {
	log.SetLevel(log.DebugLevel)
	buildcontextdir := flag.String("d", "/opt/vol", "build context. relative volume names are relative to that")
	kernelInContext := flag.String("p", "program.bin", "kernel binary name.")
	usePartitionTables := flag.Bool("part", true, "indicates whether or not to use partition tables and install grub")
	args := flag.String("a", "", "arguments to kernel")
	out := flag.String("o", "", "base name of output file")

	flag.Parse()

	kernelFile := path.Join(*buildcontextdir, *kernelInContext)
	imgFile := path.Join(*buildcontextdir, "boot.image."+uuid.New())
	defer os.Remove(imgFile)

	log.WithFields(log.Fields{"kernelFile": kernelFile, "args": *args, "imgFile": imgFile, "usePartitionTables": *usePartitionTables}).Debug("calling CreateBootImageWithSize")

	s1, err := kos.DirSize(*buildcontextdir)
	if err != nil {
		log.Fatal(err)
	}
	s2 := float64(s1) * 1.1
	size := ((int64(s2) >> 20) + 20)

	if err := kos.CopyDir(*buildcontextdir, staticFileDir); err != nil {
		log.Fatal(err)
	}

	//no need to copy twice
	os.Remove(path.Join(staticFileDir, *kernelInContext))

	if err := kos.CreateBootImageWithSize(imgFile, kos.MegaBytes(size), kernelFile, staticFileDir, *args, *usePartitionTables); err != nil {
		log.Fatal(err)
	}

	src, err := os.Open(imgFile)
	if err != nil {
		log.Fatal("failed to open produced image file "+imgFile, err)
	}
	outFile := path.Join(*buildcontextdir, *out)
	dst, err := os.OpenFile(outFile, os.O_RDWR, 0)
	if err != nil {
		log.Fatal("failed to open target output file "+outFile, err)
	}
	n, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal("failed copying produced image file to target output file", err)
	}
	log.Infof("wrote %d bytes to disk", n)
}
