// example deck
deck begin
	slide begin "white" "black"
		ctext "Deck elements" 50 90 5
		cimage "follow.jpg" "image" 72 55 640 480 50 "https://budnitzbicycles.com"

		blist 7 75 3
			li "text, image, list"
			li "rect, ellipse, polygon"
			li "line, arc, curve"
		elist

		rect    15 15 8 6              "rgb(127,0,0)"
		ellipse 27.5 15 8 6            "rgb(0,127,0)"
		polygon "37 37 45" "12 18 15"  "rgb(0,0,127)"
		line    50 15 60 15 0.25       "rgb(127,0,0)"
		arc     70 15 10 8 0 180 0.25  "rgb(0,127,0)"
		curve   80 15 95 30 90 15 0.25 "rgb(0,0,127)"
		ctext "rect"     15 10 1
		ctext "ellipse"  27.5 10 1
		ctext "polygon"  40 10 1
		ctext "line"     55 10 1
		ctext "arc"      70 10 1
		ctext "curve"    85 10 1
		dchart -left=10 -right=45 -top=50 -bottom=30 -fulldeck=f -textsize=0.7 -color=tan  -barwidth=1.5 AAPL.d  
	slide end
deck end
