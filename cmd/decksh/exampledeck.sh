	// Example deck
	deck begin
		notecolor="maroon"
		notesize=1.8
		notefont="mono"
		iw=640
		ih=480
		imscale=55
		c1="red"
		c2="green"
		c3="blue"
		slide begin "white" "black"
			ctext "Deck elements" 50 90 5
			cimage "follow.jpg" "Dreams" 72 55 iw ih imscale "https://budnitzbicycles.com"

			// List
			blist 10 75 3
				li "text, image, list"
				li "rect, ellipse, polygon"
				li "line, arc, curve"
			elist

			// Graphics
			gy=10
			notey=17
			rect    15 gy 8 6              c1
			ellipse 27.5 gy 8 6            c2
			polygon "37 37 45" "7 13 10"   c3
			line    50 gy 60 gy 0.25       c1
			arc     70 gy 10 8 0 180 0.25  c2
			curve   80 gy 95 25 90 gy 0.25 c3

			// Annotations
			ctext "text"	50 97 notesize notefont notecolor
			ctext "image"	72 80 notesize notefont notecolor
			ctext "list"	5 67 notesize notefont notecolor
			ctext "chart"	5 45 notesize notefont notecolor
			ctext "rect"	15 notey notesize notefont notecolor
			ctext "ellipse"	27.5 notey notesize notefont notecolor
			ctext "polygon"	40 notey notesize notefont notecolor
			ctext "line"	55 notey notesize notefont notecolor
			ctext "arc"		70 notey notesize notefont notecolor
			ctext "curve"	85 notey notesize notefont notecolor

			// Chart
			chartleft=10
			chartright=45
			top=50
			bottom=35
			dchart -fulldeck=f  -left chartleft -right chartright -top top -bottom bottom -textsize 1 -color tan -xlabel=2  -barwidth 1.5 AAPL.d 
		slide end
deck end
