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
	theta  = 0.0
)

func rcircles(depth int) {
	var fill string
	switch depth {
	case 1:
		fill = canvas.RGB(89, 9, 21)
	case 2:
		fill = canvas.RGB(148, 14, 25)
	case 3:
		fill = canvas.RGB(181, 86, 70)
	case 4:
		fill = canvas.RGB(199, 172, 115)
	default:
		return
	}
	canvas.Circle(0, 0, 2, fill)
	for i := 0; i < 3; i++ {
		deg := theta + (120 * (float64(i)))
		canvas.TranslateRotate(0, 1, deg)
		canvas.Scale(0.4)
		rcircles(depth + 1)
		fmt.Fprintf(os.Stderr, "theta = %.2f depth = %d, deg = %.2f\n", theta, depth, deg)
		canvas.Gend()
		canvas.Gend()
	}
}

func main() {
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(119, 112, 127))
	for i := 0; i < 100; i++ {
		canvas.Translate(250, 250)
		canvas.Scale(120)
		rcircles(1)
		canvas.Gend()
		canvas.Gend()
		theta += 1.0
	}
	canvas.End()
}
