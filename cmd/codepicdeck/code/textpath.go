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

func coord(x, y, size int) {
	offset := size * 2
	canvas.Text(x, y-offset, fmt.Sprintf("(%d,%d)", x, y),
		"font-size:50%;text-anchor:middle")
	canvas.Circle(x, y, size, "fill-opacity:0.3")
}

func makepath(x, y, sx, sy, cx, cy, ex, ey int, id, text string) {
	canvas.Def()
	canvas.Qbez(sx, sy, cx, cy, ex, ey, `id="`+id+`"`)
	canvas.DefEnd()
	canvas.Translate(x, y)
	canvas.Textpath(text, "#"+id)
	coord(sx, sy, 5)
	coord(ex, ey, 5)
	coord(cx, cy, 5)
	canvas.Gend()
}

func main() {
	message := `It's fine & "dandy" to have text on a path`
	canvas.Start(width, height)
	canvas.Gstyle("font-family:serif;font-size:21pt")
	makepath(0, 0, 70, 200, 100, 425, 425, 125, "tpath", message)
	canvas.Gend()
	canvas.End()
}
