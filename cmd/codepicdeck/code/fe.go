package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"os"
)
var canvas = svg.New(os.Stdout)
func main() {
	gutter, nc := 10, 2
	iw, ih := 200, 112
	pw, ph := (iw*nc)+gutter*(nc+1), (ih*3)+gutter*4
	canvas.Start(pw, ph)
	canvas.Def()
	canvas.Filter("f0")
	canvas.Saturate(1.0)
	canvas.Fend()
	canvas.Filter("f1")
	canvas.FeComponentTransfer()
	canvas.FeFuncTable("G", []float64{0, 0.5, 0.6, 0.85, 1.0})
	canvas.FeCompEnd()
	canvas.Fend()
	for i, b := 0, 0.0; b < 20.0; b += 2.0 {
		canvas.Filter(fmt.Sprintf("blur%d", i))
		canvas.Blur(b)
		canvas.Fend()
		i++
	}
	canvas.DefEnd()
	x, y := gutter, gutter
	canvas.Gstyle("text-anchor:middle;fill:white;font-family:sans-serif;font-size:24pt")
	for i, f := range []string{"f0", "f1", "blur1", "blur2"} {
		if i != 0 && i%nc == 0 {
			x = gutter
			y += ih + gutter
		}
		canvas.Image(x, y, iw, ih, "maple.jpg", "filter:url(#"+f+")")
		canvas.Text(x+iw/2, y+ih/2, f)
		x += iw + gutter
	}
	canvas.Gend()
	canvas.End()
}
