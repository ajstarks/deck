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
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	r := height / 2
	for g := 0; g < 250; g += 50 {
		canvas.Circle(width/2, width/2, r, canvas.RGB(g, g, g))
		r -= 50
	}
	canvas.End()
}
