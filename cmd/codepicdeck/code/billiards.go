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
	rg := []svg.Offcolor{
		{1, "powderblue", 1},
		{10, "lightskyblue", 1},
		{100, "darkblue", 1},
	}
	canvas.Start(width, height)
	canvas.Def()
	canvas.RadialGradient("rg", 50, 50, 50, 30, 30, rg)
	canvas.DefEnd()
	canvas.Rect(0, 0, width, height, "fill:lightsteelblue")
	canvas.Gstyle("fill:url(#rg)")
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
