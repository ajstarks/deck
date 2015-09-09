package main

import (
	"fmt"
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
	tiles, maxstroke := 25, 10
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	linecaps := []string{"butt", "round", "square"}
	strokefmt := "stroke-width:%d"
	lcfmt := "stroke:black;stroke-linecap:%s"
	canvas.Gstyle(fmt.Sprintf(lcfmt, linecaps[rand.Intn(3)]))
	var sw string
	for y := 0; y < tiles; y++ {
		for x := 0; x < tiles; x++ {
			px := width / tiles * x
			py := height / tiles * y
			if rand.Intn(100) > 50 {
				sw = fmt.Sprintf(strokefmt, rand.Intn(maxstroke)+1)
				canvas.Line(px, py, px+width/tiles, py+height/tiles, sw)
			} else {
				sw = fmt.Sprintf(strokefmt, rand.Intn(maxstroke)+1)
				canvas.Line(px, py+height/tiles, px+width/tiles, py, sw)
			}
		}
	}
	canvas.Gend()
	canvas.End()
}
