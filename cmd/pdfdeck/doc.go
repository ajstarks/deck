/*
pdfdeck: render deck files to pdf

```pdfdeck``` reads files in deck markup and makes pdf files.

command options:

	-author string
	  	document author
	-fontdir string
	  	directory for fonts (defaults to DECKFONTS environment variable)
	-grid float
	  	draw a percentage grid on each slide
	-mono string
	  	mono font (default "courier")
	-outdir string
	  	output directory (default ".")
	-pages string
	  	page range (first-last) (default "1-1000000")
	-pagesize string
	  	pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen (default "Letter")
	-sans string
	  	sans font (default "helvetica")
	-serif string
	  	serif font (default "times")
	-stdout
	  	output to standard output
	-symbol string
	  	symbol font (default "zapfdingbats")
	-title string
	  	document title

Fonts

```pdfdeck``` assumes a set of standard fonts (Times, Helvetica, Courier, and Zapf Dingbats) are available.
These fonts and othe TrueType fonts (Noto, Fira, Charter, Go fonts, etc) are available in the [deckfonts](https://github.com/ajstarks/deckfonts) repository.  ```pdfdeck``` also uses the DECKFONTS environment variable to indicate where fonts are stored:

	export DECKFONTS=$HOME/deckfonts
	cd $HOME
	git clone https://github.com/ajstarks/deckfonts
	...
	pdfdeck foo.xml # (use helvetica as the default)
	pdfdeck -sans NotoSans-Regular -serif NotoSerif-Regular -mono NotoMono-Regular foo.xml

To use other TrueType fonts, place them in the DECKFONTS directory and specify using the command options:

	cp FancyFont.ttf $DECKFONTS
	pdfdeck -sans FancyFont foo.xml

 Alternatively you can obtain fonts from the ```gofpdf``` package. To use the standard fonts (assuming the DECKFONTS variable has been set,
 and gofpdf is in GOPATH):


	cp $GOPATH/src/github.com/jung-kurt/gofpdf/font/*.json $DECKFONTS
	cp $GOPATH/src/github.com/jung-kurt/gofpdf/font/*.ttf $DECKFONTS



Example uses

```
pdfdeck foo.xml                                # read deck markup in foo.xml, make foo.pdf
pdfdeck -sans Go-Regular -mono Go-Mono foo.xml # use Go fonts
pdfdeck -fontdir /path/to/fonts foo.xml        # use an alternative font directory
pdfdeck -pages 10-12 foo.xml                   # only render pages 10, 11, and 12
pdfdeck -pagesize A4 foo.xml                   # use A4 page size
```

You can also read from a pipeline (for example output from the decksh command)

	decksh foo.dsh | pdfdeck -stdout - > f.pdf     # get data from another command, write to stdout
*/
package main
