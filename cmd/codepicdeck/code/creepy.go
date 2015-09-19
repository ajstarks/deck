package main

import (
	"os"

	"github.com/ajstarks/svgo"
)

var canvas = svg.New(os.Stdout)

func smile(x, y, r int) {
	r2 := r * 2
	r3 := r * 3
	r4 := r * 4
	rq := r / 4
	gray := canvas.RGB(200, 200, 200)
	red := canvas.RGB(127, 0, 0)
	canvas.Roundrect(x-r2, y-r2, r*7, r*20, r2, r2, gray)
	canvas.Circle(x, y, r, red)
	canvas.Circle(x, y, rq, "fill:white")
	canvas.Circle(x+r3, y, r)
	canvas.Arc(x-r, y+r3, rq, rq, 0, true, false, x+r4, y+r3)
}

func main() {
	canvas.Start(500, 500)
	canvas.Rect(0, 0, 500, 500, "fill:white")
	smile(200, 100, 10)
	canvas.Gtransform("rotate(30)")
	smile(200, 100, 10)
	canvas.Gend()
	canvas.Gtransform("translate(50,0) scale(2,2)")
	canvas.Gstyle("opacity:0.5")
	smile(200, 100, 30)
	canvas.Gend()
	canvas.Gend()
	canvas.End()
}
