package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"math"
	"os"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

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
	canvas.TranslateRotate(xp, yp, 54)
	canvas.Polygon(xv, yv, style)
	canvas.Gend()
}

func main() {
	x, y, size := width/2.0, height/2.0, width*30/100
	textstyle := "%s;font-size:%dpx;font-family:sans-serif;text-anchor:middle"
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(240, 240, 240))
	canvas.Circle(x, y, width/2, canvas.RGB(255, 255, 255))
	star(x, y, 5, 90, 240, canvas.RGB(200, 200, 200))
	canvas.Text(x, y+size/3, "X", fmt.Sprintf(textstyle, canvas.RGB(127, 0, 0), size))
	canvas.End()
}