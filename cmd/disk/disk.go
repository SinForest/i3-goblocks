package main

import (
	"flag"
	"fmt"

	"golang.org/x/sys/unix"

	"github.com/SinForest/i3-goblocks/bytesize"
	"github.com/SinForest/i3-goblocks/colormap"
	"github.com/SinForest/i3-goblocks/module"
)

func main() {
	tick := flag.Int("tick", 0, "for i3blocks persist mode: if > 0, update interval in seconds")
	mount := flag.String("mnt", "/", "path inside the disk in question")
	relative := flag.Bool("relative", false, "whether to additionally display the relative available space in percent")
	thresh := flag.Float64("thresh", 0.2, "the percentage [0,1] at which the output is colored green")
	flag.Parse()

	m := module.New("disk", "", *tick)
	cm := colormap.DefaultMap()
	cm.AddTopThreshold(*thresh)

	var stat = new(unix.Statfs_t)

	m.Run(func() error {
		unix.Statfs(*mount, stat)

		// Available blocks * size per block = available space in bytes
		avail := bytesize.ByteSize(stat.Bavail * uint64(stat.Bsize))
		total := bytesize.ByteSize(stat.Blocks * uint64(stat.Bsize))
		perc := avail.Over(total)

		color := cm.Eval(perc)

		out := fmt.Sprintf("<b><span color='%s'> :</span></b> %s", color, avail)
		if *relative {
			out += fmt.Sprintf(" (<span face='monospace' color='%s'>%4.1f%%</span>)", color, perc*100)
		}
		fmt.Println(out)

		return nil
	})
}
