// example deck

opts="-fulldeck"
a=100
b=200
c=a+b
thing="200"
c=c+(a*b)/100

deck begin
	canvas 1200 900
	slide begin white black
		ctext "Deck elements" 50 90 5
		cimage "follow.jpg" "Dreams" 70 60 640 480 50 "https://budnitzbicycles.com"

		blist 10 70 3
			li "text, image, list"
			li "rect, ellipse, polygon"
			li "line, arc, curve"
		elist

		rect    15 20 8 6              "rgb(127,0,0)"
		ellipse 27.5 20 8 6            "rgb(0,127,0)"
		polygon "37 37 45" "17 23 20"  "rgb(0,0,127)"
		line    50 20 60 20 0.25       "rgb(127,0,0)"
		arc     70 20 10 8 0 180 0.25  "rgb(0,127,0)"
		curve   80 20 95 30 90 20 0.25 "rgb(0,0,127)"
		ctext "rect"     15 15 1
		ctext "ellipse"  27.5 15 1
		ctext "polycon"  40 15 1
		ctext "line"     55 15 1
		ctext "arc"      70 15 1
		ctext "curve"    85 15 1
	slide end
	slide begin
		dchart -fulldeck=f AAPL.d
	slide end
deck end
