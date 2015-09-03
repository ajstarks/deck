package main

import (
	"fmt"
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
	opacity := 1.0
	for x := 0; x < width; x += 100 {
		canvas.Image(x, 0, 128, 128, "gophercolor128x128.png", fmt.Sprintf("opacity:%.2f", opacity))
		opacity -= 0.2
	}
	canvas.End()
}
