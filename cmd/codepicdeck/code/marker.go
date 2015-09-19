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
	canvas.Def()
	canvas.Marker("dot", 10, 10, 16, 16)
	canvas.Circle(10, 10, 6, "fill:black")
	canvas.MarkerEnd()

	canvas.Marker("box", 10, 10, 16, 16)
	canvas.CenterRect(10, 10, 12, 12, "fill:green")
	canvas.MarkerEnd()

	canvas.Marker("arrow", 4, 12, 26, 26)
	canvas.Path("M4,4 L4,22 L20,12 L4,4", "fill:blue")
	canvas.MarkerEnd()
	canvas.DefEnd()

	x := []int{100, 250, 100, 150}
	y := []int{100, 250, 400, 250}
	canvas.Polyline(x, y,
		`fill="none"`,
		`stroke="red"`,
		`marker-start="url(#dot)"`,
		`marker-mid="url(#arrow)"`,
		`marker-end="url(#box)"`)
	canvas.End()
}
