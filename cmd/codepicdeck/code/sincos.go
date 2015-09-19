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
func main() {
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:silver")
	pichart(50, 50, 300, 100, math.Sin)
	pichart(50, 200, 200, 200, math.Cos)
	canvas.End()
}
func pichart(left, top, w, h int, function func(float64) float64) {
	const TwoPi = 2 * math.Pi
	midx, midy := left+(w/2), (h / 2)
	labx, laby := left-10, h+15
	canvas.Translate(0, top)
	canvas.Rect(left, 0, w, h, "fill:white;stroke:black")
	canvas.Gstyle("font-family:serif;font-size:12pt")
	canvas.Text(labx, 0, "1", "text-anchor:end")
	canvas.Text(labx, h, "-1", "text-anchor:end")
	canvas.Text(labx, midy, "0", "text-anchor:end")
	canvas.Text(left, laby, "0", "text-anchor:middle")
	canvas.Text(left+w, laby, "2\u03c0", "text-anchor:middle")
	canvas.Text(midx, laby, "\u03c0", "text-anchor:middle")

	canvas.Line(midx, h, midx, 0, "stroke:gray")
	canvas.Line(left, midy, left+w, midy, "stroke:gray")
	canvas.Text(left, laby+15, "sin(x)", "fill:red")
	canvas.Text(left, laby+30, "cos(x)", "fill:blue")
	canvas.Gend()

	for x := 0.0; x < TwoPi; x += math.Pi / 25 {
		dx := int(vmap(x, 0, TwoPi, float64(left), float64(w+left)))
		dsy := int(vmap(function(x), -1, 1, 0, float64(h)))
		dcy := int(vmap(math.Cos(x), -1, 1, 0, float64(h)))

		canvas.Translate(0, (h - height))
		canvas.Circle(dx, height-dsy, 3, "fill:red")
		canvas.Circle(dx, height-dcy, 3, "fill:blue")
		canvas.Gend()
	}
	canvas.Gend()
}
