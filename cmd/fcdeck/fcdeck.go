// fcdeck: render deck markup using the fyne canvas
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/ajstarks/deck"
	"github.com/ajstarks/fc"
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
func pagedim(s string) (int, int) {
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
	return int(pw), int(ph)
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

// grid makes a percentage scale
func grid(doc *fc.Canvas, w, h float64, color string, percent float64) {
	c := fc.ColorLookup(color)
	c.A = 100
	fs := pct(1.5, w)
	for x, pl := 0.0, 0.0; x <= w; x += percent {
		doc.Line(x, 0, x, h, 0.1, c)
		if pl > 0 {
			doc.CText(x, fs, fs, fmt.Sprintf("%.0f", pl), c)
		}
		pl += percent
	}
	for y, pl := 0.0, 0.0; y <= h; y += percent {
		doc.Line(0, y, w, y, 0.1, c)
		if pl < 100 {
			doc.Text(fs, y+(fs/3), fs, fmt.Sprintf("%.0f", pl), c)
		}
		pl += percent
	}
	doc.Container.Refresh()
}

// setop sets the opacity as a truncated fraction of 255
func setop(v float64) uint8 {
	if v > 0.0 {
		return uint8(255.0 * (v / 100.0))
	}
	return 255
}

// background places a colored rectangle
func background(doc *fc.Canvas, w, h float64, color string) {
	doc.Rect(w/2, h/2, w, h, fc.ColorLookup(color))
}

// gradient sets the background color gradient
func gradient(doc *fc.Canvas, w, h float64, gc1, gc2 string, gp float64) {
}

// doline draws a line
func doline(doc *fc.Canvas, xp1, yp1, xp2, yp2, sw float64, color string, opacity float64) {
	c := fc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Line(xp1, yp1, xp2, yp2, sw, c)
}

// doarc draws an arc
func doarc(doc *fc.Canvas, x, y, w, h, a1, a2, sw float64, color string, opacity float64) {
}

// docurve draws a bezier curve
func docurve(doc *fc.Canvas, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string, opacity float64) {
}

// dorect draws a rectangle
func dorect(doc *fc.Canvas, x, y, w, h float64, color string, opacity float64) {
	c := fc.ColorLookup(color)
	c.A = setop(opacity)
	doc.Rect(x, y, w, h, c)
}

// doellipse draws an ellipse
func doellipse(doc *fc.Canvas, x, y, w, h float64, color string, opacity float64) {
}

// dopoly draws a polygon
func dopoly(doc *fc.Canvas, xc, yc string, cw, ch float64, color string, opacity float64) {
}

// dotext places text elements on the canvas according to type
func dotext(doc *fc.Canvas, x, y, fs, wp, rotation, spacing float64, tdata, font, align, ttype, color string, opacity float64) {
	td := strings.Split(tdata, "\n")
	c := fc.ColorLookup(color)
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
func loadfont(doc *fc.Canvas, s string, size float64) {
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *fc.Canvas, x, y, w, fs, leading float64, s string, color color.RGBA, font string) int {
	var factor = 0.3
	if font == "mono" {
		factor = 1.0
	}
	nbreak := 0
	loadfont(doc, font, fs)
	wordspacing := doc.TextWidth("M", fs)
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw := doc.TextWidth(s, fs)
		doc.Text(xp, yp, fs, s, color)
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
func showtext(doc *fc.Canvas, x, y float64, s string, fs float64, color color.RGBA, font, align string) {
	loadfont(doc, font, fs)
	switch align {
	case "center", "middle", "mid", "c":
		doc.CText(x, y, fs, s, color)
	case "right", "end", "e":
		doc.EText(x, y, fs, s, color)
	default:
		doc.Text(x, y, fs, s, color)
	}
}

// dolist places lists on the canvas
func dolist(doc *fc.Canvas, cw, x, y, fs, lwidth, rotation, spacing float64, list []deck.ListItem, font, ltype, align, color string, opacity float64) {
	if font == "" {
		font = "sans"
	}
	c := fc.ColorLookup(color)
	ls := listspacing * fs
	for i, tl := range list {
		loadfont(doc, font, fs)
		if len(tl.Color) > 0 {
			c = fc.ColorLookup(tl.Color)
		}
		switch ltype {
		case "number":
			showtext(doc, x, y, fmt.Sprintf("%d. ", i+1)+tl.ListText, fs, c, font, align)
		case "bullet":
			doc.Circle(x, y+fs/3, fs/2, c)
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
func showslide(doc *fc.Canvas, d deck.Deck, n int) {
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
	doc.CornerRect(0, 100, cw, ch, fc.ColorLookup(slide.Bg))

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
		doc.Image(im.Xp, im.Yp, iw, ih, im.Name)
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
			var cx float64
			switch im.Align {
			case "center", "c", "mid":
				cx = im.Xp
			case "end", "e", "right":
				cx = im.Xp + pct(float64(iw/2), cw)
			default:
				cx = im.Xp - pct(float64(iw/2), cw)
			}
			cy := im.Yp - ((float64(ih/2) / ch) * 100) - (capsize * 2)
			showtext(doc, cx, cy, im.Caption, capsize, fc.ColorLookup(im.Color), im.Font, im.Align)
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
			c := fc.ColorLookup(rect.Color)
			c.A = setop(rect.Opacity)
			doc.Rect(rect.Xp, rect.Yp, rect.Wp, rect.Wp*(cw/ch), c)
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
			c := fc.ColorLookup(ellipse.Color)
			c.A = setop(ellipse.Opacity)
			doc.Circle(ellipse.Xp, ellipse.Yp, ellipse.Wp, c)
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
	doc.Container.Refresh()

}

// ReadDeck reads the deck file, rendering to the canvas
func readDeck(filename string, w, h int) (deck.Deck, error) {
	d, err := deck.Read(filename, w, h)
	d.Canvas.Width = int(w)
	d.Canvas.Height = int(h)
	return d, err
}

// back shows the previous slide
func back(c *fc.Canvas, d deck.Deck, n *int, limit int) {
	*n--
	if *n < limit {
		*n = 0
	}
	showslide(c, d, *n)
}

// forward shows the next slide
func forward(c *fc.Canvas, d deck.Deck, n *int, limit int) {
	*n++
	if *n > limit {
		*n = 0
	}
	showslide(c, d, *n)
}

// reload reloads the content and shows the first slide
func reload(filename string, c *fc.Canvas, w, h int) {
	d, err := readDeck(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	showslide(c, d, 0)
}

// gridtoggle toggles a grid overlay
func gridtoggle(c *fc.Canvas, size float64, d deck.Deck, slidenumber int) {
	if gridstate {
		grid(c, 100, 100, d.Slide[slidenumber].Fg, size)
	} else {
		showslide(c, d, slidenumber)
	}
	gridstate = !gridstate
}

func main() {
	var (
		title    = flag.String("title", "", "slide title")
		pagesize = flag.String("pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
		sans     = flag.String("sans", "", "default font")
		initpage = flag.Int("page", 1, "initial page")
		gp       = flag.Float64("grid", 5, "grid percent")
	)
	flag.Parse()

	// set the font
	if *sans != "" {
		os.Setenv("FYNE_FONT", fmt.Sprintf("%s/%s.ttf", os.Getenv("DECKFONTS"), *sans))
	}
	// get the filename
	var filename string
	if len(flag.Args()) < 1 {
		filename = "-"
	} else {
		filename = flag.Args()[0]
	}
	if *title == "" {
		*title = filename
	}
	// initialize the canvas, read the data
	width, height := pagedim(*pagesize)
	c := fc.NewCanvas(*title, width, height)
	w := c.Window
	d, err := readDeck(filename, width, height)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// set initial values
	nslides := len(d.Slide) - 1
	if *initpage > nslides+1 || *initpage < 1 {
		*initpage = 1
	}
	canvas := &c
	slidenumber := *initpage - 1
	showslide(canvas, d, slidenumber)

	gridstate = true

	// Define keyboard shortcuts (back, forward, reload, grid, home, end, quit)
	canvas.Window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyRight, fyne.KeyDown, fyne.KeySpace, fyne.KeyEqual, fyne.KeyReturn, fyne.KeyPageDown:
			forward(canvas, d, &slidenumber, nslides)

		case fyne.KeyLeft, fyne.KeyUp, fyne.KeyMinus, fyne.KeyBackspace, fyne.KeyPageUp:
			back(canvas, d, &slidenumber, 1)

		case fyne.KeyR:
			reload(filename, canvas, width, height)

		case fyne.KeyG:
			gridtoggle(canvas, *gp, d, slidenumber)

		case fyne.KeyHome:
			slidenumber = 0
			showslide(canvas, d, slidenumber)

		case fyne.KeyEnd:
			slidenumber = nslides
			showslide(canvas, d, slidenumber)

		case fyne.KeyQ, fyne.KeyEscape:
			os.Exit(0)
		}
	})

	// define the toolbar (back, forward, reload, grid)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() { back(canvas, d, &slidenumber, 1) }),
		widget.NewToolbarAction(theme.NavigateNextIcon(), func() { forward(canvas, d, &slidenumber, nslides) }),
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() { reload(filename, canvas, width, height) }),
		widget.NewToolbarAction(theme.MediaSkipPreviousIcon(), func() { slidenumber = 0; showslide(canvas, d, slidenumber) }),
		widget.NewToolbarAction(theme.MediaSkipNextIcon(), func() { slidenumber = nslides; showslide(canvas, d, slidenumber) }),
		widget.NewToolbarAction(theme.VisibilityIcon(), func() { gridtoggle(canvas, *gp, d, slidenumber) }),
	)
	// add the content
	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar, c.Container))
	w.Resize(fyne.NewSize(width, height+toolbar.Size().Height))

	// run it!
	w.ShowAndRun()
}
