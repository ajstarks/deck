package main

import (
	"crypto/md5"
	"os"

	"github.com/ajstarks/svgo"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func colorhash(s string) (int, int, int) {
	hash := md5.New()
	hash.Write([]byte(s))
	v := hash.Sum(nil)
	return int(v[0]), int(v[1]), int(v[2])
}

func main() {
	name := "SVG is cool"
	style := "font-family:sans-serif;fill:white;text-anchor:middle"
	r, g, b := colorhash(name)
	canvas.Start(width, height)
	canvas.Gstyle(style)
	canvas.Rect(0, 0, width, height, canvas.RGB(r, g, b))
	canvas.Text(width/2, height/2, name, "font-size:36pt")
	canvas.Gend()
	canvas.End()
}
