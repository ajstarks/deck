package main

import (
	"github.com/ajstarks/svgo"
	"math"
	"os"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func plotfunc(left, top, w, h int, style string, min, max, fmin, fmax, interval float64, f func(float64) float64) {
	canvas.Translate(0, top)
	canvas.Rect(left, 0, w, h, "fill:white;stroke:gray")
	for x := min; x < max; x += interval {
		dx := int(vmap(x, min, max, float64(left), float64(w+left)))
		dy := int(vmap(f(x), fmin, fmax, 0, float64(h)))
		canvas.Translate(0, (h - height))
		canvas.Circle(dx, height-dy, 3, style)
		canvas.Gend()
	}
	canvas.Gend()
}

func main() {
	const TwoPi = 2 * math.Pi
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:silver")
	plotfunc(50, 50, 300, 100, "fill:red", 0, TwoPi, -1, 1, math.Pi/25, math.Sin)
	plotfunc(50, 200, 300, 100, "fill:blue", 0, math.Pi, 0, 1, math.Pi/25, math.Tan)
	canvas.End()
}
