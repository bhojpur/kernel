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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func (p *AwsProvider) ListInstances() ([]*types.Instance, error) {
	if len(p.state.GetInstances()) < 1 {
		return []*types.Instance{}, nil
	}

	instanceIds := []*string{}
	for instanceId := range p.state.GetInstances() {
		instanceIds = append(instanceIds, aws.String(instanceId))
	}
	param := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}
	output, err := p.newEC2().DescribeInstances(param)
	if err != nil {
		return nil, errors.New("running ec2 describe instances ", err)
	}
	updatedInstances := []*types.Instance{}
	for _, reservation := range output.Reservations {
		for _, ec2Instance := range reservation.Instances {
			logrus.WithField("ec2instance", ec2Instance).Debugf("aws returned instance %s", *ec2Instance.InstanceId)
			var instanceId string
			if ec2Instance.InstanceId != nil {
				instanceId = *ec2Instance.InstanceId
			}
			if instanceId == "" {
				logrus.Warnf("instance %v does not have readable instanceId, moving on", *ec2Instance)
				continue
			}
			instanceState := parseInstanceState(ec2Instance.State)
			if instanceState == types.InstanceState_Unknown {
				logrus.Warnf("instance %s state is unknown (%s), moving on", instanceId, *ec2Instance.State.Name)
				continue
			}
			if instanceState == types.InstanceState_Terminated {
				logrus.Warnf("instance %s state is terminated, removing it from state", instanceId)
				if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
					delete(instances, instanceId)
					return nil
				}); err != nil {
					return nil, errors.New("modifying instance map in state", err)
				}
				continue
			}
			instance, ok := p.state.GetInstances()[instanceId]
			if !ok {
				logrus.WithFields(logrus.Fields{"ec2Instance": ec2Instance}).Errorf("found an instance that Bhojpur Kernel has no record of")
				continue
			}
			instance.State = instanceState
			if ec2Instance.PublicIpAddress != nil {
				instance.IpAddress = *ec2Instance.PublicIpAddress
			}
			if err := p.state.ModifyInstances(func(instances map[string]*types.Instance) error {
				instances[instance.Id] = instance
				return nil
			}); err != nil {
				return nil, errors.New("modifying instance map in state", err)
			}
			updatedInstances = append(updatedInstances, instance)
		}
	}
	return updatedInstances, nil
}

func parseInstanceState(ec2State *ec2.InstanceState) types.InstanceState {
	if ec2State == nil {
		return types.InstanceState_Unknown
	}
	switch *ec2State.Name {
	case ec2.InstanceStateNameRunning:
		return types.InstanceState_Running
	case ec2.InstanceStateNamePending:
		return types.InstanceState_Pending
	case ec2.InstanceStateNameStopped:
		return types.InstanceState_Stopped
	case ec2.InstanceStateNameTerminated:
		return types.InstanceState_Terminated
	}
	return types.InstanceState_Unknown
}
