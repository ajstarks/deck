package main

import (
	"github.com/ajstarks/svgo"
	"math/rand"
	"os"
)

func tloc(x, y, s int, r, d float64) (int, int) {
	fx, fy, fs := float64(x), float64(y), float64(s)
	padding := 2 * fs
	return int(padding + (fx * fs) - (.5 * fs) + (r * d)),
		int(padding + (fy * fs) - (.5 * fs) + (r * d))
}

func random(n float64) float64 {
	x := rand.Float64()
	if x < 0.5 {
		return -n * x
	}
	return n * x
}

func main() {
	columns, rows, sqrsize := 12, 12, 32
	rndStep, dampen := .22, 0.45
	width, height := (columns+4)*sqrsize, (rows+4)*sqrsize
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	canvas.Gstyle("fill:rgb(0,0,127);fill-opacity:0.3")
	for y, randsum := 1, 0.0; y <= rows; y++ {
		randsum += float64(y) * rndStep
		for x := 1; x <= columns; x++ {
			tx, ty := tloc(x, y, sqrsize, random(randsum), dampen)
			canvas.TranslateRotate(tx, ty, random(randsum))
			canvas.CenterRect(0, 0, sqrsize, sqrsize)
			canvas.Gend()
		}
	}
	canvas.Gend()
	canvas.End()
}
