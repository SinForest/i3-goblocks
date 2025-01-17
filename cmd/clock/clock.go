package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/SinForest/i3-goblocks/module"
)

func colorFromDay(d time.Weekday) string {
	switch d {
	case time.Monday:
		return "cyan"
	case time.Tuesday:
		return "yellow"
	case time.Wednesday:
		return "DeepPink"
	case time.Thursday:
		return "lime"
	case time.Friday:
		return "Orange"
	case time.Saturday:
		return "blue"
	case time.Sunday:
		return "red"
	}
	return "grey"
}

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	flag.Parse()

	m := module.New("", "", *tick)

	m.Run(func() {
		t := time.Now()
		var out string
		out += fmt.Sprintf("<span face='monospace' color='%s'>%s</span> ",
			colorFromDay(t.Weekday()),
			t.Format("Mon"),
		)
		out += fmt.Sprintf("<span face='monospace' color='white'>%04d-%02d-</span>",
			t.Year(),
			t.Month(),
		)
		out += fmt.Sprintf("<b><span face='monospace' color='white'>%02d</span></b> ",
			t.Day(),
		)
		out += fmt.Sprintf("<span face='monospace' color='white'>%02d:%02d</span>",
			t.Hour(),
			t.Minute(),
		)
		out += fmt.Sprintf("<span face='monospace' color='grey'>:%02d</span>",
			t.Second(),
		)
		fmt.Println(out)
	})
}
