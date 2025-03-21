package colormap

//TODO: should this be its own package?

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/constraints"
)

var (
	Redish    = Color{240, 0, 90}
	Red       = Color{255, 0, 0}
	Orange    = Color{255, 150, 0}
	Yellow    = Color{255, 210, 0}
	Green     = Color{55, 219, 37}
	Turquoise = Color{0, 255, 170}

	White = Color{255, 255, 255}
)

func DefaultMap() *ColorMap {
	return &ColorMap{
		colorMapEntry{
			pos:   0.0,
			color: Redish,
		},
		colorMapEntry{
			pos:   0.1,
			color: Red,
		},
		colorMapEntry{
			pos:   0.3,
			color: Orange,
		},
		colorMapEntry{
			pos:   0.5,
			color: Yellow,
		},
		colorMapEntry{
			pos:   0.9,
			color: Green,
		},
		colorMapEntry{
			pos:   1.0,
			color: Turquoise,
		},
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

func clamp[T Number](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func c01(v float64) float64 {
	return clamp(v, 0, 1)
}

func cRGB[T Number](v T) uint8 {
	return uint8(clamp(int(v), 0, 255))
}

type colorMapEntry struct {
	pos   float64
	color Color
}

type Color [3]uint8

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", cRGB(c[0]), cRGB(c[1]), cRGB(c[2]))
}

type ColorMap []colorMapEntry

func New(r0, g0, b0 int, r1, g1, b1 int) *ColorMap {
	return &ColorMap{
		colorMapEntry{
			pos:   0.0,
			color: Color{cRGB(r0), cRGB(g0), cRGB(b0)},
		},
		colorMapEntry{
			pos:   1.0,
			color: Color{cRGB(r1), cRGB(g1), cRGB(b1)},
		},
	}
}

func (cm *ColorMap) Register(pos float64, r, g, b int) {
	pos = clamp(pos, 0, 1)
	R := cRGB(r)
	G := cRGB(g)
	B := cRGB(b)

	var idx int
	for ; idx < len(*cm); idx++ {
		if (*cm)[idx].pos == pos { // diff < eps, instead of ==?
			//TODO: this should be an error, or we divide by 0 later :(
		}
		if (*cm)[idx].pos > pos {
			break // idx found
		}
	}
	*cm = slices.Insert(*cm, idx, colorMapEntry{pos, Color{R, G, B}})
}

func (cm *ColorMap) Eval(pos float64) Color {
	pos = c01(pos)
	for idx := 0; idx < len(*cm); idx++ {
		if (*cm)[idx].pos == pos {
			return (*cm)[idx].color
		}
		if (*cm)[idx].pos > pos {
			bot := (*cm)[idx-1]
			top := (*cm)[idx]
			rTop := (pos - bot.pos) / (top.pos - bot.pos)
			r := float64(top.color[0])*rTop + float64(bot.color[0])*(1-rTop)
			g := float64(top.color[1])*rTop + float64(bot.color[1])*(1-rTop)
			b := float64(top.color[2])*rTop + float64(bot.color[2])*(1-rTop)
			return Color{cRGB(r), cRGB(g), cRGB(b)}
		}
	}
	return (*cm)[len(*cm)-1].color // this should never be reached
}

// AddTopThreshold stretches the first color transition, such that the second
// color is evaluated at position pos2nd, and equally compresses the remaining transition.
func (cm *ColorMap) AddTopThreshold(pos2nd float64) {
	if pos2nd <= 0 || pos2nd > 1 || len(*cm) < 2 {
		return
	}
	factor := pos2nd / (*cm)[len(*cm)-2].pos
	for idx, col := range (*cm)[:len(*cm)-1] {
		col.pos *= factor
		(*cm)[idx] = col
	}
}

func (cm *ColorMap) String() string {
	sb := strings.Builder{}
	for idx, col := range *cm {
		sb.WriteString(fmt.Sprintf("%6.2f%%: %s;", col.pos*100, col.color))
		if idx != len(*cm)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
