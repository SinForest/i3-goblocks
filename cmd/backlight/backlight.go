package main

import (
	"flag"
	"fmt"
	"math"
	"path"

	"github.com/SinForest/i3-goblocks/module"
)

const pathBL = "/sys/class/backlight/"

const (
	fMax = "max_brightness"
	fNow = "brightness"
)

func greyFromPercent(ratio float64) string {
	c := int(ratio * 255)
	return fmt.Sprintf("#%02[1]x%02[1]x%02[1]x", c)
}

func main() {
	blDir := flag.String("dir", "intel_backlight", "specify, if backlight dir is not 'intel_backlight', e.g. 'acpi_video0'")
	exponent := flag.Float64("exp", 1, "")
	flag.Parse()

	m := module.New("backlight", path.Join(pathBL, *blDir))

	blMax := m.ReadFloat(fMax)
	blNow := m.ReadFloat(fNow)
	blRatio := (blNow / blMax)
	if *exponent != 1 {
		blRatio = math.Pow(blRatio, *exponent)
	}
	fmt.Printf("<span color='%s'>%3.0f%%</span>\n", greyFromPercent(blRatio), blRatio*100)

}
