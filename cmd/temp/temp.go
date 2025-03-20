package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

func scale(min, max, current float64) float64 {
	return (current - min) / (max - min)
}

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")

	flag.Parse()

	m := module.New("cputemp", "/sys/class/thermal/", *tick)
	cm := colormap.New(255, 255, 255, 255, 176, 222)
	cm.Register(0.1, 176, 239, 255) // ice blue
	cm.Register(0.25, 39, 135, 219) // cold blue
	cm.Register(0.4, 141, 128, 171) // chill purple
	cm.Register(0.6, 199, 191, 149) // insignificant yellow
	cm.Register(0.75, 232, 83, 14)  // fire orange
	cm.Register(0.9, 255, 0, 76)    // crimson red

	m.Run(func() error {
		idx := 0
		max := 0.0
		for {
			temp, err := m.ReadFloat(fmt.Sprintf("thermal_zone%d/temp", idx))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					break
				}
				return err
			}
			if temp > max {
				max = temp
			}
			idx++
		}
		max /= 1000

		color := cm.Eval(scale(15, 90, max))
		fmt.Printf("<span color='%[1]s'> </span>: <span face='monospace' color='%[1]s'>%4.1[2]f°C</span>\n", color, max)

		return nil
	})
}
