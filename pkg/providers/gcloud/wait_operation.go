package gcloud

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
	"time"

	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

var defaultTimeout = time.Minute * 5
var defaultInterval = time.Millisecond * 250

func (p *GcloudProvider) waitOperation(operation string, global bool) error {
	errc := make(chan error)
	finished := make(chan struct{})

	backoff := int64(1)
	go func() {
		for {
			done, err := p.waitCycle(operation, global)
			if err != nil {
				errc <- err
				return
			}
			if done {
				close(finished)
				return
			}
			backoff *= 2
			time.Sleep(time.Duration(backoff) * defaultInterval)
		}
	}()

	select {
	case err := <-errc:
		return err
	case <-finished:
		return nil
	case <-time.After(defaultTimeout):
		return errors.New("timed out waiting more than "+defaultTimeout.String()+" for "+operation+" to complete", nil)
	}
}

func (p *GcloudProvider) waitCycle(operation string, global bool) (bool, error) {
	var status *compute.Operation
	var err error
	if global {
		status, err = p.compute().GlobalOperations.Get(p.config.ProjectID, operation).Do()
	} else {
		status, err = p.compute().ZoneOperations.Get(p.config.ProjectID, p.config.Zone, operation).Do()
	}
	if err != nil {
		return false, errors.New("getting status for operation "+operation, err)
	}
	logrus.Debugf("status for %v is %+v", operation, status)
	if status.Status == "DONE" {
		return true, nil
	}
	return false, nil
}
