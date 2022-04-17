//go:build nes
// +build nes

package cmd

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
	"flag"
	"fmt"
	"image"
	"io"
	"time"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/bhojpur/kernel/pkg/base/drivers/kbd"
	"github.com/bhojpur/kernel/pkg/base/drivers/vbe"
	"github.com/bhojpur/kernel/pkg/base/log"

	"github.com/fogleman/nes/nes"
	"golang.org/x/image/draw"
)

var buttonmap = map[byte]int{
	'a':  nes.ButtonLeft,
	'd':  nes.ButtonRight,
	'w':  nes.ButtonUp,
	's':  nes.ButtonDown,
	'k':  nes.ButtonA,
	'j':  nes.ButtonB,
	' ':  nes.ButtonSelect,
	'\n': nes.ButtonStart,
}

var scaleAlgs = map[string]draw.Interpolator{
	"NearestNeighbor": draw.NearestNeighbor,
	"ApproxBiLinear":  draw.ApproxBiLinear,
}

func nesmain(ctx *app.Context) error {
	if !vbe.IsEnable() {
		return errors.New("video not enabled")
	}
	var (
		flagSet  = flag.NewFlagSet(ctx.Args[0], flag.ContinueOnError)
		gameName = flagSet.String("rom", "", "game rom")
		scaleAlg = flagSet.String("scale", "NearestNeighbor", "NearestNeighbor|ApproxBiLinear")
	)

	flagSet.SetOutput(ctx.Stdout)
	err := flagSet.Parse(ctx.Args[1:])
	if err != nil {
		return err
	}
	if *gameName == "" {
		return errors.New("-rom required")
	}

	rom, err := ctx.Open(*gameName)
	if err != nil {
		return err
	}
	defer rom.Close()

	scaler, ok := scaleAlgs[*scaleAlg]
	if !ok {
		return fmt.Errorf("scale alg not found:%s", *scaleAlg)
	}
	return runGame(rom, scaler)
}

func buttons() [8]bool {
	var ret [8]bool
	for c, b := range buttonmap {
		if kbd.Pressed(c) {
			ret[b] = true
		}
	}
	return ret
}

func runGame(rom io.Reader, scaler draw.Interpolator) error {
	con, err := nes.NewConsoleReader(rom)
	if err != nil {
		return err
	}
	con.StepFrame()

	view := vbe.NewView()
	old := vbe.SaveCurrView()
	defer vbe.SetCurrView(old)

	vbe.SetCurrView(view)

	screen := view.Canvas()

	sum := time.Duration(0)
	cnt := 0

	scale := float32(1.5)
	width := int(float32(256) * scale)
	height := int(float32(240) * scale)
	point := image.Pt((screen.Bounds().Dx()-width)/2, (screen.Bounds().Dy()-height)/2)
	rect := image.Rect(0, 0, width, height).Add(point)

	for {
		if kbd.Pressed('q') {
			return nil
		}
		begin := time.Now()
		con.SetButtons1(buttons())
		con.StepFrame()
		buf := con.Buffer()
		draw.Draw(screen, screen.Bounds(), buf, image.ZP, draw.Src)
		scaler.Scale(screen, rect, buf, buf.Bounds(), draw.Src, nil)
		view.CommitRect(rect)
		used := time.Now().Sub(begin)

		sum += used
		cnt++
		if cnt%30 == 0 {
			log.Infof("used %s", sum/time.Duration(cnt))
			cnt = 0
			sum = 0
		}
	}
}

func init() {
	app.Register("nes", nesmain)
}
