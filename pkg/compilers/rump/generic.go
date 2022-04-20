package rump

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
	"fmt"
	"regexp"
	"strings"

	"github.com/bhojpur/kernel/pkg/util/errors"

	"github.com/bhojpur/kernel/pkg/types"
	kutil "github.com/bhojpur/kernel/pkg/util"
)

func execContainer(imageName string, cmds []string, binds map[string]string, privileged bool, env map[string]string) error {
	container := kutil.NewContainer(imageName).Privileged(privileged).WithVolumes(binds).WithEnvs(env)
	if err := container.Run(cmds...); err != nil {
		return errors.New("running container "+imageName, err)
	}
	return nil
}

type RumCompilerBase struct {
	DockerImage string
	CreateImage func(kernel, args string, mntPoints, bakedEnv []string, noCleanup bool) (*types.RawImage, error)
}

func (r *RumCompilerBase) runContainer(localFolder string, envPairs []string) error {
	env := make(map[string]string)
	for _, pair := range envPairs {
		split := strings.Split(pair, "=")
		if len(split) != 2 {
			return errors.New(pair+" is invaid string for env pair", nil)
		}
		env[split[0]] = split[1]
	}

	if kutil.IsDockerToolbox() {
		localFolder = kutil.GetToolboxMountPath(localFolder)
	}
	return kutil.NewContainer(r.DockerImage).WithVolume(localFolder, "/opt/code").WithEnvs(env).Run()

}

func setRumpCmdLine(c rumpConfig, prog string, argv []string, addStub bool) rumpConfig {
	if addStub {
		stub := commandLine{
			Bin:  "stub",
			Argv: []string{},
		}
		c.Rc = append(c.Rc, stub)
	}
	progrc := commandLine{
		Bin:  "program",
		Argv: argv,
	}
	c.Rc = append(c.Rc, progrc)
	return c
}

var netRegEx = regexp.MustCompile("net[1-9]")
var envRegEx = regexp.MustCompile("\"env\":\\{(.*?)\\}")
var envRegEx2 = regexp.MustCompile("env[0-9]")

// rump special json
func toRumpJson(c rumpConfig) (string, error) {

	blk := c.Blk
	c.Blk = nil

	jsonConfig, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	blks := ""
	for _, b := range blk {

		blkjson, err := json.Marshal(b)
		if err != nil {
			return "", err
		}
		blks += fmt.Sprintf("\"blk\": %s,", string(blkjson))
	}
	var jsonString string
	if len(blks) > 0 {

		jsonString = string(jsonConfig[:len(jsonConfig)-1]) + "," + blks[:len(blks)-1] + "}"

	} else {
		jsonString = string(jsonConfig)
	}

	jsonString = netRegEx.ReplaceAllString(jsonString, "net")

	jsonString = string(envRegEx.ReplaceAllString(jsonString, "$1"))

	jsonString = string(envRegEx2.ReplaceAllString(jsonString, "env"))

	return jsonString, nil

}
