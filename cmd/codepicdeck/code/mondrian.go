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

// Piet Mondrian - Composition in Red, Blue, and Yellow
func main() {
	w3 := width / 3
	w6 := w3 / 2
	w23 := w3 * 2
	canvas.Start(width, height)
	canvas.Gstyle("stroke:black;stroke-width:6")
	canvas.Rect(0, 0, w3, w3, "fill:white")
	canvas.Rect(0, w3, w3, w3, "fill:white")
	canvas.Rect(0, w23, w3, w3, "fill:blue")
	canvas.Rect(w3, 0, w23, w23, "fill:red")
	canvas.Rect(w3, w23, w23, w3, "fill:white")
	canvas.Rect(width-w6, height-w3, w3-w6, w6, "fill:white")
	canvas.Rect(width-w6, height-w6, w3-w6, w6, "fill:yellow")
	canvas.Gend()
	canvas.Rect(0, 0, width, height, "fill:none;stroke:black;stroke-width:12")
	canvas.End()
}
