package main

import (
	"flag"
	"fmt"
	"iter"
	"strings"

	"github.com/SinForest/i3-goblocks/bytesize"
	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

func ramFromLines(seq iter.Seq[string]) (free bytesize.ByteSize, perc float64, err error) {
	var total bytesize.ByteSize
	var foundTotal, foundFree bool
	for s := range seq {
		key, val, ok := strings.Cut(s, ":")
		if !ok {
			continue
		}
		switch key {
		case "MemTotal":
			total, err = bytesize.Parse(val)
			if err != nil {
				return 0, 0, err
			}
			foundTotal = true
		case "MemAvailable":
			free, err = bytesize.Parse(val)
			if err != nil {
				return 0, 0, err
			}
			foundFree = true
		}
	}
	if !(foundTotal && foundFree) {
		return 0, 0, fmt.Errorf("MemTotal or MemAvailable not found")
	}
	return free, free.Over(total), nil
}

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	thresh := flag.Float64("thresh", 0.33, "the percentage [0,1] at which the output is colored green")
	flag.Parse()

	m := module.New("ram", "/proc/", *tick)
	cm := colormap.DefaultMap()
	cm.AddTopThreshold(*thresh)

	m.Run(func() error {
		file := m.ScanSysFile("meminfo")
		free, perc, err := ramFromLines(file)
		if err != nil {
			return err
		}
		color := cm.Eval(perc)

		fmt.Printf("<span color='%[1]s'> :</span> %[2]s (<span face='monospace' color='%[1]s'>%4.1[3]f%%</span>)\n", color, free, perc*100)

		return nil
	})
}
