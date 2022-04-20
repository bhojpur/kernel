//go:build !container-binary
// +build !container-binary

package util

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
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

var containerVersions map[string]string

func InitContainers() error {
	versionData, err := versiondata.Asset("containers/versions.json")
	if err != nil {
		return errors.New("failed to get version data from containers/versions.json", err)
	}
	if err := json.Unmarshal(versionData, &containerVersions); err != nil {
		return errors.New("failed to unmarshall version data "+string(versionData), err)
	}
	logrus.WithField("versions", containerVersions).Info("using container versions")
	return nil
}

type Container struct {
	env           map[string]string
	privileged    bool
	volumes       map[string]string
	interactive   bool
	network       string
	containerName string
	name          string
	entrypoint    string
}

func NewContainer(imageName string) *Container {
	c := &Container{}

	c.name = imageName
	c.env = make(map[string]string)
	c.volumes = make(map[string]string)

	return c
}

func (c *Container) WithEntrypoint(entrypoint string) *Container {
	c.entrypoint = entrypoint
	return c
}

func (c *Container) WithVolume(hostdir, containerdir string) *Container {
	c.volumes[hostdir] = containerdir
	return c
}

func (c *Container) WithVolumes(vols map[string]string) *Container {
	for k, v := range vols {
		c.WithVolume(k, v)
	}
	return c
}

func (c *Container) WithEnv(key, value string) *Container {
	c.env[key] = value
	return c
}

func (c *Container) WithEnvs(vars map[string]string) *Container {
	for k, v := range vars {
		c.WithEnv(k, v)
	}
	return c
}

func (c *Container) WithNet(net string) *Container {
	c.network = net
	return c
}

func (c *Container) WithName(name string) *Container {
	c.containerName = name
	return c
}

func (c *Container) Interactive(i bool) *Container {
	c.interactive = i
	return c
}

func (c *Container) Privileged(p bool) *Container {
	c.privileged = p
	return c
}

func (c *Container) Run(arguments ...string) error {
	cmd := c.BuildCmd(arguments...)

	LogCommand(cmd, true)

	return cmd.Run()
}

func (c *Container) Output(arguments ...string) ([]byte, error) {
	return c.BuildCmd(arguments...).Output()
}

func (c *Container) CombinedOutput(arguments ...string) ([]byte, error) {
	return c.BuildCmd(arguments...).CombinedOutput()
}

func (c *Container) Stop() error {
	return exec.Command("docker", "stop", c.containerName).Run()
}

func (c *Container) BuildCmd(arguments ...string) *exec.Cmd {
	if c.containerName == "" {
		c.containerName = uuid.New()
	}

	args := []string{"run", "--rm"}
	if c.privileged {
		args = append(args, "--privileged")
	}
	if c.interactive {
		args = append(args, "-i")
	}
	if c.network != "" {
		args = append(args, fmt.Sprintf("--net=%s", c.network))
	}
	for key, val := range c.env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, val))
	}
	for key, val := range c.volumes {
		if IsDockerToolbox() {
			key = GetToolboxMountPath(key)
		}
		args = append(args, "-v", fmt.Sprintf("%s:%s", key, val))
	}

	if c.entrypoint != "" {
		args = append(args, "--entrypoint", c.entrypoint)
	}

	args = append(args, fmt.Sprintf("--name=%s", c.containerName))

	containerVer, ok := containerVersions[c.name]
	if !ok {
		logrus.Warnf("version for container %s not found, using version 'latest'", c.name)
		containerVer = "latest"
	}

	finalName := c.name + ":" + containerVer
	if !strings.Contains(finalName, "/") { /*bhojpur container*/
		finalName = "bhojpur/" + finalName
	}

	args = append(args, finalName)
	args = append(args, arguments...)

	logrus.WithField("args", args).Info("Build cmd for container ", finalName)

	cmd := exec.Command("docker", args...)

	return cmd
}

func IsDockerToolbox() bool {
	return runtime.GOOS == "windows" && os.Getenv("DOCKER_TOOLBOX_INSTALL_PATH") != ""
}

func GetToolboxMountPath(path string) string {
	path = strings.Replace(path, "\\", "/", -1)
	if len(path) >= 2 && path[1] == ':' {
		path = "/" + strings.ToLower(string(path[0])) + path[2:]
	}
	return path
}
