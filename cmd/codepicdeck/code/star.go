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
func star(xp, yp, n int, inner, outer float64, style string) {
	xv, yv := make([]int, n*2), make([]int, n*2)
	angle := math.Pi / float64(n)
	for i := 0; i < n*2; i++ {
		fi := float64(i)
		if i%2 == 0 {
			xv[i] = int(math.Cos(angle*fi) * outer)
			yv[i] = int(math.Sin(angle*fi) * outer)
		} else {
			xv[i] = int(math.Cos(angle*fi) * inner)
			yv[i] = int(math.Sin(angle*fi) * inner)
		}
	}
	canvas.Translate(xp, yp)
	canvas.Polygon(xv, yv, style)
	canvas.Gend()
}

func main() {
	canvas.Start(width, height)
	for x, op, i := 50, 1.0, 5; i <= 10; i++ {
		star(x, 200, i*2, 20, 40, canvas.RGBA(0, 0, 127, op))
		star(x, 300, i, 20, 40, canvas.RGBA(127, 0, 127, op))
		x += 80
		op -= 0.15
	}
	canvas.End()
}
