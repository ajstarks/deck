package main

import (
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func main() {
	y := 20
	v := 10
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:rgb(240,240,240)")
	canvas.Gstyle("stroke:white")
	for x := 20; x < 450; x += 30 {
		op := 0.1
		for i := 0; i < 100; i += 10 {
			canvas.Square(x, y, 20*2, canvas.RGBA(v, 0, 0, op))
			y += 30
			op += 0.1
		}
		y = 20
		v += 25
	}
	canvas.Gend()
	canvas.End()
}
