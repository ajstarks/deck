// vgdeck: slide decks for OpenVG
package main

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck"
	"github.com/ajstarks/openvg"
)

var wintrans, _ = charset.TranslatorTo("windows-1252")
var codemap = strings.NewReplacer("\t", "    ")

// dodeck sets up the graphics environment and kicks off the interaction
func dodeck(filename string, slidenum int) {
	w, h := openvg.Init()
	defer openvg.Finish()
	d, err := deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	openvg.FillRGB(200, 200, 200, 1)
	openvg.Rect(0, 0, openvg.VGfloat(w), openvg.VGfloat(h))
	showslide(d, slidenum)
	openvg.SaveEnd(fmt.Sprintf("%s-slide-%04d", filename, slidenum))	
}


// pct computes percentages
func pct(p float64, m openvg.VGfloat) openvg.VGfloat {
	return openvg.VGfloat((p / 100.0)) * m
}

func pctwidth(p float64, p1, p2 openvg.VGfloat) openvg.VGfloat {
	pw := deck.Pwidth(p, float64(p1), float64(p2))
	return openvg.VGfloat(pw)

}

func fromUTF8(s string) string {
	_, b, err := wintrans.Translate([]byte(s), true)
	if err != nil {
		return s
	}
	return string(b)
}

//showtext displays text
func showtext(x, y openvg.VGfloat, s, align, font string, fs openvg.VGfloat) {
	t := fromUTF8(s)
	fontsize := int(fs)
	switch align {
	case "center", "middle", "mid", "c":
		openvg.TextMid(x, y, t, font, fontsize)
	case "right", "end", "e":
		openvg.TextEnd(x, y, t, font, fontsize)
	default:
		openvg.Text(x, y, t, font, fontsize)
	}
}

// dimen returns device dimemsion from percentages
func dimen(d deck.Deck, x, y, s float64) (xo, yo, so openvg.VGfloat) {
	xf, yf, sf := deck.Dimen(d.Canvas, x, y, s)
	xo, yo, so = openvg.VGfloat(xf), openvg.VGfloat(yf), openvg.VGfloat(sf)*0.8
	return
}

