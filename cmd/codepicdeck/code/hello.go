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
	style := "fill:white;font-size:48pt;text-anchor:middle"
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	canvas.Circle(width/2, height, width/2, "fill:rgb(44, 77, 232)")
	canvas.Text(width/2, height/3, "hello, world", style)
	canvas.End()
}
