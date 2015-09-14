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
	name := "SVGo"
	style := "fill:white;text-anchor:middle;font-size:72pt"
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(colorhash(name)))
	canvas.Text(width/2, height/2, name, style)
	canvas.End()
}
