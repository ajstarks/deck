/*
pngdeck converts deck markup to PNG

Command options:
	-fontdir string
	  	directory for fonts (defaults to DECKFONTS environment variable) (default "/home/ajstarks/TTF")
	-grid float
	  	draw a percentage grid on each slide
	-mono string
	  	mono font (default "FiraMono-Regular")
	-outdir string
	  	output directory (default ".")
	-pages string
	  	page range (first-last) (default "1-1000000")
	-pagesize string
	  	pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen (default "Letter")
	-sans string
	  	sans font (default "FiraSans-Regular")
	-serif string
	  	serif font (default "Charter-Regular")
	-symbol string
	  	symbol font (default "ZapfDingbats")
*/
package main
