package vsphere

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
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/bhojpur/kernel/pkg/config"
	"github.com/bhojpur/kernel/pkg/providers/common"
	"github.com/bhojpur/kernel/pkg/providers/vsphere/vsphereclient"
	"github.com/bhojpur/kernel/pkg/state"
	"github.com/bhojpur/kernel/pkg/util/errors"
	"github.com/sirupsen/logrus"
)

func VsphereStateFile() string {
	return filepath.Join(config.Internal.KernelHome, "vsphere/state.json")
}

var VsphereImagesDirectory = "bhojpur/vsphere/images/"
var VsphereVolumesDirectory = "bhojpur/vsphere/volumes/"

const VsphereKernelInstanceListener = "VsphereKernelInstanceListener"
const instanceListenerPrefix = "bhojpur_vsphere"

type VsphereProvider struct {
	config             config.Vsphere
	state              state.State
	u                  *url.URL
	instanceListenerIp string
}

func NewVsphereProvier(config config.Vsphere) (*VsphereProvider, error) {
	rawUrl := "https://" + config.VsphereUser + ":" + config.VspherePassword + "@" + strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(config.VsphereURL, "http://"), "https://"), "/sdk") + "/sdk"
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.New("parsing vsphere url", err)
	}

	p := &VsphereProvider{
		config: config,
		state:  state.NewBasicState(VsphereStateFile()),
		u:      u,
	}

	p.getClient().Mkdir("kernctl")
	p.getClient().Mkdir("bhojpur/vsphere")
	p.getClient().Mkdir("bhojpur/vsphere/images")
	p.getClient().Mkdir("bhojpur/vsphere/volumes")

	if err := p.deployInstanceListener(); err != nil {
		return nil, errors.New("deploying vSphere instance listener", err)
	}

	instanceListenerIp, err := common.GetInstanceListenerIp(instanceListenerPrefix, timeout)
	if err != nil {
		return nil, errors.New("failed to retrieve instance listener IP. is Bhojpur Kernel instance listener running?", err)
	}

	p.instanceListenerIp = instanceListenerIp

	tmpDir := filepath.Join(os.Getenv("HOME"), ".bhojpur", "tmp")
	os.Setenv("TMPDIR", tmpDir)
	logrus.Infof("Creating directory %s", tmpDir)
	os.MkdirAll(tmpDir, 0755)

	// begin update instances cycle
	go func() {
		for {
			if err := p.syncState(); err != nil {
				logrus.Error("error updating vSphere state:", err)
			}
			time.Sleep(time.Second)
		}
	}()

	return p, nil
}

func (p *VsphereProvider) WithState(state state.State) *VsphereProvider {
	p.state = state
	return p
}

func (p *VsphereProvider) getClient() *vsphereclient.VsphereClient {
	return vsphereclient.NewVsphereClient(p.u, p.config.Datastore, p.config.Datacenter)
}

//just for consistency
func getInstanceDatastoreDir(instanceName string) string {
	return instanceName
}

func getImageDatastoreDir(imageName string) string {
	return filepath.Join(VsphereImagesDirectory, imageName+"/")
}

func getImageDatastorePath(imageName string) string {
	return filepath.Join(getImageDatastoreDir(imageName), "boot.vmdk")
}

func getVolumeDatastoreDir(volumeName string) string {
	return filepath.Join(VsphereVolumesDirectory, volumeName+"/")
}

func getVolumeDatastorePath(volumeName string) string {
	return filepath.Join(getVolumeDatastoreDir(volumeName), "data.vmdk")
}
