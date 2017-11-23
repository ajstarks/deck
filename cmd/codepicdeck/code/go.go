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

	blues := "stroke:blue"
	reds := "stroke:red"
	greens := "stroke:green"
	organges := "stroke:orange"
	
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	
	canvas.Gstyle("fill:none;stroke-opacity:0.5;stroke-width:35;stroke-linecap:round")
	// g
	canvas.Arc(20, 200, 30, 30, 0, false, true, 220, 200, blues)
	canvas.Arc(20, 200, 30, 30, 0, false, false, 220, 200, reds)
	canvas.Line(220, 100, 220, 300, greens)
	canvas.Arc(20, 320, 30, 30, 0, false, false, 220, 300, organges)
	// o
	canvas.Arc(280, 200, 30, 30, 0, false, true, 480, 200, reds)
	canvas.Arc(280, 200, 30, 30, 0, false, false, 480, 200, blues)
	canvas.Gend()
	canvas.End()
}
