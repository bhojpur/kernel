package fileutils

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
	"archive/tar"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/bhojpur/kernel/pkg/util/errors"
	log "github.com/sirupsen/logrus"
)

func ExtractTar(tarArchive io.ReadCloser, localFolder string) error {
	tr := tar.NewReader(tarArchive)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return errors.New("reading tar archive", err)
		}
		log.WithField("file", hdr.Name).Debug("Extracting file")
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(path.Join(localFolder, hdr.Name), 0755)
			if err != nil {
				return errors.New("making folder", err)
			}

		case tar.TypeReg:
			fallthrough
		case tar.TypeRegA:
			dir, _ := path.Split(hdr.Name)
			if err := os.MkdirAll(path.Join(localFolder, dir), 0755); err != nil {
				return errors.New("making parent folder for file", err)
			}

			outputFile, err := os.Create(path.Join(localFolder, hdr.Name))
			if err != nil {
				return errors.New("creating output file", err)
			}

			if _, err := io.Copy(outputFile, tr); err != nil {
				outputFile.Close()
				return errors.New("writing output file", err)
				return err
			}
			outputFile.Close()

		default:
			continue
		}
	}

	return nil
}

func Compress(source, destination string) error {
	tarCmd := exec.Command("tar", "cf", destination, "-C", source, ".")
	if out, err := tarCmd.Output(); err != nil {
		return errors.New("running tar command: "+string(out), err)
	}
	return nil
}
