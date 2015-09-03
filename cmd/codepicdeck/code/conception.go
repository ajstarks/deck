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

func male(x, y, w int) {
	canvas.Ellipse(x, y, w, w/2, "fill:blue")
	canvas.Bezier(
		x-(w*8), y,
		x-(w*4), y-(w*4),
		x-(w*4), y+w,
		x-w, y, "stroke:blue;fill:none")
}

func female(x, y, w int) {
	canvas.Circle(x, y, w, "fill:pink")
}

func main() {
	msize := 5
	fsize := msize * 40
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	female(width, height-50, fsize)
	male(100, 200, msize)
	canvas.End()
}