// showlide displays slides
func showslide(d deck.Deck, n int) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	slide := d.Slide[n]
	if slide.Bg == "" {
		slide.Bg = "white"
	}
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	openvg.Start(d.Canvas.Width, d.Canvas.Height)
	cw := openvg.VGfloat(d.Canvas.Width)
	ch := openvg.VGfloat(d.Canvas.Height)
	openvg.FillColor(slide.Bg)
	openvg.Rect(0, 0, cw, ch)
	var x, y, fs openvg.VGfloat

	// every image in the slide
	for _, im := range slide.Image {
		x = pct(im.Xp, cw)
		y = pct(im.Yp, ch)
		midx := openvg.VGfloat(im.Width / 2)
		midy := openvg.VGfloat(im.Height / 2)
		openvg.Image(x-midx, y-midy, im.Width, im.Height, im.Name)
		if len(im.Caption) > 0 {
			capfs := pctwidth(im.Sp, cw, cw/100)
			if im.Font == "" {
				im.Font = "sans"
			}
			if im.Color == "" {
				openvg.FillColor(slide.Fg)
			} else {
				openvg.FillColor(im.Color)
			}
			if im.Align == "" {
				im.Align = "center"
			}
			switch im.Align {
			case "left", "start":
				x -= midx
			case "right", "end":
				x += midx
			}
			showtext(x, y-((midy)+(capfs*2.0)), im.Caption, im.Align, im.Font, capfs)
		}
	}

	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	const defaultSw = 1.5
	var strokeopacity float64
	// line
	for _, line := range slide.Line {
		if line.Color == "" {
			line.Color = slide.Fg // defaultColor
		}
		if line.Opacity == 0 {
			strokeopacity = 1
		} else {
			strokeopacity = line.Opacity / 100.0
		}
		x1, y1, sw := dimen(d, line.Xp1, line.Yp1, line.Sp)
		x2, y2, _ := dimen(d, line.Xp2, line.Yp2, 0)
		openvg.StrokeColor(line.Color, openvg.VGfloat(strokeopacity))
		if sw == 0 {
			sw = defaultSw
		}
		openvg.StrokeWidth(openvg.VGfloat(sw))
		openvg.StrokeColor(line.Color)
		openvg.Line(x1, y1, x2, y2)
		openvg.StrokeWidth(0)
	}
	// ellipse
	for _, ellipse := range slide.Ellipse {
		x, y, _ = dimen(d, ellipse.Xp, ellipse.Yp, 0)
		var w, h openvg.VGfloat
		w = pct(ellipse.Wp, cw)
		if ellipse.Hr == 0 { // if relative height not specified, base height on overall height
			h = pct(ellipse.Hp, ch)
		} else {
			h = pct(ellipse.Hr, w)
		}
		if ellipse.Color == "" {
			ellipse.Color = defaultColor
		}
		if ellipse.Opacity == 0 {
			ellipse.Opacity = 1
		} else {
			ellipse.Opacity /= 100
		}
		openvg.FillColor(ellipse.Color, openvg.VGfloat(ellipse.Opacity))
		openvg.Ellipse(x, y, w, h)
	}
	// rect
	for _, rect := range slide.Rect {
		x, y, _ = dimen(d, rect.Xp, rect.Yp, 0)
		var w, h openvg.VGfloat
		w = pct(rect.Wp, cw)
		if rect.Hr == 0 { // if relative height not specified, base height on overall height
			h = pct(rect.Hp, ch)
		} else {
			h = pct(rect.Hr, w)
		}
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		if rect.Opacity == 0 {
			rect.Opacity = 1
		} else {
			rect.Opacity /= 100
		}
		openvg.FillColor(rect.Color, openvg.VGfloat(rect.Opacity))
		openvg.Rect(x-(w/2), y-(h/2), w, h)
	}
	// curve
	for _, curve := range slide.Curve {
		if curve.Color == "" {
			curve.Color = defaultColor
		}
		if curve.Opacity == 0 {
			strokeopacity = 1
		} else {
			strokeopacity = curve.Opacity / 100.0
		}
		x1, y1, sw := dimen(d, curve.Xp1, curve.Yp1, curve.Sp)
		x2, y2, _ := dimen(d, curve.Xp2, curve.Yp2, 0)
		x3, y3, _ := dimen(d, curve.Xp3, curve.Yp3, 0)
		openvg.StrokeColor(curve.Color, openvg.VGfloat(strokeopacity))
		openvg.FillColor(slide.Bg, openvg.VGfloat(curve.Opacity))
		if sw == 0 {
			sw = defaultSw
		}
		openvg.StrokeWidth(sw)
		openvg.Qbezier(x1, y1, x2, y2, x3, y3)
		openvg.StrokeWidth(0)
	}

	// arc
	for _, arc := range slide.Arc {
		if arc.Color == "" {
			arc.Color = defaultColor
		}
		if arc.Opacity == 0 {
			strokeopacity = 1
		} else {
			strokeopacity = arc.Opacity / 100.0
		}
		ax, ay, sw := dimen(d, arc.Xp, arc.Yp, arc.Sp)
		w := pct(arc.Wp, cw)
		h := pct(arc.Hp, cw)
		openvg.StrokeColor(arc.Color, openvg.VGfloat(strokeopacity))
		openvg.FillColor(slide.Bg, openvg.VGfloat(arc.Opacity))
		if sw == 0 {
			sw = defaultSw
		}
		openvg.StrokeWidth(sw)
		openvg.Arc(ax, ay, w, h, openvg.VGfloat(arc.A1), openvg.VGfloat(arc.A2))
		openvg.StrokeWidth(0)
	}

	// polygon
	for _, poly := range slide.Polygon {
		if poly.Color == "" {
			poly.Color = defaultColor
		}
		if poly.Opacity == 0 {
			poly.Opacity = 1
		} else {
			poly.Opacity /= 100
		}
		xs := strings.Split(poly.XC, " ")
		ys := strings.Split(poly.YC, " ")
		if len(xs) != len(ys) {
			continue
		}
		if len(xs) < 3 || len(ys) < 3 {
			continue
		}
		px := make([]openvg.VGfloat, len(xs))
		py := make([]openvg.VGfloat, len(ys))
		for i := 0; i < len(xs); i++ {
			x, err := strconv.ParseFloat(xs[i], 32)
			if err != nil {
				px[i] = 0
			} else {
				px[i] = pct(x, cw)
			}
			y, err := strconv.ParseFloat(ys[i], 32)
			if err != nil {
				py[i] = 0
			} else {
				py[i] = pct(y, ch)
			}
		}
		openvg.FillColor(poly.Color, openvg.VGfloat(poly.Opacity))
		openvg.Polygon(px, py)
	}

	openvg.FillColor(slide.Fg)

	// every list in the slide
	var offset, textopacity openvg.VGfloat
	const blinespacing = 2.4
	for _, l := range slide.List {
		if l.Font == "" {
			l.Font = "sans"
		}
		x, y, fs = dimen(d, l.Xp, l.Yp, l.Sp)
		if l.Type == "bullet" {
			offset = 1.2 * fs
		} else {
			offset = 0
		}
		if l.Opacity == 0 {
			textopacity = 1
		} else {
			textopacity = openvg.VGfloat(l.Opacity / 100)
		}
		// every list item
		var li, lifont string
		for ln, tl := range l.Li {
			if len(l.Color) > 0 {
				openvg.FillColor(l.Color, textopacity)
			} else {
				openvg.FillColor(slide.Fg)
			}
			if l.Type == "bullet" {
				boffset := fs / 2
				openvg.Ellipse(x, y+boffset, boffset, boffset)
				//openvg.Rect(x, y+boffset/2, boffset, boffset)
			}
			if l.Type == "number" {
				li = fmt.Sprintf("%d. ", ln+1) + tl.ListText
			} else {
				li = tl.ListText
			}
			if len(tl.Color) > 0 {
				openvg.FillColor(tl.Color, textopacity)
			}
			if len(tl.Font) > 0 {
				lifont = tl.Font
			} else {
				lifont = l.Font
			}
			showtext(x+offset, y, li, l.Align, lifont, fs)
			y -= fs * blinespacing
		}
	}
	openvg.FillColor(slide.Fg)

	// every text in the slide
	const linespacing = 1.8

	var tdata string
	for _, t := range slide.Text {
		if t.File != "" {
			tdata = includefile(t.File)
		} else {
			tdata = t.Tdata
		}
		if t.Font == "" {
			t.Font = "sans"
		}
		if t.Opacity == 0 {
			textopacity = 1
		} else {
			textopacity = openvg.VGfloat(t.Opacity / 100)
		}
		x, y, fs = dimen(d, t.Xp, t.Yp, t.Sp)
		td := strings.Split(tdata, "\n")
		if t.Type == "code" {
			t.Font = "mono"
			tdepth := ((fs * linespacing) * openvg.VGfloat(len(td))) + fs
			openvg.FillColor("rgb(240,240,240)")
			openvg.Rect(x-20, y-tdepth+(fs*linespacing), pctwidth(t.Wp, cw, cw-x-20), tdepth)
		}
		if t.Color == "" {
			openvg.FillColor(slide.Fg, textopacity)
		} else {
			openvg.FillColor(t.Color, textopacity)
		}
		if t.Type == "block" {
			textwrap(x, y, pctwidth(t.Wp, cw, cw/2), tdata, t.Font, fs, fs*linespacing, 0.3)
		} else {
			// every text line
			for _, txt := range td {
				showtext(x, y, txt, t.Align, t.Font, fs)
				y -= (fs * linespacing)
			}
		}
	}
	openvg.FillColor(slide.Fg)
	openvg.End()
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '-'
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(x, y, w openvg.VGfloat, s string, font string, fs, leading, factor openvg.VGfloat) {
	size := int(fs)
	if font == "mono" {
		factor = 1.0
	}
	wordspacing := openvg.TextWidth("m", font, size)
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw := openvg.TextWidth(s, font, size)
		openvg.Text(xp, yp, s, font, size)
		xp += tw + (wordspacing * factor)
		if xp > edge {
			xp = x
			yp -= leading
		}
	}
}

// includefile returns the contents of a file as string
func includefile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}


// for every file, make a deck
func main() {
	var slidenum = flag.Int("slide", 0, "initial slide")
	flag.Parse()
	for _, f := range flag.Args() {
		dodeck(f, *slidenum)

	}
}
