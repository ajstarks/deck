package main

import (
	"math/rand"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func main() {
	canvas.Start(width, height)
	canvas.Gstyle("fill:white")
	var color string
	radius := 80
	step := 8
	for i := 0; i < 200; i++ {
		if i%4 == 0 {
			color = "rgb(127,0,0)"
		} else {
			color = "rgb(0,127,0)"
		}
		x, y := rand.Intn(width), rand.Intn(height)
		for r, nc := radius, 0; nc < 10; nc++ {
			canvas.Circle(x, y, r, "stroke:"+color)
			r -= step
		}
	}
	canvas.Gend()
	canvas.End()
}
