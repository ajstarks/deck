// gcdeck: render deck markup using the gio canvas
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/deck"
	gc "github.com/ajstarks/giocanvas"
)

const (
	mm2pt       = 2.83464 // mm to pt conversion
	linespacing = 1.8
	listspacing = 1.8
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

var gridstate bool

// pagedim converts a named pagesize to width, height
func pagedim(s string) (float32, float32) {
	var pw, ph float64
	nd, err := fmt.Sscanf(s, "%g,%g", &pw, &ph)
	if nd != 2 || err != nil {
		pw, ph = 0.0, 0.0
	}
	if pw == 0 && ph == 0 {
		p, ok := pagemap[s]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = p.width * p.unit
		ph = p.height * p.unit
	}
	return float32(pw), float32(ph)
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

// pct converts percentages to canvas measures
func pct(p, m float64) float64 {
	return (p / 100.0) * m
}

// fontlookup maps font aliases to implementation font names
func fontlookup(s string) string {
	font, ok := fontmap[s]
	if ok {
		return font
	}
	return "sans"
}

// setop sets the opacity as a truncated fraction of 255
func setop(v float64) uint8 {
	if v > 0.0 {
		return uint8(255.0 * (v / 100.0))
	}
	return 255
}

// gradient sets the background color gradient
func gradient(doc *gc.Canvas, w, h float64, gc1, gc2 string, gp float64) {
}

// doline draws a line
func doline(doc *gc.Canvas, xp1, yp1, xp2, yp2, sw float64, color string, opacity float64) {
	c := gc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Line(float32(xp1), float32(yp1), float32(xp2), float32(yp2), float32(sw), c)
}

// doarc draws an arc
func doarc(doc *gc.Canvas, x, y, w, h, a1, a2, sw float64, color string, opacity float64) {
}

// docurve draws a bezier curve
func docurve(doc *gc.Canvas, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string, opacity float64) {
	c := gc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Curve(float32(xp1), float32(yp1), float32(xp2), float32(yp2), float32(xp3), float32(yp3), c)
}

// dorect draws a rectangle
func dorect(doc *gc.Canvas, x, y, w, h float64, color string, opacity float64) {
	c := gc.ColorLookup(color)
	c.A = setop(opacity)
	doc.CenterRect(float32(x), float32(y), float32(w), float32(h), c)
}

// doellipse draws an ellipse
func doellipse(doc *gc.Canvas, x, y, w, h float64, color string, opacity float64) {
	c := gc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Ellipse(float32(x), float32(y), float32(w), float32(h), c)
}

// dopoly draws a polygon
func dopoly(doc *gc.Canvas, xc, yc string, cw, ch float64, color string, opacity float64) {
	xs := strings.Split(xc, " ")
	ys := strings.Split(yc, " ")
	if len(xs) != len(ys) {
		return
	}
	if len(xs) < 3 || len(ys) < 3 {
		return
	}
	px := make([]float32, len(xs))
	py := make([]float32, len(xs))
	for i := 0; i < len(xs); i++ {
		x, err := strconv.ParseFloat(xs[i], 32)
		if err != nil {
			px[i] = 0
		} else {
			px[i] = float32(pct(x, cw))
		}
		y, err := strconv.ParseFloat(ys[i], 32)
		if err != nil {
			py[i] = 0
		} else {
			py[i] = float32(pct(100-y, ch))
		}
	}
	c := gc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Polygon(px, py, c)
}

// textwidth returns the width of text
func textwidth(c *gc.Canvas, s string, size float64) float64 {
	return float64(len(s)) * (size * 0.65)
}

// dotext places text elements on the canvas according to type
func dotext(doc *gc.Canvas, x, y, fs, wp, rotation, spacing float64, tdata, font, align, ttype, color string, opacity float64) {
	td := strings.Split(tdata, "\n")
	c := gc.ColorLookup(color)
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * spacing * fs
		bx := (x + (wp / 2))
		by := (y - (ch / 2)) + (spacing * fs)
		dorect(doc, bx, by, wp+fs, ch+fs, "rgb(240,240,240)", 100)
	}
	if ttype == "block" {
		textwrap(doc, x, y, wp, fs, fs*spacing, tdata, c, font)
	} else {
		ls := spacing * fs
		for _, t := range td {
			showtext(doc, x, y, t, fs, c, font, align)
			y -= ls
		}
	}
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// loadfont loads a font at the specified size
func loadfont(doc *gc.Canvas, s string, size float64) {
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *gc.Canvas, x, y, w, fs, leading float64, s string, color color.RGBA, font string) int {
	var factor float64 = 0.03
	if font == "mono" {
		factor = 1.0
	}
	nbreak := 0
	loadfont(doc, font, fs)
	wordspacing := textwidth(doc, "M", fs)
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw := textwidth(doc, s, fs)
		doc.Text(float32(xp), float32(yp), float32(fs), s, color)
		xp += tw + (wordspacing * factor)
		if xp > edge {
			xp = x
			yp -= leading
			nbreak++
		}
	}
	return nbreak
}

// showtext places fully attributed text at the specified location
func showtext(doc *gc.Canvas, x, y float64, s string, fs float64, color color.RGBA, font, align string) {
	loadfont(doc, font, fs)
	tx := float32(x)
	ty := float32(y)
	tfs := float32(fs)
	switch align {
	case "center", "middle", "mid", "c":
		doc.TextMid(tx, ty, tfs, s, color)
	case "right", "end", "e":
		doc.TextEnd(tx, ty, tfs, s, color)
	default:
		doc.Text(tx, ty, tfs, s, color)
	}
}

// dolist places lists on the canvas
func dolist(doc *gc.Canvas, cw, x, y, fs, lwidth, rotation, spacing float64, list []deck.ListItem, font, ltype, align, color string, opacity float64) {
	if font == "" {
		font = "sans"
	}
	c := gc.ColorLookup(color)
	ls := listspacing * fs
	for i, tl := range list {
		loadfont(doc, font, fs)
		if len(tl.Color) > 0 {
			c = gc.ColorLookup(tl.Color)
		}
		switch ltype {
		case "number":
			showtext(doc, x, y, fmt.Sprintf("%d. ", i+1)+tl.ListText, fs, c, font, align)
		case "bullet":
			doc.Circle(float32(x), float32(y+fs/3), float32(fs/4), c)
			showtext(doc, x+fs, y, tl.ListText, fs, c, font, align)
		case "center":
			showtext(doc, x, y, tl.ListText, fs, c, font, align)
		default:
			showtext(doc, x, y, tl.ListText, fs, c, font, align)
		}
		y -= ls
	}
}

// showslide shows a slide
func showslide(doc *gc.Canvas, d *deck.Deck, n int) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	cw := float64(d.Canvas.Width)
	ch := float64(d.Canvas.Height)
	slide := d.Slide[n]
	// set default background
	if slide.Bg == "" {
		slide.Bg = "white"
	}
	doc.Background(gc.ColorLookup(slide.Bg))

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
		doc.Image(im.Name, float32(im.Xp), float32(im.Yp), im.Width, im.Height, float32(im.Scale))
		if len(im.Caption) > 0 {
			capsize := 1.5
			if im.Font == "" {
				im.Font = "sans"
			}
			if im.Color == "" {
				im.Color = slide.Fg
			}
			if im.Align == "" {
				im.Align = "center"
			}
			var cx, cy float64
			iw := float64(im.Width) * (im.Scale / 100)
			ih := float64(im.Height) * (im.Scale / 100)
			cimx := im.Xp
			switch im.Align {
			case "center", "c", "mid":
				cx = im.Xp
			case "end", "e", "right":
				cx = cimx + pct((iw/2), cw)
			default:
				cx = cimx - pct((iw/2), cw)
			}
			cy = im.Yp - (ih/2)/ch*100 - (capsize * 2)
			showtext(doc, cx, cy, im.Caption, capsize, gc.ColorLookup(im.Color), im.Font, im.Align)
		}
	}
	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	// rect
	for _, rect := range slide.Rect {
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		if rect.Hr == 100 {
			c := gc.ColorLookup(rect.Color)
			c.A = setop(rect.Opacity)
			doc.Rect(float32(rect.Xp), float32(rect.Yp), float32(rect.Wp), float32((rect.Wp)*(cw/ch)), c)
		} else {
			dorect(doc, rect.Xp, rect.Yp, rect.Wp, rect.Hp, rect.Color, rect.Opacity)
		}
	}
	// ellipse
	for _, ellipse := range slide.Ellipse {
		if ellipse.Color == "" {
			ellipse.Color = defaultColor
		}
		if ellipse.Hr == 100 {
			c := gc.ColorLookup(ellipse.Color)
			c.A = setop(ellipse.Opacity)
			doc.Circle(float32(ellipse.Xp), float32(ellipse.Yp), float32(ellipse.Wp)/2, c)
		} else {
			doellipse(doc, ellipse.Xp, ellipse.Yp, ellipse.Wp, ellipse.Hp, ellipse.Color, ellipse.Opacity)
		}
	}
	// curve
	for _, curve := range slide.Curve {
		if curve.Color == "" {
			curve.Color = defaultColor
		}
		if curve.Sp == 0 {
			curve.Sp = 2.0
		}
		docurve(doc, curve.Xp1, curve.Yp1, curve.Xp2, curve.Yp2, curve.Xp3, curve.Yp3, curve.Sp, curve.Color, curve.Opacity)
	}
	// arc
	for _, arc := range slide.Arc {
		if arc.Color == "" {
			arc.Color = defaultColor
		}
		w := arc.Wp
		h := arc.Hp
		if arc.Sp == 0 {
			arc.Sp = 2.0
		}
		doarc(doc, arc.Xp, arc.Yp, w/2, h/2, arc.A1, arc.A2, arc.Sp, arc.Color, arc.Opacity)
	}
	// line
	for _, line := range slide.Line {
		if line.Color == "" {
			line.Color = defaultColor
		}
		if line.Sp == 0 {
			line.Sp = 0.2
		}
		doline(doc, line.Xp1, line.Yp1, line.Xp2, line.Yp2, line.Sp, line.Color, line.Opacity)
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
		if t.File != "" {
			tdata = includefile(t.File)
		} else {
			tdata = t.Tdata
		}
		if t.Lp == 0 {
			t.Lp = linespacing
		}
		dotext(doc, t.Xp, t.Yp, t.Sp, t.Wp, t.Rotation, t.Lp*1.2, tdata, t.Font, t.Align, t.Type, t.Color, t.Opacity)
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
		dolist(doc, cw, l.Xp, l.Yp, l.Sp, l.Wp, l.Rotation, l.Lp, l.Li, l.Font, l.Type, l.Align, l.Color, l.Opacity)
	}

}

