package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ajstarks/deck/generate"
)

var codemap = strings.NewReplacer(
	"\t", "    ",
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;")

func includefile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "specify a code file")
		return
	}
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	slide := 0
	for _, codefile := range os.Args[1:] {
		i := strings.LastIndex(codefile, ".go")
		if i < 0 {
			fmt.Fprintf(os.Stderr, "cannot get the basename for %s\n", codefile)
			continue
		}
		imagefile := codefile[:i] + ".png"
		slide++
		deck.StartSlide()
		deck.Image(75, 68, 500, 500, imagefile)
		deck.Text(2.5, 95, includefile(codefile), "mono", 1.2, "black")
		deck.TextEnd(90, 2.5, codefile, "sans", 2, "black")
		deck.TextEnd(95, 2.5, fmt.Sprintf("[%d]", slide), "sans", 2, "gray")
		deck.EndSlide()
	}
	deck.EndDeck()

}
