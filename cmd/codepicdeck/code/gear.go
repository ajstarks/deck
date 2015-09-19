package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func gear(x, y, w, h, n, l, m int, color string) {
	canvas.Gstyle(fmt.Sprintf("fill:none;stroke:%s;stroke-width:%d", color, n/2))
	canvas.Circle(x+w/2, y+h/2, n)
	canvas.Circle(x+w/2, y+h/2, n/5, "fill:"+color)
	ai := 360 / float64(m)
	for a := 0.0; a <= 360.0; a += ai {
		canvas.TranslateRotate(x+w/2, y+h/2, a)
		canvas.Line(n-l, n-l, n+l, n+l)
		canvas.Gend()
	}
	canvas.Gend()
}

func main() {
	canvas.Start(width, height)
	gear(0, 0, 250, 250, 60, 10, 8, "black")
	gear(100, 160, 250, 250, 60, 10, 8, "red")
	gear(300, 140, 100, 100, 20, 6, 8, "blue")
	canvas.End()
}
