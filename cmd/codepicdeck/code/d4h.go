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

func dot(x, y, d int) {
	canvas.Circle(x, y, d/2, "fill:rgb(128,0,128)")
}

// Composition from "Design for Hackers, pg. 129
func main() {
	d1 := height
	d2 := d1 / 4
	d3 := (d2 * 3) / 4
	d4 := (d3 * 3) / 4

	coffset := height / 8
	hoffset := height / (height / 10)
	voffset := -width / 10

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:lightblue")
	dot(width-coffset, height-coffset, d1)
	dot(width/2, height/3, d2)
	dot(width/4, height*2/3, d3)
	dot(width/4+hoffset, height/3+voffset, d4)
	canvas.Grid(0, 0, width, height, width/4, "stroke:red")
	canvas.End()
}
