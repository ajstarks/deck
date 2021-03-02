// pngdeck: render deck markup into a series of PNG images
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ajstarks/deck"
	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
)

const (
	mm2pt       = 2.83464 // mm to pt conversion
	linespacing = 1.4
	listspacing = 2.0
	fontfactor  = 1.0
	listwrap    = 95.0
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

// pagerange returns the begin and end using a "-" string
func pagerange(s string) (int, int) {
	p := strings.Split(s, "-")
	if len(p) != 2 {
		return 0, 0
	}
	b, berr := strconv.Atoi(p[0])
	e, err := strconv.Atoi(p[1])
	if berr != nil || err != nil {
		return 0, 0
	}
	if b > e {
		return 0, 0
	}
	return b, e
}

// includefile returns the contents of a file as string
func includefile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

// pct converts percentages to canvas measures
func pct(p, m float64) float64 {
	return (p / 100.0) * m
}

// dimen returns canvas dimensions from percentages
func dimen(w, h, xp, yp, sp float64) (float64, float64, float64) {
	return pct(xp, w), pct(100-yp, h), pct(sp, w) * fontfactor
}

// fontlookup maps font aliases to implementation font names
func fontlookup(s string) string {
	font, ok := fontmap[s]
	if ok {
		return font
	}
	return "sans"
}

// grid makes a percentage scale
func grid(doc *gg.Context, w, h float64, color string, percent float64) {
	pw := w * (percent / 100)
	ph := h * (percent / 100)
	r, g, b := colorlookup(color)
	doc.SetRGB255(r, g, b)
	doc.SetLineWidth(0.25)
	fs := pct(1, w)
	for x, pl := 0.0, 0.0; x <= w; x += pw {
		doc.DrawLine(x, 0, x, h)
		doc.Stroke()
		if pl > 0 {
			showtext(doc, x, h-fs, fmt.Sprintf("%.0f", pl), fs, "sans", "center")
		}
		pl += percent
	}
	for y, pl := 0.0, 0.0; y <= h; y += ph {
		doc.DrawLine(0, y, w, y)
		doc.Stroke()
		if pl < 100 {
			showtext(doc, fs, y+(fs/3), fmt.Sprintf("%.0f", 100-pl), fs, "sans", "center")
		}
		pl += percent
	}
}

// setop sets the opacity as a truncated fraction of 255
func setop(v float64) int {
	if v > 0.0 {
		return int(255.0 * (v / 100.0))
	}
	return 255
}

// bullet draws a bullet
func bullet(doc *gg.Context, x, y, size float64, color string) {
	rs := size / 2
	r, g, b := colorlookup(color)
	doc.SetRGB255(r, g, b)
	doc.DrawCircle(x-size*2, y-rs, rs)
	doc.Fill()
}

// background places a colored rectangle
func background(doc *gg.Context, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetRGB255(r, g, b)
	doc.Clear()
}

// gradient sets the background color gradient
func gradient(doc *gg.Context, w, h float64, gc1, gc2 string, gp float64) {
	r1, g1, b1 := colorlookup(gc1)
	r2, g2, b2 := colorlookup(gc2)
	gp /= 100.0
	grad := gg.NewLinearGradient(0, 0, 100, 100)
	grad.AddColorStop(0, color.RGBA{uint8(r1), uint8(g1), uint8(b1), 1})
	grad.AddColorStop(0, color.RGBA{uint8(r2), uint8(g2), uint8(b2), 1})
	doc.SetFillStyle(grad)
	doc.DrawRectangle(0, 0, w, h)
	doc.Fill()
}

// doline draws a line
func doline(doc *gg.Context, xp1, yp1, xp2, yp2, sw float64, color string, opacity float64) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.SetLineCapButt()
	doc.DrawLine(xp1, yp1, xp2, yp2)
	doc.Stroke()
}

// doarc draws an arc
func doarc(doc *gg.Context, x, y, w, h, a1, a2, sw float64, color string, opacity float64) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.SetLineCapButt()
	doc.DrawEllipticalArc(x, y, w, h, gg.Radians(360-a1), gg.Radians(360-a2))
	doc.Stroke()
}

