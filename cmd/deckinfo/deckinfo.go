// deckinfo: count deck elements
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ajstarks/deck"
)

func show(name string, value int) {
	if value > 0 {
		fmt.Printf("%s %d\n", name, value)
	}
}

var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

var ltype = map[string]string{"": "", "bullet": "b", "number": "n"}

func showitems(s deck.Slide, n int) {
	fmt.Printf("\n// slide %d\nslide ", n+1)
	if len(s.Bg) > 0 {
		fmt.Printf("%q", s.Bg)
	}
	if len(s.Fg) > 0 {
		fmt.Printf(" %q", s.Fg)
	}
	fmt.Println()

	for _, i := range s.Text {
		var textalign string
		var font, color string
		var opacity float64
		switch i.Align {
		case "center", "c", "middle":
			textalign = "c"
		case "end", "right", "r", "e":
			textalign = "e"
		}
		if i.Font == "" {
			font = "sans"
		} else {
			font = i.Font
		}
		if i.Color == "" {
			color = s.Fg
		} else {
			color = i.Color
		}
		if i.Opacity == 0 {
			opacity = 100
		} else {
			opacity = i.Opacity
		}

		if i.Type == "block" {
			if i.Wp == 0 {
				i.Wp = 50
			}
			fmt.Printf("\ttextblock\t%q %v %v %v %v %q %q %v\n", xmlesc(i.Tdata), i.Xp, i.Yp, i.Wp, i.Sp, font, color, opacity)
		} else {
			if len(i.File) > 0 {
				fmt.Printf("\ttextfile\t%q %v %v %v %q %q %v\n", i.File, i.Xp, i.Yp, i.Sp, font, color, opacity)
			} else {
				fmt.Printf("\t%stext\t%q %v %v %v %q %q %v\n", textalign, xmlesc(i.Tdata), i.Xp, i.Yp, i.Sp, font, color, opacity)
			}
		}
	}
	var scale float64
	for _, i := range s.Image {
		if i.Scale == 0 {
			scale = 100.0
		} else {
			scale = i.Scale
		}
		if len(i.Caption) == 0 {
			fmt.Printf("\timage\t%q %v %v %d %d %v\n", i.Name, i.Xp, i.Yp, i.Width, i.Height, scale)
		} else {
			fmt.Printf("\tcimage\t%q %q %v %v %d %d %v\n", i.Name, xmlesc(i.Caption), i.Xp, i.Yp, i.Width, i.Height, scale)
		}
	}
	for _, i := range s.List {
		fmt.Printf("\t%slist\t%v %v %v\n", ltype[i.Type], i.Xp, i.Yp, i.Sp)
		for _, li := range i.Li {
			fmt.Printf("\t\tli %q\n", xmlesc(li.ListText))
		}
		fmt.Println("\telist")
	}
	for _, i := range s.Line {
		fmt.Printf("\tline\t%v %v %v %v\n", i.Xp1, i.Yp1, i.Xp2, i.Yp2)
	}
	for _, i := range s.Curve {
		fmt.Printf("\tcurve\t%v %v %v %v %v %v\n", i.Xp1, i.Yp1, i.Xp2, i.Yp2, i.Xp3, i.Yp3)
	}
	for _, i := range s.Arc {
		fmt.Printf("\tarc\t%v %v %v %v %v %v\n", i.Xp, i.Yp, i.Wp, i.Hp, i.A1, i.A2)
	}
	for _, i := range s.Polygon {
		fmt.Printf("\tpolygon\t%q %q\n", i.XC, i.YC)
	}
	fmt.Println("eslide")
}

func main() {
	var showit = flag.Bool("v", false, "verbose")
	flag.Parse()
	for _, file := range flag.Args() {
		d, err := deck.Read(file, 0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
			continue
		}
		fmt.Printf("// Elements of %s\n", file)
		if *showit {
			fmt.Println("deck")
		}
		var texts, images, lists, arcs, lines, ellipses, rects, curves, polygons, links int
		show("// slide count", len(d.Slide))
		for ns, s := range d.Slide {
			if *showit {
				showitems(s, ns)
			}
			texts += len(s.Text)
			images += len(s.Image)
			lists += len(s.List)
			lines += len(s.Line)
			rects += len(s.Rect)
			arcs += len(s.Arc)
			ellipses += len(s.Ellipse)
			curves += len(s.Curve)
			polygons += len(s.Polygon)
		}

		show("// text", texts)
		show("// image", images)
		show("// link", links)
		show("// list", lists)
		show("// line", lines)
		show("// rect", rects)
		show("// ellipse", ellipses)
		show("// arc", arcs)
		show("// curve", curves)
		show("// polygon", polygons)
	}
	if *showit {
		fmt.Println("edeck")
	}
}
