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

func randarc(aw, ah, sw int, f1, f2 bool) {
	colors := []string{"red", "green", "blue", "gray"}
	afmt := "stroke:%s;stroke-opacity:%.2f;stroke-width:%dpx;fill:none"
	begin, arclength := rand.Intn(aw), rand.Intn(aw)
	end := begin + arclength
	baseline := ah / 2
	al, cl := arclength/2, len(colors)
	canvas.Arc(begin, baseline, al, al, 0, f1, f2, end, baseline,
		fmt.Sprintf(afmt, colors[rand.Intn(cl)], rand.Float64(), rand.Intn(sw)))
}

func main() {
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	aw := width / 2
	maxstroke := height / 10
	for i := 0; i < 20; i++ {
		randarc(aw, height, maxstroke, false, true)
		randarc(aw, height, maxstroke, false, false)
	}
	canvas.End()
}
