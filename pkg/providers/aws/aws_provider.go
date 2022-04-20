package aws

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
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/sirupsen/logrus"
)

func AwsStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "aws", "state.json")
}

type AwsProvider struct {
	config config.Aws
	state  state.State
}

func NewAwsProvier(config config.Aws) *AwsProvider {
	logrus.Infof("state file: %s", AwsStateFile())
	return &AwsProvider{
		config: config,
		state:  state.NewBasicState(AwsStateFile()),
	}
}

func (p *AwsProvider) WithState(state state.State) *AwsProvider {
	p.state = state
	return p
}

func (p *AwsProvider) newEC2() *ec2.EC2 {
	sess := session.New(&aws.Config{
		Region: aws.String(p.config.Region),
	})
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		if r != nil {
			logrus.WithFields(logrus.Fields{"params": r.Params}).Debugf("request sent to EC2")
		}
	})
	return ec2.New(sess)
}

func (p *AwsProvider) newS3() *s3.S3 {
	sess := session.New(&aws.Config{
		Region: aws.String(p.config.Region),
	})
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		if r != nil {
			logrus.WithFields(logrus.Fields{"params": r.Params}).Debugf("request sent to S3")
		}
	})
	return s3.New(sess)
}
