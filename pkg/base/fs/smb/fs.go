package smb

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
	"io/fs"
	"net"
	"os"

	"github.com/hirochachacha/go-smb2"
	"github.com/spf13/afero"
)

// assert that smb.Fs implements afero.Fs.
var _ afero.Fs = (*Fs)(nil)

type Config struct {
	Host     string
	User     string
	Password string
	Mount    string
}

type Fs struct {
	config Config
	conn   net.Conn
	sess   *smb2.Session
	*smb2.Share
}

func New(config *Config) (afero.Fs, error) {
	conn, err := net.Dial("tcp", config.Host)
	if err != nil {
		return nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.User,
			Password: config.Password,
		},
	}

	sess, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}
	share, err := sess.Mount(config.Mount)
	if err != nil {
		return nil, err
	}

	return &Fs{
		config: *config,
		sess:   sess,
		conn:   conn,
		Share:  share,
	}, nil
}

func (f *Fs) Close() {
	f.sess.Logoff()
	f.conn.Close()
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func (f *Fs) Create(name string) (afero.File, error) {
	return f.Share.Create(name)
}

// Open opens a file, returning it or an error, if any happens.
func (f *Fs) Open(name string) (afero.File, error) {
	return f.Share.Open(name)
}

// OpenFile opens a file using the given flags and the given mode.
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return f.Share.OpenFile(name, flag, perm)
}

// The name of this FileSystem
func (f *Fs) Name() string {
	return "smbfs"
}

// Chown changes the uid and gid of the named file.
func (f *Fs) Chown(name string, uid, gid int) error {
	// NOTE: go-smb2 doesn't implement the CAP_UNIX extensions, which are
	// required to handle Unix uid/guid.
	return fs.ErrInvalid
}