// docurve draws a bezier curve
func docurve(doc *gg.Context, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string, opacity float64) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetLineCapButt()
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.MoveTo(xp1, yp1)
	doc.QuadraticTo(xp2, yp2, xp3, yp3)
	doc.Stroke()
}

// dorect draws a rectangle
func dorect(doc *gg.Context, x, y, w, h float64, color string, opacity float64) {
	r, g, b := colorlookup(color)
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.DrawRectangle(x, y, w, h)
	doc.Fill()
}

// doellipse draws a rectangle
func doellipse(doc *gg.Context, x, y, w, h float64, color string, opacity float64) {
	r, g, b := colorlookup(color)
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.DrawEllipse(x, y, w, h)
	doc.Fill()
}

// dopoly draws a polygon
func dopoly(doc *gg.Context, xc, yc string, cw, ch float64, color string, opacity float64) {
	xs := strings.Split(xc, " ")
	ys := strings.Split(yc, " ")
	if len(xs) != len(ys) {
		return
	}
	if len(xs) < 3 || len(ys) < 3 {
		return
	}
	doc.NewSubPath()
	for i := 0; i < len(xs); i++ {
		x, err := strconv.ParseFloat(xs[i], 64)
		if err != nil {
			x = 0
		} else {
			x = pct(x, cw)
		}
		y, err := strconv.ParseFloat(ys[i], 64)
		if err != nil {
			y = 0
		} else {
			y = pct(100-y, ch)
		}
		doc.LineTo(x, y)
	}
	doc.ClosePath()
	r, g, b := colorlookup(color)
	doc.SetRGBA255(r, g, b, setop(opacity))
	doc.Fill()
}

// dotext places text elements on the canvas according to type
func dotext(doc *gg.Context, cw, x, y, fs, wp, rotation, spacing float64, tdata, font, align, ttype, color string, opacity float64) {
	var tw float64

	td := strings.Split(tdata, "\n")
	red, green, blue := colorlookup(color)
	if rotation > 0 {
		doc.Push()
		doc.RotateAbout(gg.Radians(360-rotation), x, y)
	}
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * spacing * fs
		tw = deck.Pwidth(wp, cw, cw-x-20)
		dorect(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)", 100)
	}
	doc.SetRGBA255(red, green, blue, setop(opacity))
	if ttype == "block" {
		tw = deck.Pwidth(wp, cw, cw/2)
		textwrap(doc, x, y, tw, fs, fs*spacing, tdata, font)
	} else {
		ls := spacing * fs
		for _, t := range td {
			showtext(doc, x, y, t, fs, font, align)
			y += ls
		}
	}
	if rotation > 0 {
		doc.Pop()
	}
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// loadfont loads a font at the specified size
func loadfont(doc *gg.Context, s string, size float64) {
	f, err := gg.LoadFontFace(fontlookup(s), size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pngdeck %v\n", err)
		return
	}
	doc.SetFontFace(f)
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *gg.Context, x, y, w, fs, leading float64, s, font string) int {
	var factor = 0.3
	if font == "mono" {
		factor = 1.0
	}
	nbreak := 0
	loadfont(doc, font, fs)
	wordspacing, _ := doc.MeasureString("M")
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw, _ := doc.MeasureString(s)
		doc.DrawString(s, xp, yp)
		xp += tw + (wordspacing * factor)
		if xp > edge {
			xp = x
			yp += leading
			nbreak++
		}
	}
	return nbreak
}

// showtext places fully attributed text at the specified location
func showtext(doc *gg.Context, x, y float64, s string, fs float64, font, align string) {
	offset := 0.0
	loadfont(doc, font, fs)
	t := s
	tw, _ := doc.MeasureString(t)
	switch align {
	case "center", "middle", "mid", "c":
		offset = (tw / 2)
	case "right", "end", "e":
		offset = tw
	}
	doc.DrawString(t, x-offset, y)
}

