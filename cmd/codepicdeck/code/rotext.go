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
	a := 1.0
	ai := 0.03
	ti := 10.0

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	canvas.Gstyle("font-family:serif;font-size:244pt")
	for t := 0.0; t <= 360.0; t += ti {
		canvas.TranslateRotate(width/2, height/2, t)
		canvas.Text(0, 0, "s", canvas.RGBA(255, 255, 255, a))
		canvas.Gend()
		a -= ai
	}
	canvas.Gend()
	canvas.End()
}
