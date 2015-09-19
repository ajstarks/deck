package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

type Measure struct {
	name  string
	value int
}

func (data *Measure) meter(x, y, w, h int) {
	corner := h / 2
	inset := corner / 2
	canvas.Text(x-10, y+h/2, data.name, "text-anchor:end;baseline-shift:-33%")
	canvas.Roundrect(x, y, w, h, corner, corner, "fill:rgb(240,240,240)")
	canvas.Roundrect(x+corner, y+inset, data.value, h-(inset*2), inset, inset, "fill:darkgray")
	canvas.Circle(x+inset+data.value, y+corner, inset, "fill:red;fill-opacity:0.3")
	canvas.Text(x+inset+data.value+inset+2, y+h/2, fmt.Sprintf("%-3d", data.value),
		"font-size:75%;text-anchor:start;baseline-shift:-33%")
}

func main() {
	items := []Measure{{"Cost", 100}, {"Timing", 250}, {"Sourcing", 50}, {"Technology", 175}}
	x, y, gutter, mh := 100, 50, 20, 50
	canvas.Start(width, height)
	canvas.Gstyle("font-family:sans-serif;font-size:12pt")
	for _, data := range items {
		data.meter(x, y, width-100, mh)
		y += mh + gutter
	}
	canvas.Gend()
	canvas.End()
}
