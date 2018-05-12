// svgdeck: make SVG slide decks
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajstarks/deck"
	"github.com/ajstarks/svgo/float"
)

const (
	mm2pt   = 2.83464 // mm to pt conversion
	namefmt = "%s-%05d.svg"
)

// PageDimen describes page dimensions
// the unit field is used to convert to pt.
type PageDimen struct {
	width, height, unit float64
}

// fontmap maps generic font names to specific implementation names
var fontmap = map[string]string{}

// pagemap defines page dimensions
var pagemap = map[string]PageDimen{
	"Letter":     {792, 612, 1},
	"Legal":      {1008, 612, 1},
	"Tabloid":    {1224, 792, 1},
	"ArchA":      {864, 648, 1},
	"Widescreen": {1152, 648, 1},
	"4R":         {432, 288, 1},
	"Index":      {360, 216, 1},
	"A2":         {420, 594, mm2pt},
	"A3":         {420, 297, mm2pt},
	"A4":         {297, 210, mm2pt},
	"A5":         {210, 148, mm2pt},
}

var codemap = strings.NewReplacer("\t", "    ")

// grid makes a labeled grid
func grid(doc *svg.SVG, w, h float64, color string, percent float64) {
	pw := w * (percent / 100)
	ph := h * (percent / 100)
	fs := pct(1, w)
	pl := 0.0
	doc.Gstyle(fmt.Sprintf("fill:%s;font-family:%s;font-size:%.2f;text-anchor:center", color, fontlookup("sans"), fs))
	for x := 0.0; x <= w; x += pw {
		doc.Line(x, 0, x, h, "stroke-width:0.5; stroke:"+color)
		if pl > 0 {
			doc.Text(x, h-fs, fmt.Sprintf("%.0f", pl))
		}
		pl += percent
	}
	pl = 0.0
	for y := 0.0; y <= h; y += ph {
		doc.Line(0, y, w, y, "stroke-width:0.5; stroke:"+color)
		if pl < 100 {
			doc.Text(fs, y+(fs/3), fmt.Sprintf("%.0f", 100-pl))
		}
		pl += percent
	}
	doc.Gend()
}

// pct converts percentages to canvas measures
func pct(p float64, m float64) float64 {
	return ((p / 100.0) * m)
}

// radians converts degrees to radians
func radians(deg float64) float64 {
	return (deg * math.Pi) / 180.0
}

// polar returns the euclidian corrdinates from polar coordinates
func polar(x, y, r, angle float64) (float64, float64) {
	px := (r * math.Cos(radians(angle))) + x
	py := (r * math.Sin(radians(angle))) + y
	return px, py
}

// dimen returns canvas dimensions from percentages
func dimen(w, h float64, xp, yp, sp float64) (float64, float64, float64) {
	return pct(xp, w), pct(100-yp, h), pct(sp, w)
}

// setop sets the alpha value:
// 0 == default value (opaque)
// -1 == fully transparent
// > 0 set opacity percent
func setop(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 0:
		return v / 100
	case v == 0:
		return 1
	}
	return v
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// fontlookup maps font aliases to implementation font names
func fontlookup(s string) string {
	font, ok := fontmap[s]
	if ok {
		return font
	}
	return "sans"
}

// bullet draws a bullet
func bullet(doc *svg.SVG, x, y, size float64, color string) {
	rs := size / 2
	doc.Circle(x-size, y-(rs*2)/3, rs/2, "fill:"+color)
	// dorect(doc, x-size, y-rs-(rs/2), rs, rs, color, 0)
}

// background places a colored rectangle
func background(doc *svg.SVG, w, h float64, color string) {
	dorect(doc, 0, 0, w, h, color, 0)
}

// doline draws a line
func doline(doc *svg.SVG, xp1, yp1, xp2, yp2, sw float64, color string, opacity float64) {
	doc.Line(xp1, yp1, xp2, yp2, fmt.Sprintf("strokewidth:%.2f;stroke:%s;stroke-opacity:%.2f", sw, color, setop(opacity)))
}

// doarc draws a line
func doarc(doc *svg.SVG, x, y, w, h, a1, a2, sw float64, color string, opacity float64) {
	sx, sy := polar(x, y, w, -a1)
	ex, ey := polar(x, y, h, -a2)
	doc.Arc(sx, sy, w, h, 0, false, false, ex, ey, fmt.Sprintf("fill:none;strokewidth:%.2f;stroke:%s;stroke-opacity:%.2f", sw, color, setop(opacity)))
}

