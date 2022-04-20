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
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bhojpur/kernel/pkg/util/errors"

	log "github.com/sirupsen/logrus"

	"io"

	kos "github.com/bhojpur/kernel/pkg/os"
	"github.com/pborman/uuid"
)

type volumeslice []kos.RawVolume

func (m *volumeslice) String() string {

	return fmt.Sprintf("%v", ([]kos.RawVolume)(*m))
}

// The second method is Set(value string) error
func (m *volumeslice) Set(value string) error {

	volparts := strings.Split(value, ",")

	if (len(volparts) != 1) && (len(volparts) != 2) {
		return errors.New("bad format", nil)
	}

	folder := volparts[0]

	var size int64
	if len(volparts) >= 2 {
		size, _ = strconv.ParseInt(volparts[1], 0, 64)
	}
	*m = append(*m, kos.RawVolume{Path: folder, Size: size})

	return nil
}

func verifyPreConditions() {
	_, err := os.Stat("/dev/loop0")
	if os.IsNotExist(err) {
		log.Fatal("No loop device found. if running from docker use \"--privileged -v /dev/:/dev/\"")
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	var volumes volumeslice
	partitionTable := flag.String("p", "true", "create partition table")
	buildcontextdir := flag.String("d", "/opt/vol", "build context. relative volume names are relative to that")
	volType := flag.String("t", "ext2", "type of volume 'mirage-fat', 'fat' or 'ext2'")
	flag.Var(&volumes, "v", "volumes folder[,size]")
	out := flag.String("o", "", "base name of output file")

	flag.Parse()

	if len(volumes) == 0 {
		log.Fatal("No volumes provided")
	}

	imgFile := path.Join(*buildcontextdir, "data.image."+uuid.New())
	defer os.Remove(imgFile)

	for i := range volumes {
		volumes[i].Path = path.Join(*buildcontextdir, volumes[i].Path)
	}

	verifyPreConditions()
	if *volType == "mirage-fat" {
		if *partitionTable == "true" {
			log.Fatal("Can't create mirage-fat volume with a partition table.")
		}

		if len(volumes) != 1 {
			log.Fatal("Can only create one volume with no partition table")
		}
		volume := volumes[0]
		if volume.Size == 0 {
			var err error
			volume.Size, err = kos.GetDirSize(volume.Path)
			if err != nil {
				log.Panic(err)
			}
		}

		sizeKb := volume.Size >> 10
		// add 8 kb to handle edge cases
		sizeKb += 8

		log.WithFields(log.Fields{"sizeKb": sizeKb, "size": volume.Size}).Info("Creating mirage fat volume.")

		err := kos.RunLogCommand("fat", "create", imgFile, fmt.Sprintf("%d%s", sizeKb, "KiB"))
		if err != nil {
			log.Panic(err)
		}

		if volume.Path != "" {
			err := kos.RunLogCommand("/bin/bash", "-c", fmt.Sprintf("cd \"%s\" && fat add %s *", volume.Path, imgFile))
			if err != nil {
				log.Panic(err)
			}
		}

	} else {

		if *partitionTable == "true" {
			log.Info("Creating volume with partition table")

			diskLabelGen := func(device string) kos.Partitioner { return &kos.DiskLabelPartioner{device} }

			// rump so we use disklabel
			err := kos.CreateVolumes(imgFile, *volType, []kos.RawVolume(volumes), diskLabelGen)

			if err != nil {
				panic(err)
			}
		} else {
			log.Info("Creating volume with no partition table")

			if len(volumes) != 1 {
				log.Fatal("Can only create one volume with no partition table")
			}

			err := kos.CreateSingleVolume(imgFile, *volType, volumes[0])

			if err != nil {
				panic(err)
			}
		}
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
