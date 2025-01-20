package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SinForest/i3-goblocks/module"
)

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	flag.Parse()

	m := module.New("wifi", "/proc/net", *tick)

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
				}
				return fmt.Errorf("iwgetid returned with status %d: %v", eErr.ProcessState.ExitCode(), eErr.Stderr)
			}
			return fmt.Errorf("unknown error when running iwgetid: %v", err)
		}

		sysOut := m.ReadSysFile("wireless")
		for _, line := range strings.Split(sysOut, "\n") {
			//TODO: do this!
			line = line
		}

		return nil
	})
}
