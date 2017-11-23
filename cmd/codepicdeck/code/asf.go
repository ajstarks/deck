package main

import (
	"github.com/ajstarks/svgo/float"
	"math"
	"os"
)

var canvas = svg.New(os.Stdout)

const PI = math.Pi

func polar(r, theta float64) (float64, float64) {
	return r * math.Cos(theta), r * math.Sin(theta)
}
func star(x, y float64, n int, inner, outer float64, style string) {
	xv, yv := make([]float64, n*2), make([]float64, n*2)
	angle := PI / float64(n)
	for i := 0; i < n*2; i++ {
		fi := float64(i)
		if i%2 == 0 {
			xv[i] = (math.Cos(angle*fi) * outer)
			yv[i] = (math.Sin(angle*fi) * outer)
		} else {
			xv[i] = (math.Cos(angle*fi) * inner)
			yv[i] = (math.Sin(angle*fi) * inner)
		}
	}
	canvas.Translate(x, y)
	canvas.Polygon(xv, yv, style)
	canvas.Gend()
}
func main() {
	w, h, cx, cy := 500.0, 500.0, w/2, h/2
	canvas.Start(w, h)
	canvas.Circle(cx, cy, w*.4, "fill:red")
	canvas.Gstyle("fill:white")
	for t := PI / 6; t < PI*2; t += PI / 3 {
		x, y := polar(w*.2, t)
		canvas.Circle(cx+x, cy+y, 65)
	}
	canvas.Gend()
	canvas.TranslateRotate(cx, cy, 55)
	star(0, 0, 5, 45, 120, "fill:blue")
	canvas.Gend()
	canvas.End()
}
