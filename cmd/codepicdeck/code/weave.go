package main

import (
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	pw     = 500
	ph     = 500
)

func main() {
	width, height := pw, ph
	canvas.Start(pw, ph)
	canvas.Rect(0, 0, pw, ph, "fill:gray")
	for i := 0; i < height; i += 20 {
		canvas.Rect(0, i, width, 10, "fill:black")
		canvas.Rect(i, 0, 10, height, "fill:white")
	}
	canvas.End()
}
