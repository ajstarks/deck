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

func branch(x, y, r, level int) {
	astyle := fmt.Sprintf("fill:none;stroke:rgb(0,130,164);stroke-width:%dpx", level*2)
	canvas.Arc(x-r, y, r, r, 0, true, true, x+r, y, astyle)
	if level > 0 {
		branch(x-r, y+r/2, r/2, level-1)
		branch(x+r, y+r/2, r/2, level-1)
	}
}

// Example from "Generative Design", pg 414
func main() {
	canvas.Start(width, height)
	branch(0, 0, 250, 6)
	canvas.End()
}
