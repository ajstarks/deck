// example deck

deck begin
	canvas 1200 900
	slide begin "rgb(240,240,240)"
		text "one" 		20 80 4
		text "two" 		20 70 4 "serif"
		text "three" 	20 60 4 "mono" "red"
		text "four"  	20 50 4 "sans" "blue" 50


		ctext "one" 	50 80 4
		ctext "two" 	50 70 4 "serif"
		ctext "three" 	50 60 4 "mono" "red"
		ctext "four"  	50 50 4 "sans" "blue" 30
		

		etext "one"		80 80 4
		etext "two"		80 70 4 "serif"
		etext "three"	80 60 4 "mono" "red"
		etext "four"	80 50 4 "sans" "blue" 10
	slide end
	
	slide begin
		image "follow.jpg" 50 50 640 480
		image "follow.jpg" 50 50 640 480 50
		image "follow.jpg" 50 50 640 480 20 "https://budnitzbicycles.com"
	slide end
	
	slide begin
		cimage "follow.jpg" "BIG" 50 50 640 480
		cimage "follow.jpg" "MED" 50 50 640 480 50
		cimage "follow.jpg" "SMALL" 50 50 640 480 20 "https://budnitzbicycles.com"
	slide end
	
	slide begin
		list 10 90 2
			li "one"
			li "two"
			li "three"
		elist
		
		list 30 90 2
			li "one"
			li "two"
			li "three"
		elist
		
		nlist 50 90 2
			li "one"
			li "two"
			li "three"
		elist
		
		list 10 70 2 "sans"
			li "one"
			li "two"
			li "three"
		elist
		
		list 30 70 2 "serif"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist 50 70 2 "mono"
			li "one"
			li "two"
			li "three"
		elist
		
		list 10 50 2 "sans" "red"
			li "one"
			li "two"
			li "three"
		elist
		
		list 30 50 2 "serif" "green"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist 50 50 2 "mono" "blue"
			li "one"
			li "two"
			li "three"
		elist
		
		list 10 30 2 "sans" "red" 50
			li "one"
			li "two"
			li "three"
		elist
		
		list 30 30 2 "serif" "green" 30
			li "one"
			li "two"
			li "three"
		elist
		
		nlist 50 30 2 "mono" "blue" 10
			li "one"
			li "two"
			li "three"
		elist
	slide end
	
	slide begin
		polygon	"15 20 25" "90 95 90"
		polygon	"35 40 45" "90 95 90" "red"
		polygon	"55 60 65" "90 95 90" "blue" 30
		
		rect	20 80 10 5
		rect	40 80 10 5 "red"
		rect	60 80 10 5 "blue" 30
		
		square	20 70 5
		square	40 70 5 "red"
		square	60 70 5 "blue" 30
		
		ellipse	20 60 10 5
		ellipse	40 60 10 5 "red"
		ellipse	60 60 10 5 "blue" 30
		
		circle	20 50 5
		circle	40 50 5 "red"
		circle	60 50 5 "blue" 30
		
		line	15 35 25 40
		line	35 35 45 40 1 "red"
		line	55 35 65 40 1 "blue"
		line	75 35 85 40 1 "green" 30
		
		arc		20 25 10 5 0 180
		arc		40 25 10 5 0 180 1 "red"
		arc		60 25 10 5 0 180 1 "blue"
		arc		80 25 10 5 0 180 1 "green" 30
		
		curve	15 15 10 25 25 15
		curve	35 15 30 25 45 15 1
		curve	55 15 45 25 65 15 1 "blue"
		curve	75 15 65 25 85 15 1 "green" 30
	slide end
	
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