// dolists places lists on the canvas
func dolist(doc *gg.Context, cw, x, y, fs, lwidth, rotation, spacing float64, list []deck.ListItem, font, ltype, align, color string, opacity float64) {
	if font == "" {
		font = "sans"
	}
	red, green, blue := colorlookup(color)

	if ltype == "bullet" {
		x += fs * 1.2
	}
	ls := spacing * fs
	tw := deck.Pwidth(lwidth, cw, cw/2)

	if rotation > 0 {
		doc.Push()
		doc.RotateAbout(gg.Radians(360-rotation), x, y)
	}
	var t string
	for i, tl := range list {
		loadfont(doc, font, fs)
		doc.SetRGB255(red, green, blue)
		if ltype == "number" {
			t = fmt.Sprintf("%d. ", i+1) + tl.ListText
		} else {
			t = tl.ListText
		}
		if ltype == "bullet" {
			bullet(doc, x, y, fs/2, color)
		}
		if len(tl.Color) > 0 {
			tlred, tlgreen, tlblue := colorlookup(tl.Color)
			doc.SetRGB255(tlred, tlgreen, tlblue)
		}
		if len(tl.Font) > 0 {
			loadfont(doc, tl.Font, fs)
		}
		if align == "center" || align == "c" {
			showtext(doc, x, y, t, fs, font, align)
			y += ls
		} else {
			yw := textwrap(doc, x, y, tw, fs, ls, t, font)
			y += ls
			if yw >= 1 {
				y += ls * float64(yw)
			}
		}
	}
	if rotation > 0 {
		doc.Pop()
	}
}

