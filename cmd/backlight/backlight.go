package main

import (
	"flag"
	"fmt"
	"math"
	"path"

	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

const pathBL = "/sys/class/backlight/"

const wheelDelta = 0.02
const minWheelBrightness = 0.02

const (
	fMax = "max_brightness"
	fNow = "brightness"
)

type brightnessModule struct {
	*module.Module
	exponent float64
}

func (bm *brightnessModule) ratio() float64 {
	blMax := bm.ReadFloat(fMax)
	blNow := bm.ReadFloat(fNow)
	blRatio := (blNow / blMax)
	if bm.exponent != 1 {
		blRatio = math.Pow(blRatio, bm.exponent)
	}
	return blRatio
}

func (bm *brightnessModule) setRatio(to float64) {
	if to < minWheelBrightness {
		to = minWheelBrightness
	}
	if bm.exponent != 1 {
		to = math.Pow(to, 1/bm.exponent)
	}
	blMax := bm.ReadFloat(fMax)
	to *= blMax

	if to > blMax {
		to = blMax
	}
	if to < 0 {
		to = 0
	}

	bm.WriteSysFile(fNow, fmt.Sprintf("%d", int(to+0.5)))
}

func main() {
	blDir := flag.String("dir", "intel_backlight", "specify, if backlight dir is not 'intel_backlight', e.g. 'acpi_video0'")
	exponent := flag.Float64("exp", 1, "exponent to scale percentage, to make low percentages more differentiable")
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	flag.Parse()

	m := module.New("backlight", path.Join(pathBL, *blDir), *tick)
	bm := brightnessModule{m, *exponent}
	m.RegisterClickHandler(func(m *module.Module, click module.Click) {
		switch click {
		case module.BtnWheelUp:
			bm.setRatio(bm.ratio() + wheelDelta)
		case module.BtnWheelDown:
			bm.setRatio(bm.ratio() - wheelDelta)
		case module.BtnMiddle:
			if bm.ratio() >= 0.9 {
				bm.setRatio(0.4)
			} else {
				bm.setRatio(1)
			}
		}
	})
	cm := colormap.New(120, 120, 120, 255, 255, 255)
	cm.Register(0.25, 100, 100, 100)
	m.Run(func() error {
		blRatio := bm.ratio()
		fmt.Printf("<span color='#ffff00'>󱠂</span> : <span color='%s'>%3.0f%%</span>\n", cm.Eval(blRatio), blRatio*100)
		return nil
	})

}
