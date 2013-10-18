// Package deck makes slide decks
package deck

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// Deck defines the structure of a presentation deck
// The size of the canvas, and series of slides
type Deck struct {
	Canvas canvas  `xml:"canvas"`
	Slide  []slide `xml:"slide"`
}

type canvas struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

type slide struct {
	Bg      string    `xml:"bg,attr"`
	Fg      string    `xml:"fg,attr"`
	List    []List    `xml:"list"`
	Text    []Text    `xml:"text"`
	Image   []Image   `xml:"image"`
	Ellipse []Ellipse `xml:"ellipse"`
	Line    []Line    `xml:"line"`
	Rect    []Rect    `xml:"rect"`
	Curve   []Curve   `xml:"curve"`
	Arc     []Arc     `xml:"arc"`
}

// CommonAttr are the common attributes for text and list
type CommonAttr struct {
	Xp      float64 `xml:"xp,attr"`
	Yp      float64 `xml:"yp,attr"`
	Sp      float64 `xml:"sp,attr"`
	Type    string  `xml:"type,attr"`
	Align   string  `xml:"align,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
	Font    string  `xml:"font,attr"`
}

// Dimension describes a graphics object with width and height
type Dimension struct {
	CommonAttr
	Wp float64 `xml:"wp,attr"`
	Hp float64 `xml:"hp,attr"`
	Hr float64 `xml:"hr,attr"`
	Hw float64 `xml:"hw,attr"`
}

// ListItem describes a list item
type ListItem struct {
	Color    string  `xml:"color,attr"`
	Opacity  float64 `xml:"opacity,attr"`
	Font     string  `xml:"font,attr"`
	ListText string  `xml:",chardata"`
}

// List describes the list element
type List struct {
	CommonAttr
	Li []ListItem `xml:"li"`
}

// Text describes the text element
type Text struct {
	CommonAttr
	Wp    float64 `xml:"wp,attr"`
	Tdata string  `xml:",chardata"`
}

// Image describes an image
type Image struct {
	CommonAttr
	Width   int    `xml:"width,attr"`
	Height  int    `xml:"height,attr"`
	Name    string `xml:"name,attr"`
	Caption string `xml:"caption,attr"`
}

// Ellipse describes a rectangle with x,y,w,h
type Ellipse struct {
	Dimension
}

// Rect describes a rectangle with x,y,w,h
type Rect struct {
	Dimension
}

// Line defines a straight line
type Line struct {
	Xp1     float64 `xml:"xp1,attr"`
	Yp1     float64 `xml:"yp1,attr"`
	Xp2     float64 `xml:"xp2,attr"`
	Yp2     float64 `xml:"yp2,attr"`
	Sp      float64 `xml:"sp,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Curve defines a quadratic Bezier curve
type Curve struct {
	Xp1     float64 `xml:"xp1,attr"`
	Yp1     float64 `xml:"yp1,attr"`
	Xp2     float64 `xml:"xp2,attr"`
	Yp2     float64 `xml:"yp2,attr"`
	Xp3     float64 `xml:"xp3,attr"`
	Yp3     float64 `xml:"yp3,attr"`
	Sp      float64 `xml:"sp,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Arc defines an elliptical arc
type Arc struct {
	Dimension
	A1      float64 `xml:"a1,attr"`
	A2      float64 `xml:"a2,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Read reads the deck description file
func Read(filename string, w, h int) (Deck, error) {
	var d Deck
	r, err := os.Open(filename)
	if err != nil {
		return d, err
	}
	err = xml.NewDecoder(r).Decode(&d)
	if d.Canvas.Width == 0 {
		d.Canvas.Width = w
	}
	if d.Canvas.Height == 0 {
		d.Canvas.Height = h
	}
	r.Close()
	return d, err
}

// Dimen computes the coordinates and size of an object
func Dimen(c canvas, xp, yp, sp float64) (x, y, s float64) {
	x = (xp / 100) * float64(c.Width)
	y = (yp / 100) * float64(c.Height)
	s = (sp / 100) * float64(c.Width)
	return
}

// Pwidth computes the percent width based on canvas size
func Pwidth(wp, cw, defval float64) float64 {
	if wp == 0 {
		return defval
	}
	return (wp / 100) * cw
}

// Search searches the deck for the specified text, returning the slide number if found
func Search(d Deck, s string) int {
	// for every slide...
	for i, slide := range d.Slide {
		// search lists
		for _, l := range slide.List {
			for _, ll := range l.Li {
				if strings.Contains(ll.ListText, s) {
					return i
				}
			}
		}
		// search text
		for _, t := range slide.Text {
			if strings.Contains(t.Tdata, s) {
				return i
			}
		}
	}
	return -1
}

// Dump shows the decoded description
func Dump(d Deck) {
	fmt.Printf("Canvas = %v\n", d.Canvas)
	for i, s := range d.Slide {
		fmt.Printf("Slide [%d] = %#v %#v\n", i, s.Bg, s.Fg)
		for j, l := range s.List {
			fmt.Printf("\tList [%d] = %#v\n", j, l)
		}
		for k, t := range s.Text {
			fmt.Printf("\tText [%d] = %#v\n", k, t)
		}
		for m, im := range s.Image {
			fmt.Printf("\tImage [%d] = %#v\n", m, im)
		}
		for l, line := range s.Line {
			fmt.Printf("\tLine [%d] = %#v\n", l, line)
		}
		for r, rect := range s.Rect {
			fmt.Printf("\tRect [%d] = %#v\n", r, rect)
		}
		for a, arc := range s.Arc {
			fmt.Printf("\tArc [%d] = %#v\n", a, arc)
		}
		for c, curve := range s.Curve {
			fmt.Printf("\tCurve [%d] = %#v\n", c, curve)
		}
		for e, ellipse := range s.Ellipse {
			fmt.Printf("\tEllipse [%d] = %#v\n", e, ellipse)
		}
	}
}
