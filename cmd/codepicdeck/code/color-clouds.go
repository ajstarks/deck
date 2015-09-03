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

func cloud(x, y, r int, style string) {
	small := r / 2
	medium := (r * 6) / 10
	canvas.Gstyle(style)
	canvas.Circle(x, y, r)
	canvas.Circle(x+r, y+small, small)
	canvas.Circle(x-r-small, y+small, small)
	canvas.Circle(x-r, y, medium)
	canvas.Rect(x-r-small, y, r*2+small, r)
	canvas.Gend()
}

func main() {
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	for i := 0; i < 50; i++ {
		red := rand.Intn(255)
		green := rand.Intn(255)
		blue := rand.Intn(255)
		size := rand.Intn(60)
		x := rand.Intn(width)
		y := rand.Intn(height)
		cloud(x, y, size, canvas.RGB(red, green, blue))
	}
	canvas.End()
}
