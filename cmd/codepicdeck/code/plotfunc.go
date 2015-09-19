package main

import (
	"math"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func vmap(value float64, l1 float64, h1 float64,
	l2 float64, h2 float64) float64 {
	return l2 + (h2-l2)*(value-l1)/(h1-l1)
}

func plotfunc(left, top, w, h int, min, max, fmin, fmax,
	interval float64, f func(float64) float64, style ...string) {
	canvas.Translate(0, top)
	canvas.Rect(left, 0, w, h, "fill:white;stroke:gray")
	for x := min; x < max; x += interval {
		dx := int(vmap(x, min, max, float64(left), float64(w+left)))
		dy := int(vmap(f(x), fmin, fmax, 0, float64(h)))
		canvas.Translate(0, (h - height))
		canvas.Circle(dx, height-dy, 2, style...)
		canvas.Gend()
	}
	canvas.Gend()
}

func main() {
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:rgb(240,240,240)")
	plotfunc(80, 20, 360, 120, 0, 6*math.Pi, -1, 1, math.Pi/20, math.Sin)
	plotfunc(80, 180, 360, 120, 0, 12*math.Pi, -1, 1, math.Pi/20, math.Cos, "fill:red")
	plotfunc(80, 340, 360, 120, -3, 3, -2, 20, 0.2, math.Exp, "fill:green")
	canvas.End()
}
