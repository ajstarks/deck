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

// See: http://vormplus.be/blog/article/processing-month-day-4-stars
func stars(n int, inner, outer float64) ([]int, []int) {
	xv := make([]int, n*2)
	yv := make([]int, n*2)
	angle := math.Pi / float64(n)
	var x, y float64
	for i := 0; i < n*2; i++ {
		fi := float64(i)
		if i%2 == 0 {
			x = math.Cos(angle*fi) * outer
			y = math.Sin(angle*fi) * outer
		} else {
			x = math.Cos(angle*fi) * inner
			y = math.Sin(angle*fi) * inner
		}
		xv[i] = int(x)
		yv[i] = int(y)
	}
	return xv, yv
}

func main() {
	canvas.Start(width, height)
	canvas.Translate(width/2, height/2)
	canvas.Polygon(stars(10, 50, 200))
	canvas.Gend()
	canvas.End()
}
