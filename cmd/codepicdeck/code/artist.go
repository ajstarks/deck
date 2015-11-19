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
	canvas.Text(x, y, s,fmt.Sprintf("font-size:%gpt",size))
}

func main() {
	x, y := 50, 40
	canvas.Start(width, height)
	canvas.Rect(0,0,width,height, canvas.RGB(72, 45, 77))
	canvas.Gstyle("font-family:Roboto;fill:white")
	tf(x, y, "A MAN WHO WORKS WITH HIS HANDS IS A LABORER",12); y += 70
	tf(x, y, "A MAN WHO", 52); y += 100
	tf(x, y, "WORKS", 87); y += 40
	tf(x, y, "WITH HIS HANDS AND HIS BRAIN IS A CRAFTSMAN", 13); y += 50
	tf(x, y, "BUT A MAN WHO",37); y += 30
	tf(x, y, "WORKS WITH HIS HANDS AND HIS BRAIN", 15); y += 50
	tf(x, y, "AND HIS HEART IS", 35.4); y += 80
	tf(x, y, "AN ARTIST", 60)
	canvas.Gend()
	canvas.End()
}