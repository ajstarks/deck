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

func main() {
	nr := 6
	radius := width / 10
	x := width / 2
	y := height / 2
	fgcolor := "white"
	bgcolor := "lightsteelblue"
	sw := width / 50
	sfmt := "fill:%s;;stroke:%s;stroke-width:%dpx"

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:"+bgcolor)
	canvas.Gstyle(fmt.Sprintf(sfmt, fgcolor, bgcolor, sw))
	for r := 0; r < nr; r++ {
		xc := x
		for c := 0; c < r+1; c++ {
			canvas.Circle(xc, y, radius)
			xc += radius * 2
		}
		x -= radius
		y += radius
	}
	canvas.Gend()
	canvas.End()
}