// docurve draws a bezier curve
func docurve(doc *svg.SVG, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string, opacity float64) {
	doc.Qbez(xp1, yp1, xp2, yp2, xp3, yp3, fmt.Sprintf("fill:none;strokewidth:%.2f;stroke:%s;stroke-opacity:%.2f", sw, color, setop(opacity)))
}

// dorect draws a rectangle
func dorect(doc *svg.SVG, x, y, w, h float64, color string, opacity float64) {
	doc.Rect(x, y, w, h, fmt.Sprintf("fill:%s;fill-opacity:%.2f", color, setop(opacity)))
}

// doellipse draws a rectangle
func doellipse(doc *svg.SVG, x, y, w, h float64, color string, opacity float64) {
	doc.Ellipse(x, y, w, h, fmt.Sprintf("fill:%s;fill-opacity:%.2f", color, setop(opacity)))
}

// dotext places text elements on the canvas according to type
func dotext(doc *svg.SVG, cw, x, y, fs, wp float64, tdata, font, color string, opacity float64, align, ttype string) {
	var tw float64
	const emsperpixel = 10
	ls := fs + ((fs * 4) / 10)
	td := strings.Split(tdata, "\n")
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * ls
		tw = cw - x - 20
		dorect(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)", opacity)
	}
	if ttype == "block" {
		if wp == 0 {
			tw = 50
		} else {
			tw = (cw * (wp / 100.0)) / emsperpixel
		}
		textwrap(doc, x, y, tw, fs, ls, tdata, font, color, opacity)
	} else {
		for _, t := range td {
			showtext(doc, x, y, t, fs, font, color, align)
			y += ls
		}
	}
}

// textalign returns the SVG text alignment operator
func textalign(s string) string {
	switch s {
	case "center", "middle", "mid", "c":
		return "middle"
	case "left", "start", "l":
		return "start"
	case "right", "end", "e":
		return "end"
	}
	return "start"
}

// showtext places fully attributed text at the specified location
func showtext(doc *svg.SVG, x, y float64, s string, fs float64, font, color, align string) {
	doc.Text(x, y, s, `xml:space="preserve"`, fmt.Sprintf("fill:%s;font-size:%.2fpx;font-family:%s;text-anchor:%s", color, fs, fontlookup(font), textalign(align)))
}

// dolists places lists on the canvas
func dolist(doc *svg.SVG, x, y, fs float64, tlist []deck.ListItem, font, color string, opacity float64, ltype string) {
	if font == "" {
		font = "sans"
	}
	doc.Gstyle(fmt.Sprintf("fill-opacity:%.2f;fill:%s;font-family:%s;font-size:%.2fpx", setop(opacity), color, fontlookup(font), fs))
	if ltype == "bullet" {
		x += fs
	}
	ls := fs * 2
	var t string
	for i, tl := range tlist {
		if ltype == "number" {
			t = fmt.Sprintf("%d. ", i+1) + tl.ListText
		} else {
			t = tl.ListText
		}
		if ltype == "bullet" {
			bullet(doc, x, y, fs, color)
		}
		lifmt := ""
		if len(tl.Color) > 0 {
			lifmt += "fill:" + tl.Color
		}
		if len(tl.Font) > 0 {
			lifmt += ";font-family:" + tl.Font
		}
		if len(lifmt) > 0 {
			doc.Text(x, y, t, `xml:space="preserve"`, lifmt)
		} else {
			doc.Text(x, y, t, `xml:space="preserve"`)
		}
		y += ls
	}
	doc.Gend()
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *svg.SVG, x, y, w, fs float64, leading float64, s, font, color string, opacity float64) {
	doc.Gstyle(fmt.Sprintf("fill-opacity:%.2f;fill:%s;font-family:%s;font-size:%.2fpx", setop(opacity), color, fontlookup(font), fs))
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	var line string
	for _, s := range words {
		line += s + " "
		if float64(len(line)) > w {
			doc.Text(xp, yp, line)
			yp += leading
			line = ""
		}
	}
	if len(line) > 0 {
		doc.Text(xp, yp, line)
	}
	doc.Gend()
}

