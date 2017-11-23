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
func randcolor() string {
	return canvas.RGB(rand.Intn(255),rand.Intn(255),rand.Intn(255))
}
func rcube(x, y, l int) {
	l2, l3, l4, l6, l8 := l*2, l*3, l*4, l*6, l*8
	tx := []int{x, x + (l3), x, x - (l3), x}
	ty := []int{y, y + (l2), y + (l4), y + (l2), y}
	lx := []int{x - (l3), x, x, x - (l3), x - (l3)}
	ly := []int{y + (l2), y + (l4), y + (l8), y + (l6), y + (l2)}
	rx := []int{x + (l3), x + (l3), x, x, x + (l3)}
	ry := []int{y + (l2), y + (l6), y + (l8), y + (l4), y + (l2)}
	canvas.Polygon(tx, ty, randcolor())
	canvas.Polygon(lx, ly, randcolor())
	canvas.Polygon(rx, ry, randcolor())
}
func main() {
	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	xp, y := width/10, height/10
	n, hspace, vspace, size := 3, width/5, height/4, width/40
	for r := 0; r < n; r++ {
		for x := xp; x < width; x += hspace {
			rcube(x, y, size)
		}
		y += vspace
	}
	canvas.End()
}
