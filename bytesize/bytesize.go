package bytesize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var factorSizes = map[string]int64{
	"B":  1,
	"kB": 1024,
	"MB": 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
	"TB": 1024 * 1024 * 1024 * 1024,
	"PB": 1024 * 1024 * 1024 * 1024 * 1024,
}
var sizeFactors = []string{"B", "kB", "MB", "GB", "TB", "PB"}

type ByteSize uint64

func Parse(s string) (ByteSize, error) {
	s = strings.TrimSpace(s)
	am, unit, ok := strings.Cut(s, " ")
	if !ok {
		return 0, fmt.Errorf("%q is not a valid byte size", s)
	}
	amount, err := strconv.ParseInt(strings.TrimSpace(am), 10, 64)
	if err != nil {
		return 0, err
	}
	unit = strings.TrimSpace(unit)
	factor, ok := factorSizes[unit]
	if !ok {
		return 0, fmt.Errorf("%q is not a valid byte size unit", unit)
	}
	return ByteSize(amount * factor), nil
}

func (bs ByteSize) String() string {
	if bs < 1024 {
		return fmt.Sprintf("%d B", bs)
	}
	running := float64(bs) / 1024
	var fac string
	for _, fac = range sizeFactors[1:] {
		if running < 1024 {
			break
		}
		running /= 1024
	}
	return fmt.Sprintf("%.2f %s", running, fac)
}

func (bs ByteSize) Over(total ByteSize) float64 {
	if total == 0 {
		return math.Inf(1)
	}
	return float64(bs) / float64(total)
}
