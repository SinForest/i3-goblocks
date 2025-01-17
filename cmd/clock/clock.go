package main

import (
	"fmt"
	"time"
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
}
