package state

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
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/bhojpur/kernel/pkg/util/errors"
)

type basicState struct {
	imagesLock    sync.RWMutex
	instancesLock sync.RWMutex
	volumesLock   sync.RWMutex
	saveLock      sync.Mutex
	saveFile      string
	Images        map[string]*types.Image    `json:"Images"`
	Instances     map[string]*types.Instance `json:"Instances"`
	Volumes       map[string]*types.Volume   `json:"Volumes"`
}

func NewBasicState(saveFile string) *basicState {
	return &basicState{
		saveFile:  saveFile,
		Images:    make(map[string]*types.Image),
		Instances: make(map[string]*types.Instance),
		Volumes:   make(map[string]*types.Volume),
	}
}

func BasicStateFromFile(saveFile string) (*basicState, error) {
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return nil, errors.New("error reading save file "+saveFile, err)
	}
	var s basicState
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, errors.New("failed to unmarshal data "+string(data)+" to memory state", err)
	}
	if s.Images == nil {
		s.Images = make(map[string]*types.Image)
	}
	if s.Instances == nil {
		s.Instances = make(map[string]*types.Instance)
	}
	if s.Volumes == nil {
		s.Volumes = make(map[string]*types.Volume)
	}
	s.saveFile = saveFile
	return &s, nil
}

func (s *basicState) GetImages() map[string]*types.Image {
	s.imagesLock.RLock()
	defer s.imagesLock.RUnlock()
	imagesCopy := make(map[string]*types.Image)
	for id, image := range s.Images {
		imageCopy := *image
		imagesCopy[id] = &imageCopy
	}
	return imagesCopy
}

func (s *basicState) GetInstances() map[string]*types.Instance {
	s.instancesLock.RLock()
	defer s.instancesLock.RUnlock()
	instancesCopy := make(map[string]*types.Instance)
	for id, instance := range s.Instances {
		instanceCopy := *instance
		instancesCopy[id] = &instanceCopy
	}
	return instancesCopy
}

func (s *basicState) GetVolumes() map[string]*types.Volume {
	s.volumesLock.RLock()
	defer s.volumesLock.RUnlock()
	volumesCopy := make(map[string]*types.Volume)
	for id, volume := range s.Volumes {
		volumeCopy := *volume
		volumesCopy[id] = &volumeCopy
	}
	return volumesCopy
}

func (s *basicState) ModifyImages(modify func(images map[string]*types.Image) error) error {
	s.imagesLock.Lock()
	defer s.imagesLock.Unlock()
	if err := modify(s.Images); err != nil {
		return errors.New("modifying Images", err)
	}
	return s.save()
}

func (s *basicState) ModifyInstances(modify func(instances map[string]*types.Instance) error) error {
	s.instancesLock.Lock()
	defer s.instancesLock.Unlock()
	if err := modify(s.Instances); err != nil {
		return errors.New("modifying Instances", err)
	}
	return s.save()
}

func (s *basicState) ModifyVolumes(modify func(volumes map[string]*types.Volume) error) error {
	s.volumesLock.Lock()
	defer s.volumesLock.Unlock()
	if err := modify(s.Volumes); err != nil {
		return errors.New("modifying Volumes", err)
	}
	return s.save()
}

func (s *basicState) save() error {
	s.saveLock.Lock()
	defer s.saveLock.Unlock()
	data, err := json.Marshal(s)
	if err != nil {
		return errors.New("failed to marshal memory state to json", err)
	}
	os.MkdirAll(filepath.Dir(s.saveFile), 0755)
	err = ioutil.WriteFile(s.saveFile, data, 0644)
	if err != nil {
		return errors.New("writing save file "+s.saveFile, err)
	}
	return nil
}

func (s *basicState) RemoveImage(image *types.Image) error {
	if err := s.ModifyImages(func(images map[string]*types.Image) error {
		delete(images, image.Id)
		return nil
	}); err != nil {
		return errors.New("modifying image map in state", err)
	}
	return nil
}

func (s *basicState) RemoveInstance(instance *types.Instance) error {
	if err := s.ModifyInstances(func(instances map[string]*types.Instance) error {
		delete(instances, instance.Id)
		return nil
	}); err != nil {
		return errors.New("modifying image map in state", err)
	}
	volumesToDetach := []*types.Volume{}
	volumes := s.GetVolumes()
	for _, volume := range volumes {
		if volume.Attachment == instance.Id {
			volumesToDetach = append(volumesToDetach, volume)
		}
	}
	for _, volume := range volumesToDetach {
		if err := s.ModifyVolumes(func(volumes map[string]*types.Volume) error {
			volume, ok := volumes[volume.Id]
			if !ok {
				return errors.New("no record of "+volume.Id+" in the state", nil)
			}
			volume.Attachment = ""
			return nil
		}); err != nil {
			return errors.New("modifying volume map in state", err)
		}
	}
	return nil
}

func (s *basicState) RemoveVolume(volume *types.Volume) error {
	if err := s.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		delete(volumes, volume.Id)
		return nil
	}); err != nil {
		return errors.New("modifying volume map in state", err)
	}
	return nil
}
