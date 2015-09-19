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

func tro() {
	canvas.Rect(100, 100, 100, 100)
	canvas.TranslateRotate(100, 100, 30)
	canvas.Grid(0, 0, width, height, 50, "stroke:pink")
	canvas.Rect(0, 0, 100, 100, "fill:red")
	canvas.Gend()
}

func main() {
	canvas.Start(width, height)
	canvas.Grid(0, 0, width, height, 50, "stroke:lightsteelblue")
	tro()
	canvas.End()
}
