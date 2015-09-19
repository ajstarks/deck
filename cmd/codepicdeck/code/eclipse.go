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
	h2 := height / 2
	r := width / 20
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	for x, y := 50, h2; x < 450; x += 75 {
		canvas.Circle(x, h2, r+1, "fill:white")
		canvas.Circle(x, y, r, "fill:black")
		y += 10
	}
	canvas.End()
}
