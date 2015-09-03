package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func meter(x, y, w, h, value int, label string) {
	corner := h / 2
	inset := corner / 2
	canvas.Text(x-10, y+h/2, label, "text-anchor:end;baseline-shift:-33%")
	canvas.Roundrect(x, y, w, h, corner, corner, "fill:rgb(240,240,240)")
	canvas.Roundrect(x+corner, y+inset, value, h-(inset*2), inset, inset, "fill:darkgray")
	canvas.Circle(x+inset+value, y+corner, inset, "fill:red;fill-opacity:0.3")
	canvas.Text(x+inset+value+inset+2, y+h/2, fmt.Sprintf("%-3d", value), "font-size:75%;text-anchor:start;baseline-shift:-33%")
}

func main() {
	rand.Seed(time.Now().Unix())
	items := []string{"Cost", "Timing", "Sourcing", "Technology"}
	mh, gutter := 50, 20
	x, y := 100, 50
	canvas.Start(width, height)
	canvas.Gstyle("font-family:sans-serif;font-size:12pt")
	for _, data := range items {
		meter(x, y, width-100, mh, rand.Intn(300), data)
		y += mh + gutter
	}
	canvas.Gend()
	canvas.End()
}
