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

func vmap(value float64, l1 float64, h1 float64, l2 float64, h2 float64) float64 {
	return l2 + (h2-l2)*(value-l1)/(h1-l1)
}
func main() {
	const TwoPi = 2 * math.Pi
	left, top, w, h := 50, 50, 400, 400 // 125, 100, 350, 200
	right, midx, midy, labx, laby := left+w, left+(w/2), h/2, left-10, h+18
	canvas.Start(width, height)
	canvas.Translate(0, top)
	canvas.Rect(left, 0, w, h, "stroke:gray;fill:white")
	canvas.Gstyle("font-family:serif;font-size:14pt")
	canvas.Text(left, laby, "0", "text-anchor:middle")
	canvas.Text(midx, laby, "\u03c0", "text-anchor:middle")
	canvas.Text(right, laby, "2\u03c0", "text-anchor:middle")
	canvas.Text(labx, 0, "1", "text-anchor:end")
	canvas.Text(labx, midy, "0", "text-anchor:end")
	canvas.Text(labx, h, "-1", "text-anchor:end")
	canvas.Gend()
	for x := 0.0; x < TwoPi; x += math.Pi / 25 {
		dx := int(vmap(x, 0, TwoPi, float64(left), float64(right)))
		dsy := int(vmap(math.Sin(x), -1, 1, 0, float64(h)))
		dcy := int(vmap(math.Cos(x), -1, 1, 0, float64(h)))
		canvas.Translate(0, h-height)
		canvas.Circle(dx, height-dsy, 3, "fill:red")
		canvas.Circle(dx, height-dcy, 3, "fill:blue")
		canvas.Gend()
	}
	canvas.Gend()
	canvas.End()
}
