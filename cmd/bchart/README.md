
# bchart - make barcharts in the deck format #

bchart reads data from the standard input file, generating deck markup to standard output.
 
Input is a tab-separated list of text,data pairs
where text is an arbitrary string, and data is intepreted as a floating point value.
A line beginning with "#" is parsed as a title, with the title text beginning after the "#".

For example:

	# PDF File Sizes
	casino.pdf	410907
	countdown.pdf	157784
	deck-12x8.pdf	837831
	deck-dejavu.pdf	1601595
	deck-fira-4x3.pdf	1196167
	deck-fira.pdf	1195517
	deck-gg.pdf	978688
	deck-gofont.pdf	1044627


bchart is useful in pipeline:

	ls -l *.pdf | awk 'BEGIN { print "# PDF File Sizes" } NF > 5 { print $NF "\t" $5 }' | sort -nr -k2 | bchart > f.xml && pdfdeck f.xml

 The command line options are:

	  -color barcolor (default "rgb(175,175,175)")
	  -datafmt data format (default "%.1f")
	  -dmin zero minimum
	  -dot draw a line and dot instead of a solid bar
	  -left left margin (default 20)
	  -textsize text size (default 1.2)
	  -top top of the chart (default 90)