// hup processes the hangup (SIGHUP) signal
// re-read the input, adjust the number of slides, show the first slide
func hup(sigch chan os.Signal, filename string, c *gc.Canvas, d *deck.Deck, w, h float32, n *int) {
	for range sigch {
		newdeck, err := readDeck(filename, w, h)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		nd := len(newdeck.Slide)
		d.Slide = make([]deck.Slide, nd)
		copy(d.Slide, newdeck.Slide)
		*n = nd - 1
		showslide(c, d, 0)
	}
}

// ReadDeck reads the deck file, rendering to the canvas
func readDeck(filename string, w, h float32) (deck.Deck, error) {
	d, err := deck.Read(filename, int(w), int(h))
	d.Canvas.Width = int(w)
	d.Canvas.Height = int(h)
	return d, err
}

// back shows the previous slide
func back(c *gc.Canvas, d *deck.Deck, n *int, limit int) {
	*n--
	if *n < 0 {
		*n = limit
	}
	showslide(c, d, *n)
}

// forward shows the next slide
func forward(c *gc.Canvas, d *deck.Deck, n *int, limit int) {
	*n++
	if *n > limit {
		*n = 0
	}
	showslide(c, d, *n)
}

// reload reloads the content and shows the first slide
func reload(filename string, c *gc.Canvas, w, h, n int) (deck.Deck, int) {
	d, err := readDeck(filename, float32(w), float32(h))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return d, 0
	}
	showslide(c, &d, n)
	return d, len(d.Slide) - 1
}

