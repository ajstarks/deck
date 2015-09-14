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

func main() {
	canvas.Start(width, height)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		fill := canvas.RGBA(
			rand.Intn(255),
			rand.Intn(255),
			rand.Intn(255),
			rand.Float64())
		canvas.Rect(
			rand.Intn(width),
			rand.Intn(height),
			rand.Intn(100),
			rand.Intn(100),
			fill)
	}
	canvas.End()
}
