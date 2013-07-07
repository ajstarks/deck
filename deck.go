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
	Bg    string  `xml:"bg,attr"`
	Fg    string  `xml:"fg,attr"`
	List  []list  `xml:"list"`
	Text  []text  `xml:"text"`
	Image []image `xml:"image"`
}

// CommonAttr are the common attributes for text and list
type CommonAttr struct {
	Xp    float64 `xml:"xp,attr"`
	Yp    float64 `xml:"yp,attr"`
	Sp    float64 `xml:"sp,attr"`
	Type  string  `xml:"type,attr"`
	Align string  `xml:"align,attr"`
	Color string  `xml:"color,attr"`
	Font  string  `xml:"font,attr"`
}

type list struct {
	CommonAttr
	Li []string `xml:"li"`
}

type text struct {
	CommonAttr
	Wp    float64 `xml:"wp,attr"`
	Tdata string  `xml:",chardata"`
}

type image struct {
	Xp     float64 `xml:"xp,attr"`
	Yp     float64 `xml:"yp,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Name   string  `xml:"name,attr"`
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
func Dimen(c canvas, xp, yp, sp float64) (x, y float64, s int) {
	x = (xp / 100) * float64(c.Width)
	y = (yp / 100) * float64(c.Height)
	s = int((sp / 100) * float64(c.Width))
	return
}


// Pwidth computes the percent width based on canvas size
func Pwidth(wp, cw, defval float64) float64 {
	if wp == 0 {
		return defval
	}
	return (wp/100)  * cw
}

// Search searches the deck for the specified text, returning the slide number if found
func Search(d Deck, s string) int {
	// for every slide...
	for i, slide := range d.Slide {
		// search lists
		for _, l := range slide.List {
			for _, ll := range l.Li {
				if strings.Contains(ll, s) {
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
	}
}

