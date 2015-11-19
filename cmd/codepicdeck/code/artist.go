package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"os"
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
	bgfill := canvas.RGB(72, 45, 77)
	left := 50
	y := 40
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, bgfill)
	canvas.Gstyle("font-family:Roboto;fill:white")
	tf(left, y, "A MAN WHO WORKS WITH HIS HANDS IS A LABORER", 12)
	y += 70
	tf(left, y, "A MAN WHO", 52)
	y += 100
	tf(left, y, "WORKS", 87)
	y += 40
	tf(left, y, "WITH HIS HANDS AND HIS BRAIN IS A CRAFTSMAN", 13)
	y += 50
	tf(left, y, "BUT A MAN WHO", 37)
	y += 30
	tf(left, y, "WORKS WITH HIS HANDS AND HIS BRAIN", 15)
	y += 50
	tf(left, y, "AND HIS HEART IS", 35.4)
	y += 80
	tf(left, y, "AN ARTIST", 60)
	canvas.Gend()
	canvas.End()
}