// doslides reads the deck file, making the SVG version
func doslides(outname, filename, title string, width, height float64, gp float64) {
	var d deck.Deck
	var err error

	d, err = deck.Read(filename, int(width), int(height))
	if err != nil {
		fmt.Fprintf(os.Stderr, "svgdeck: %v\n", err)
		return
	}
	d.Canvas.Width = int(width)
	d.Canvas.Height = int(height)

	if outname == "" {
		doc := svg.New(os.Stdout)
		for i := 0; i < len(d.Slide); i++ {
			doc.Start(width, height)
			svgslide(doc, d, i, gp, outname, title)
			doc.End()
		}
		return
	}

	for i := 0; i < len(d.Slide); i++ {
		out, err := os.Create(fmt.Sprintf(namefmt, outname, i+1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "svgdeck: %v\n", err)
			continue
		}
		doc := svg.New(out)
		doc.Start(width, height)
		svgslide(doc, d, i, gp, outname, title)
		doc.End()
		out.Close()
	}
}

// svgslide makes one slide per SVG page
func svgslide(doc *svg.SVG, d deck.Deck, n int, gp float64, outname, title string) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	var x, y, fs float64

	cw := float64(d.Canvas.Width)
	ch := float64(d.Canvas.Height)
	slide := d.Slide[n]

	// insert navigation links:
	// the full slide links to the next one in sequence,
	// the last slide links to the first
	if len(outname) > 0 {
		var link int
		if n < len(d.Slide)-1 {
			link = n + 1
		} else {
			link = 0
		}
		doc.Link(fmt.Sprintf(namefmt, outname, link), fmt.Sprintf("Link to slide %03d", link))
	}
	// insert title, if specified
	if len(title) > 0 {
		doc.Title(fmt.Sprintf("%s: Slide %d", title, n))
	}
	// set background, if specified
	if len(slide.Bg) > 0 {
		background(doc, cw, ch, slide.Bg)
	}
	// set gradient background, if specified
	if len(slide.Gradcolor1) > 0 && len(slide.Gradcolor2) > 0 {
		oc := []svg.Offcolor{
			{0, slide.Gradcolor1, 1.0},
			{100, slide.Gradcolor2, 1.0},
		}
		doc.Def()
		doc.LinearGradient("slidegrad", 0, 0, 0, 100, oc)
		doc.DefEnd()
		doc.Rect(0, 0, cw, ch, "fill:url(#slidegrad)")
	}
	// set the default foreground
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	// for every image on the slide...
	for _, im := range slide.Image {
		x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
		iw, ih := float64(im.Width), float64(im.Height)

		if im.Scale > 0 {
			iw *= (im.Scale / 100)
			ih *= (im.Scale / 100)
		}
		// scale the image to fit the canvas width
		if im.Autoscale == "on" && iw < cw {
			ih = (cw / iw) * ih
			iw = cw
		}

		midx := iw / 2
		midy := ih / 2
		doc.Image(x-midx, y-midy, int(iw), int(ih), im.Name)
		if len(im.Caption) > 0 {
			capsize := deck.Pwidth(im.Sp, float64(cw), float64(pct(2.0, cw)))
			if im.Font == "" {
				im.Font = "sans"
			}
			if im.Color == "" {
				im.Color = slide.Fg
			}
			if im.Align == "" {
				im.Align = "center"
			}
			showtext(doc, x, y+midy+(capsize*2), im.Caption, capsize, im.Font, im.Color, im.Align)
		}
	}
	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	// rect
	for _, rect := range slide.Rect {
		x, y, _ := dimen(cw, ch, rect.Xp, rect.Yp, 0)
		var w, h float64
		w = pct(rect.Wp, cw)
		if rect.Hr == 0 {
			h = pct(rect.Hp, ch)
		} else {
			h = pct(rect.Hr, w)
		}
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		dorect(doc, x-(w/2), y-(h/2), w, h, rect.Color, rect.Opacity)
	}
	// ellipse
	for _, ellipse := range slide.Ellipse {
		x, y, _ := dimen(cw, ch, ellipse.Xp, ellipse.Yp, 0)
		var w, h float64
		w = pct(ellipse.Wp, cw)
		if ellipse.Hr == 0 {
			h = pct(ellipse.Hp, ch)
		} else {
			h = pct(ellipse.Hr, w)
		}
		if ellipse.Color == "" {
			ellipse.Color = defaultColor
		}
		doellipse(doc, x, y, w/2, h/2, ellipse.Color, ellipse.Opacity)
	}
	// curve
	for _, curve := range slide.Curve {
		if curve.Color == "" {
			curve.Color = defaultColor
		}
		x1, y1, sw := dimen(cw, ch, curve.Xp1, curve.Yp1, curve.Sp)
		x2, y2, _ := dimen(cw, ch, curve.Xp2, curve.Yp2, 0)
		x3, y3, _ := dimen(cw, ch, curve.Xp3, curve.Yp3, 0)
		if sw == 0 {
			sw = 2.0
		}
		docurve(doc, x1, y1, x2, y2, x3, y3, sw, curve.Color, curve.Opacity)
	}
	// arc
	for _, arc := range slide.Arc {
		if arc.Color == "" {
			arc.Color = defaultColor
		}
		x, y, sw := dimen(cw, ch, arc.Xp, arc.Yp, arc.Sp)
		w := pct(arc.Wp, cw)
		h := pct(arc.Hp, cw)
		if sw == 0 {
			sw = 2.0
		}
		doarc(doc, x, y, w/2, h/2, arc.A1, arc.A2, sw, arc.Color, arc.Opacity)
	}
	// line
	for _, line := range slide.Line {
		if line.Color == "" {
			line.Color = defaultColor
		}
		x1, y1, sw := dimen(cw, ch, line.Xp1, line.Yp1, line.Sp)
		x2, y2, _ := dimen(cw, ch, line.Xp2, line.Yp2, 0)
		if sw == 0 {
			sw = 2.0
		}
		doline(doc, x1, y1, x2, y2, sw, line.Color, line.Opacity)
	}
	// for every text element...
	var tdata string
	for _, t := range slide.Text {
		if t.Color == "" {
			t.Color = slide.Fg
		}
		if t.Font == "" {
			t.Font = "sans"
		}
		if t.File != "" {
			tdata = includefile(t.File)
		} else {
			tdata = t.Tdata
		}
		x, y, fs = dimen(cw, ch, t.Xp, t.Yp, t.Sp)
		dotext(doc, cw, x, y, fs, t.Wp, tdata, t.Font, t.Color, t.Opacity, t.Align, t.Type)
	}
	// for every list element...
	for _, l := range slide.List {
		if l.Color == "" {
			l.Color = slide.Fg
		}
		x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
		dolist(doc, x, y, fs, l.Li, l.Font, l.Color, l.Opacity, l.Type)
	}
	// add a grid, if specified
	if gp > 0 {
		grid(doc, cw, ch, slide.Fg, gp)
	}
	// complete the link
	if len(outname) > 0 {
		doc.LinkEnd()
	}
}

