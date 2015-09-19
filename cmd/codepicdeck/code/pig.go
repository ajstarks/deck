package main

import (
	"github.com/ajstarks/svgo"
	"os"
)

var canvas = svg.New(os.Stdout)

func defilter(id string) {

	cm := [20]float64{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 5, -4}
	gb1 := svg.Filterspec{Result: "result1"}
	gb2 := svg.Filterspec{In: "result4", Result: "result5"}
	cm1 := svg.Filterspec{Result: "result2"}
	cp1 := svg.Filterspec{In: "result2", In2: "result2", Result: "result3"}
	cp2 := svg.Filterspec{In: "result1", In2: "result3", Result: "result4"}
	cp3 := svg.Filterspec{In: "result6", In2: "result4", Result: "result7"}
	cp4 := svg.Filterspec{In: "result4", In2: "result7", Result: "result8"}
	cp5 := svg.Filterspec{In2: "result8", Result: "result9"}
	sps := svg.Filterspec{In: "result5", Result: "result6"}
	bls := svg.Filterspec{In: "result9", In2: "result9"}

	canvas.Filter(id)
	canvas.FeGaussianBlur(gb1, 5, 5)
	canvas.FeColorMatrix(cm1, cm)
	canvas.FeComposite(cp1, "atop", 0, 0, 0, 0)
	canvas.FeComposite(cp2, "in", 0, 0, 0, 0)
	canvas.FeGaussianBlur(gb2, 5, 5)
	canvas.FeSpecularLighting(sps, 2, 2.5, 55, "white")
	canvas.FeDistantLight(svg.Filterspec{}, 255, 60)
	canvas.FeSpecEnd()
	canvas.FeComposite(cp3, "in", 0, 0, 0, 0)
	canvas.FeComposite(cp4, "arithmetic", 0, 1, 1, 0)
	canvas.FeComposite(cp5, "in", 0, 0, 0, 0)
	canvas.FeBlend(bls, "multiply")
	canvas.Fend()
}

func main() {
	width := 600
	height := 600
	canvas.Start(width, height)
	id := "ink"
	canvas.Def()
	defilter(id)
	canvas.Gid("pic")
	canvas.Ellipse(0, 0, 100, 50)
	canvas.Gend()
	canvas.DefEnd()

	canvas.Use(width/2, height-100, "#pic", `fill="red"`)
	canvas.Use(width/2, height/2, "#pic", `fill="red"`, `filter="url(#ink)"`)
	canvas.End()
}
