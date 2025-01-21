package main

import (
	"flag"
	"fmt"

	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

const batPath = "/sys/class/power_supply/BAT0/"
const (
	fStatus = "status"
	fFull   = "_full"
	fNow    = "_now"
	sFull   = "Full"
	sDis    = "Discharging"
	sChr    = "Charging"
	sIdle   = "Not charging"
)

func symbolFromStatus(perc float64, status string) string {
	if status == sFull {
		return ""
	}
	if status == sIdle {
		return "󰔟"
	}
	res := ""
	switch {
	case perc > 0.6:
		res = ""
	case perc > 0.3:
		res = ""
	case perc > 0.15:
		res = ""
	default:
		res = ""
	}
	if status == sChr {
		return "" + res
	}
	return res
}

func main() {
	var warnBelow float64
	var filePrefix string
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	flag.Float64Var(&warnBelow, "warn-below", 0.3, "percentage, under which the whole output is colorized")
	flag.StringVar(&filePrefix, "prefix", "energy", "file prefix of energy files, usually 'charge' or 'energy'")
	flag.Parse()

	m := module.New("battery", batPath, *tick)
	cm := colormap.DefaultMap()

	m.Run(func() error {
		chStatus := m.ReadSysFile(fStatus)
		chNow := m.ReadFloat(filePrefix + fNow)
		chFull := m.ReadFloat(filePrefix + fFull)
		chPerc := chNow / chFull
		color := cm.Eval(chPerc)
		textColor := color
		if chPerc > warnBelow {
			textColor = colormap.White
		}

		fmt.Printf("<span color='%s'>%s :</span> <span face='monospace' color='%s'>%3.1f%%</span>\n", color, symbolFromStatus(chPerc, chStatus), textColor, chPerc*100)

		return nil
	})
}