// dodeck turns deck input files into SVG
// if the sflag is set, all output goes to the standard output file,
// otherwise, SVG is written the destination directory, to filenames based on the input name.
func dodeck(files []string, sflag bool, pw, ph float64, outdir, title string, gp float64) {
	if sflag { // combined output to standard output
		for _, filename := range files {
			doslides("", filename, title, pw, ph, gp)
		}
	} else { // output to individual files
		for _, filename := range files {
			base := strings.Split(filepath.Base(filename), ".xml")
			outname := filepath.Join(outdir, base[0])
			doslides(outname, filename, title, pw, ph, gp)
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
	var (
		sansfont = flag.String("sans", "Helvetica", "sans font")
		serifont = flag.String("serif", "Times-Roman", "serif font")
		monofont = flag.String("mono", "Courier", "mono font")
		outdir   = flag.String("outdir", ".", "output directory")
		stdout   = flag.Bool("stdout", false, "output to standard output")
		pagesize = flag.String("pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
		title    = flag.String("title", "", "document title")
		gridpct  = flag.Float64("grid", 0, "place percentage grid on each slide")
	)
	flag.Parse()
	var pw, ph float64
	nd, err := fmt.Sscanf(*pagesize, "%g,%g", &pw, &ph)
	if nd != 2 || err != nil {
		pw, ph = 0, 0
	}
	if pw == 0 && ph == 0 {
		p, ok := pagemap[*pagesize]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = (p.width * p.unit)
		ph = (p.height * p.unit)
	}
	fontmap["sans"] = *sansfont
	fontmap["serif"] = *serifont
	fontmap["mono"] = *monofont
	dodeck(flag.Args(), *stdout, pw, ph, *outdir, *title, *gridpct)
}
