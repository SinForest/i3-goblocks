package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	dev := flag.String("dev", "", "wifi device name , e.g. 'wlp0s20f3' or 'wifi0'; if not given, choose first device")
	flag.Parse()

	m := module.New("wifi", "/proc/net", *tick)
	cm := colormap.DefaultMap()

	m.Run(func() error {
		cmd := exec.Command("iwgetid", "-r")
		sb := new(strings.Builder)
		cmd.Stdout = sb
		err := cmd.Run()
		if err != nil {
			var eErr *exec.ExitError
			if errors.As(err, &eErr) {
				if eErr.ProcessState.ExitCode() == 255 {
					fmt.Printf("<span color='#ff0000'> </span>: no wifi\n")
					return nil
				}
				return fmt.Errorf("iwgetid returned with status %d: %v", eErr.ProcessState.ExitCode(), eErr.Stderr)
			}
			return fmt.Errorf("unknown error when running iwgetid: %v", err)
		}

		sysOut := m.MustReadSysFile("wireless")
		perc := 0.0
		for _, line := range strings.Split(sysOut, "\n")[2:] /*skip header*/ {
			if *dev != "" && !strings.HasPrefix(line, *dev+":") {
				continue
			}
			f, err := strconv.ParseFloat(strings.Fields(line)[2], 64)
			if err != nil {
				return fmt.Errorf("could not parse float from line %q: %w", line, err)
			}
			perc = f / 70
		}

		//TODO: ip?

		fmt.Printf("<b><span color='%[1]s'> :</span></b> <span>%[2]s</span> <b><span color='%[1]s' face='monospace'>(%3.[3]f%%)</span></b>\n", cm.Eval(perc), strings.TrimSuffix(sb.String(), "\n"), perc*100)

		return nil
	})
}
