package os

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
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/bhojpur/kernel/pkg/util/errors"
)

type DiskSize interface {
	ToPartedFormat() string
	ToBytes() Bytes
}

type Bytes int64

func (s Bytes) ToPartedFormat() string {
	return fmt.Sprintf("%dB", uint64(s))
}

func (s Bytes) ToBytes() Bytes {
	return s
}

// ToMegaBytes returns lowest whole number of size_MB so that size_MB >= (size_B / 1024^2)
func (s Bytes) ToMegaBytes() MegaBytes {
	return MegaBytes(int(math.Ceil(float64(s) / float64(MegaBytes(1).ToBytes()))))
}

type MegaBytes int64

func (s MegaBytes) ToPartedFormat() string {
	return fmt.Sprintf("%dMiB", uint64(s))
}

func (s MegaBytes) ToBytes() Bytes {
	return Bytes(s << 20)
}

type GigaBytes int64

func (s GigaBytes) ToPartedFormat() string {
	return fmt.Sprintf("%dGiB", uint64(s))
}

func (s GigaBytes) ToBytes() Bytes {
	return Bytes(s << 30)
}

type Sectors int64

const SectorSize = 512

func (s Sectors) ToPartedFormat() string {
	return fmt.Sprintf("%ds", uint64(s))
}

func (s Sectors) ToBytes() Bytes {
	return Bytes(s * SectorSize)
}

func ToSectors(b DiskSize) (Sectors, error) {
	inBytes := b.ToBytes()
	if inBytes%SectorSize != 0 {
		return 0, errors.New("can't convert to sectors", nil)
	}
	return Sectors(inBytes / SectorSize), nil
}

type BlockDevice string

func (b BlockDevice) Name() string {
	return string(b)
}

type Partitioner interface {
	MakeTable() error
	MakePart(partType string, start, size DiskSize) error
}

type Resource interface {
	Acquire() (BlockDevice, error)
	Release() error
}

type Part interface {
	Resource

	Size() DiskSize
	Offset() DiskSize

	Get() BlockDevice
}

func IsExists(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}

// ParseSize parses disk size string (e.g. "10GB" or "150MB") into MegaBytes
// If no unit string is provided, megabytes are assumed
func ParseSize(sizeStr string) (MegaBytes, error) {
	r, _ := regexp.Compile("^([0-9]+)(m|mb|M|MB|g|gb|G|GB)?$")
	match := r.FindStringSubmatch(sizeStr)
	if len(match) != 3 {
		return -1, fmt.Errorf("%s: unrecognized size", sizeStr)
	}
	size, _ := strconv.ParseInt(match[1], 10, 64)
	unit := match[2]
	switch unit {
	case "g", "gb", "G", "GB":
		size *= 1024
	}
	if size == 0 {
		return -1, fmt.Errorf("%s: size must be larger than zero", sizeStr)
	}
	return MegaBytes(size), nil
}