// pngslide makes a slide, one slide per generated PNG
func pngslide(doc *gg.Context, d deck.Deck, n int, gp float64, showslide bool, dest string) {
	if n < 0 || n > len(d.Slide)-1 || !showslide {
		return
	}

	var x, y, fs float64

	cw := float64(d.Canvas.Width)
	ch := float64(d.Canvas.Height)
	slide := d.Slide[n]
	// set default background
	if slide.Bg == "" {
		slide.Bg = "white"
	}
	background(doc, cw, ch, slide.Bg)

	if slide.GradPercent <= 0 || slide.GradPercent > 100 {
		slide.GradPercent = 100
	}
	// set gradient background, if specified. You need both colors
	if len(slide.Gradcolor1) > 0 && len(slide.Gradcolor2) > 0 {
		gradient(doc, cw, ch, slide.Gradcolor1, slide.Gradcolor2, slide.GradPercent)
	}
	// set the default foreground
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	// for every image on the slide...
	for _, im := range slide.Image {
		x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
		iw, ih := im.Width, im.Height
		// scale the image by the specified percentage
		if im.Scale > 0 {
			iw = int(float64(iw) * (im.Scale / 100))
			ih = int(float64(ih) * (im.Scale / 100))
		}
		// scale the image to fit the canvas width
		if im.Autoscale == "on" && iw < d.Canvas.Width {
			ih = int((float64(d.Canvas.Width) / float64(iw)) * float64(ih))
			iw = d.Canvas.Width
		}

		img, err := gg.LoadImage(im.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pngdeck: slide %d (%v)\n", n+1, err)
			return
		}
		bounds := img.Bounds()
		if iw == (bounds.Max.X-bounds.Min.X) && ih == (bounds.Max.Y-bounds.Min.Y) {
			doc.DrawImageAnchored(img, int(x), int(y), 0.5, 0.5)
		} else {
			g := gift.New(gift.Resize(iw, ih, gift.BoxResampling))
			resized := image.NewRGBA(g.Bounds(img.Bounds()))
			g.Draw(resized, img)
			doc.DrawImageAnchored(resized, int(x), int(y), 0.5, 0.5)
		}
		if len(im.Caption) > 0 {
			capsize := deck.Pwidth(im.Sp, cw, pct(2, cw))
			if im.Font == "" {
				im.Font = "sans"
			}
			if im.Color == "" {
				im.Color = slide.Fg
			}
			if im.Align == "" {
				im.Align = "center"
			}
			midx := float64(iw) / 2
			midy := float64(ih) / 2
			switch im.Align {
			case "left", "start":
				x -= midx
			case "right", "end":
				x += midx
			}
			capr, capg, capb := colorlookup(im.Color)
			doc.SetRGB255(capr, capg, capb)
			showtext(doc, x, y+(midy)+(capsize*1.5), im.Caption, capsize, im.Font, im.Align)
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
	// polygon
	for _, poly := range slide.Polygon {
		if poly.Color == "" {
			poly.Color = defaultColor
		}
		dopoly(doc, poly.XC, poly.YC, cw, ch, poly.Color, poly.Opacity)
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
		x, y, fs = dimen(cw, ch, t.Xp, t.Yp, t.Sp)
		if t.File != "" {
			tdata = includefile(t.File)
		} else {
			tdata = t.Tdata
		}
		if t.Lp == 0 {
			t.Lp = linespacing
		}
		dotext(doc, cw, x, y, fs, t.Wp, t.Rotation, t.Lp, tdata, t.Font, t.Align, t.Type, t.Color, t.Opacity)
	}
	// for every list element...
	for _, l := range slide.List {
		if l.Color == "" {
			l.Color = slide.Fg
		}
		if l.Lp == 0 {
			l.Lp = listspacing
		}
		if l.Wp == 0 {
			l.Wp = listwrap
		}
		x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
		dolist(doc, cw, x, y, fs, l.Wp, l.Rotation, l.Lp, l.Li, l.Font, l.Type, l.Align, l.Color, l.Opacity)
	}
	// add a grid, if specified
	if gp > 0 {
		grid(doc, cw, ch, slide.Fg, gp)
	}
	doc.SavePNG(fmt.Sprintf("%s-%05d.png", dest, n+1))
}

// doslides reads the deck file, making a series of PNGs
func doslides(outname, filename string, w, h int, gp float64, begin, end int) {
	var d deck.Deck
	var err error
	d, err = deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pngdeck: %v\n", err)
		return
	}
	d.Canvas.Width = w
	d.Canvas.Height = h

	for i := 0; i < len(d.Slide); i++ {
		pngslide(gg.NewContext(w, h), d, i, gp, (i+1 >= begin && i+1 <= end), outname)
	}
}

// dodeck turns deck input files into PNG files
// PNGs are written to the destination directory, to filenames based on the input name.
func dodeck(files []string, w, h float64, outdir string, gp float64, begin, end int) {
	for _, filename := range files {
		base := strings.Split(filepath.Base(filename), ".xml")
		outname := filepath.Join(outdir, base[0])
		doslides(outname, filename, int(w), int(h), gp, begin, end)
	}
}

// for every file, make a deck
func main() {
	var (
		sansfont   = flag.String("sans", "FiraSans-Regular", "sans font")
		serifont   = flag.String("serif", "Charter-Regular", "serif font")
		monofont   = flag.String("mono", "FiraMono-Regular", "mono font")
		symbolfont = flag.String("symbol", "ZapfDingbats", "symbol font")
		pagesize   = flag.String("pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
		fontdir    = flag.String("fontdir", os.Getenv("DECKFONTS"), "directory for fonts (defaults to DECKFONTS environment variable)")
		outdir     = flag.String("outdir", ".", "output directory")
		gridpct    = flag.Float64("grid", 0, "draw a percentage grid on each slide")
		pr         = flag.String("pages", "1-1000000", "page range (first-last)")
	)
	flag.Parse()

	var pw, ph float64
	nd, err := fmt.Sscanf(*pagesize, "%g,%g", &pw, &ph)
	begin, end := pagerange(*pr)
	if nd != 2 || err != nil {
		pw, ph = 0.0, 0.0
	}
	if pw == 0 && ph == 0 {
		p, ok := pagemap[*pagesize]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = p.width * p.unit
		ph = p.height * p.unit
	}
	fontmap["sans"] = filepath.Join(*fontdir, *sansfont+".ttf")
	fontmap["serif"] = filepath.Join(*fontdir, *serifont+".ttf")
	fontmap["mono"] = filepath.Join(*fontdir, *monofont+".ttf")
	fontmap["symbol"] = filepath.Join(*fontdir, *symbolfont+".ttf")
	dodeck(flag.Args(), pw, ph, *outdir, *gridpct, begin, end)
}
