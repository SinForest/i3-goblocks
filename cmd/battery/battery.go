package main

import (
	"fmt"

	"github.com/SinForest/i3-goblocks/module"
)

const batPath = "/sys/class/power_supply/BAT0/"
const (
	fStatus = "status"
	fFull   = "energy_full"
	fNow    = "energy_now"
	sFull   = "Full"
	sDis    = "Discharging"
	sChr    = "Charging"     //TODO: check
	sNot    = "Not charging" //TODO: what is this for?

)

func colorFromPercent(perc float64) string {
	if perc > 60 {
		return "#0af548" // green
	}
	if perc > 30 {
		return "#cef50a" // yellow
	}
	if perc > 15 {
		return "#f58f0a" // orange
	}
	return "#f50a0a" // red
}

func symbolFromStatus(perc float64, status string) string {
	if status == sFull {
		return ""
	}
	if status == sNot {
		return "󰔟"
	}
	res := ""
	switch {
	case perc > 60:
		res = ""
	case perc > 30:
		res = ""
	case perc > 15:
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
	m := module.New("battery", batPath)
	chStatus := m.ReadSysFile(fStatus)
	chNow := m.ReadFloat(fNow)
	chFull := m.ReadFloat(fFull)
	chPerc := (chNow * 100.0) / chFull

	fmt.Printf("<span color='%s'>%s %3.1f%%</span>\n", colorFromPercent(chPerc), symbolFromStatus(chPerc, chStatus), chPerc)
}
