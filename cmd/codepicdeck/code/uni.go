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

func main() {
	top, left, fontsize := 50, 100, 16
	xoffset, yoffset := 25, 25
	rows, cols := 16, 16
	glyph := 0x2600
	font := "Symbola"
	stylefmt := "font-family:%s;font-size:%dpx;text-anchor:middle"
	canvas.Start(width, height)
	canvas.Gstyle(fmt.Sprintf(stylefmt, font, fontsize))

	x, y := left, top
	for r := 0; r < rows; r++ {
		canvas.Text(x-yoffset, y, fmt.Sprintf("%X", glyph), "text-anchor:end;fill:gray")
		for c := 0; c < cols; c++ {
			if r == 0 {
				canvas.Text(x, y-yoffset, fmt.Sprintf("%02X", c), "fill:gray")
			}
			canvas.Text(x, y, string(glyph))
			glyph++
			x += xoffset
		}
		x = left
		y += yoffset
	}
	canvas.Gend()
	canvas.End()
}
