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
		canvas.Image(x, 100, 100, 124, "gopher.jpg", fmt.Sprintf("opacity:%.2f", opacity))
		opacity -= 0.15
	}
	canvas.End()
}
