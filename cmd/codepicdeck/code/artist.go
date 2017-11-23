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

func tf(x, y int, s string, size float64) {
	canvas.Text(x, y, s, fmt.Sprintf("font-size:%gpt", size))
}

func main() {
	x, y := width/2, 35
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(72, 45, 77))
	canvas.Gstyle("font-family:Roboto;fill:white;text-anchor:middle")
	tf(x, y, "A MAN WHO WORKS WITH HIS HANDS IS A LABORER", 14)
	y += 70
	tf(x, y, "A MAN WHO", 60)
	y += 105
	tf(x, y, "WORKS", 90)
	y += 35
	tf(x, y, "WITH HIS HANDS AND HIS BRAIN IS A CRAFTSMAN", 15)
	y += 60
	tf(x, y, "BUT A MAN WHO", 42)
	y += 40
	tf(x, y, "WORKS WITH HIS HANDS AND HIS BRAIN", 16)
	y += 55
	tf(x, y, "AND HIS HEART IS", 36)
	y += 85
	tf(x, y, "AN ARTIST", 64)
	canvas.Gend()
	canvas.End()
}
