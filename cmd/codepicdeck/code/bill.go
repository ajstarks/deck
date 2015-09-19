package main

import (
	"github.com/ajstarks/svgo"
	"os"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func main() {
	nr := 10
	radius := width / 10
	x := width / 2
	y := radius

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:lightsteelblue")
	canvas.Gstyle("fill:white;stroke:gray")
	for r := 0; r < nr; r++ {
		xc := x
		for c := 0; c < r+1; c++ {
			canvas.Ellipse(xc, y, radius, radius)
			xc += radius * 2
		}
		x -= radius
		y += radius
	}
	canvas.Gend()
	canvas.End()
}
