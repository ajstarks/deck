package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width = 500
	height = 500
)

// inspired by Gerhard Richter's 256 colors, 1974
func main() {
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	rw := 32
	rh := 18
	margin := 5
	for i, x := 0, 20; i < 16; i++ {
		x += (rw + margin)
		for j, y := 0, 20; j < 16; j++ {
			canvas.Rect(x, y, rw, rh, canvas.RGB(rand.Intn(255), rand.Intn(255), rand.Intn(255)))
			y += (rh + margin)
		}
	}
	canvas.End()
}
