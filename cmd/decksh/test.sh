// example deck
deck begin
	canvas 1200 900
	slide begin white black
		ctext "Deck elements" 50 90 5
		cimage "follow.jpg" "Dreams" 70 60 640 480 80

		blist 10 70 3
			li "text, image, list"
			li "rect, ellipse, polygon"
			li "line, arc, curve"
		elist

		rect    15 20 8 6             "rgb(127,0,0)"
		ellipse 27.5 20 8 6           "rgb(0,127,0)"
		polygon "37 37 45" "17 23 20" "rgb(0,0,127)"
		line    50 20 60 20
		arc     70 20 10 8 0 180 0   "rgb(0,0,127)"
		curve   80 20 95 30 90 20 0  "rgb(127,127,127)"
	slide end
deck end
