package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

// inspired by Gerhard Richter's 256 colors, 1974
func main() {
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)

	w, h, gutter := 24, 18, 5
	rows, cols := 16, 16
	top, left := 20, 20

	for r, x := 0, left; r < rows; r++ {
		for c, y := 0, top; c < cols; c++ {
			canvas.Rect(x, y, w, h,
				canvas.RGB(rand.Intn(255), rand.Intn(255), rand.Intn(255)))
			y += (h + gutter)
		}
		x += (w + gutter)
	}
	canvas.End()
}
