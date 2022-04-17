package cowsay

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
	"errors"
	"fmt"

	"github.com/dj456119/go-cowsay/gocowsay"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/bhojpur/kernel/pkg/base/app/cowsay/animal"
)

var animalTempalteMap = make(map[string]AnimalTemplate)

type AnimalTemplate interface {
	Get() string
}

func cowsay(ctx *app.Context) error {
	err := ctx.ParseFlags()
	if err != nil {
		return err
	}

	if ctx.Flag().NArg() == 0 {
		return errors.New("no input")
	}
	info := ctx.Flag().Arg(0)
	animal := ctx.Flag().Arg(1)
	animalTemplate, err := GetAnimal(animal)
	if err != nil {
		return err
	}
	fmt.Print(gocowsay.Format(animalTemplate.Get(), info))
	return nil
}

func GetAnimal(animalType string) (AnimalTemplate, error) {
	if animalTemplate, ok := animalTempalteMap[animalType]; !ok {
		return nil, errors.New("no support animal " + animalType)
	} else {
		return animalTemplate, nil
	}
}

func RegisterAnimalTemplate(animalType string, animalTemplate AnimalTemplate) {
	animalTempalteMap[animalType] = animalTemplate
}

func init() {
	RegisterAnimalTemplate("", animal.CowTemplate{})
	RegisterAnimalTemplate("cow", animal.CowTemplate{})
	RegisterAnimalTemplate("sheep", animal.SheepTemplate{})
	RegisterAnimalTemplate("demon", animal.DemonTemplate{})
	RegisterAnimalTemplate("pig", animal.PigTemplate{})
	RegisterAnimalTemplate("monkey", animal.MonkeyTemplate{})
	app.Register("cowsay", cowsay)
}
