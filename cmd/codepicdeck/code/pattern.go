package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"os"
)
var (
	canvas = svg.New(os.Stdout)
	width, height = 500, 500
)

func main() {
	pct := 5
	pw, ph := (width*pct)/100, (height*pct)/100
	canvas.Start(width, height)

	// define the pattern
	canvas.Def()
	canvas.Pattern("hatch", 0, 0, pw, ph, "user")
	canvas.Gstyle("fill:none;stroke-width:1")
	canvas.Path(fmt.Sprintf("M0,0 l%d,%d", pw, ph), "stroke:red")
	canvas.Path(fmt.Sprintf("M%d,0 l-%d,%d", pw, pw, ph), "stroke:blue")
	canvas.Gend()
	canvas.PatternEnd()
	canvas.DefEnd()

	// use the pattern
	canvas.Gstyle("stroke:black; stroke-width:2")
	canvas.Circle(width/2, height/2, height/8, "fill:url(#hatch)")
	canvas.CenterRect((width*4)/5, height/2, height/4, height/4, "fill:url(#hatch)")
	canvas.Gend()
	canvas.End()
}
