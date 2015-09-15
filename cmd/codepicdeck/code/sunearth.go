package main

import (
	"math/rand"
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
	canvas.Rect(0, 0, width, height, "fill:black")
	for i := 0; i < width; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		canvas.Line(x, y, x, y+1, "stroke:white")
	}
	earth := 4
	sun := earth * 109
	canvas.Circle(150, 50, earth, "fill:blue")
	canvas.Circle(width, height, sun, "fill:rgb(255, 248, 231)")
	canvas.End()
}
