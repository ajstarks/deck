package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas   = svg.New(os.Stdout)
	width    = 500
	height   = 500
	maxlevel = 5
	colors   = []string{"red", "orange", "yellow", "green", "blue"}
)

func branch(x, y, r, level int) {
	astyle := fmt.Sprintf("fill:none;stroke:%s;stroke-width:%dpx", colors[level%maxlevel], level*2)
	canvas.Arc(x-r, y, r, r, 0, true, true, x+r, y, astyle)
	if level > 0 {
		branch(x-r, y+r/2, r/2, level-1)
		branch(x+r, y+r/2, r/2, level-1)
	}
}

// Example from "Generative Design", pg 414
func main() {
	canvas.Start(width, height)
	branch(0, 0, width/2, 6)
	canvas.End()
}