// gridtoggle toggles a grid overlay
func gridtoggle(c *gc.Canvas, size float32, d *deck.Deck, slidenumber int) {
	if gridstate {
		gcolor := gc.ColorLookup(d.Slide[slidenumber].Fg)
		gcolor.A = 100
		c.Grid(0, 0, 100, 100, 0.1, size, gcolor)
	} else {
		showslide(c, d, slidenumber)
	}
	gridstate = !gridstate
}

func main() {
	var (
		title    = flag.String("title", "", "slide title")
		pagesize = flag.String("pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
		initpage = flag.Int("page", 1, "initial page")
	)
	flag.Parse()

	// get the filename
	var filename string
	if len(flag.Args()) < 1 {
		filename = "-"
		*title = "Standard Input"
	} else {
		filename = flag.Args()[0]
	}
	if *title == "" {
		*title = filename
	}
	go slidedeck(*title, *initpage, filename, *pagesize)
	app.Main()
}

func slidedeck(s string, initpage int, filename, pagesize string) {
	gofont.Register()

	width, height := pagedim(pagesize)
	deck, err := readDeck(filename, width, height)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// set initial values
	nslides := len(deck.Slide) - 1
	if initpage > nslides+1 || initpage < 1 {
		initpage = 1
	}
	slidenumber := initpage - 1

	gridstate = true
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGHUP)
	win := app.NewWindow(app.Title(s), app.Size(unit.Dp(width), unit.Dp(height)))
	for we := range win.Events() {
		if e, ok := we.(system.FrameEvent); ok {
			canvas := gc.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			go hup(sigch, filename, canvas, &deck, width, height, &nslides)
			showslide(canvas, &deck, slidenumber)
			e.Frame(canvas.Context.Ops)
		}
	}
}
