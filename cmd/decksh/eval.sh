// comprehensive tests
l1=20
l2=50
l3=80
op1=50
op2=30
op3=10
ts1=4
deck begin
	canvas 1200 900
	slide begin "rgb(240,240,240)"
		text "one" 		l1 80 ts1
		text "two" 		l1 70 ts1 "serif"
		text "three" 	l1 60 ts1 "mono" "red"
		text "four"  	l1 50 ts1 "sans" "blue" 50


		ctext "one" 	l2 80 ts1
		ctext "two" 	l2 70 ts1 "serif"
		ctext "three" 	l2 60 ts1 "mono" "red"
		ctext "four"  	l2 50 ts1 "sans" "blue" 30
		

		etext "one"		l3 80 ts1
		etext "two"		l3 70 ts1 "serif"
		etext "three"	l3 60 ts1 "mono" "red"
		etext "four"	l3 50 ts1 "sans" "blue" 10
	slide end
	
	midx=50
	midy=50
	iw=640
	ih=480
	s1=50
	s2=20
	imfile="follow.jpg"
	imlink="https://budnitzbicycles.com"
	slide begin
		image imfile midx midy iw ih
		image imfile midx midy iw ih s1
		image imfile midx midy iw ih s2 imlink
	slide end
	
	slide begin
		cimage "follow.jpg" "BIG" midx midy iw ih
		cimage "follow.jpg" "MED" midx midy iw ih s1
		cimage "follow.jpg" "SMALL" midx midy iw ih s2 imlink
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
	
	c1="red"
	c2="blue"
	c3="green"
	slide begin
		polygon	"15 20 25" "90 95 90"
		polygon	"35 40 45" "90 95 90" c1
		polygon	"55 60 65" "90 95 90" c2 30
		
		rect	l1 80 10 5
		rect	40 80 10 5 c1
		rect	60 80 10 5 c2 30
		
		square	l1 70 5
		square	40 70 5 c1
		square	60 70 5 c2 30
		
		ellipse	l1 60 10 5
		ellipse	40 60 10 5 c1
		ellipse	60 60 10 5 c2 30
		
		circle	l1 50 5
		circle	40 50 5 c1
		circle	60 50 5 c2 30
		
		line	15 35 25 40
		line	35 35 45 40 1 c1
		line	55 35 65 40 1 c2
		line	75 35 85 40 1 c3 30
		
		arc		20 25 10 5 0 180
		arc		40 25 10 5 0 180 1 c1
		arc		60 25 10 5 0 180 1 c2
		arc		80 25 10 5 0 180 1 c3 30
		
		curve	15 15 10 25 25 15
		curve	35 15 30 25 45 15 1
		curve	55 15 45 25 65 15 1 c2
		curve	75 15 65 25 85 15 1 c3 30
	slide end
	
	slide begin "white" "black"
		ctext "Deck elements" 50 90 5
		cimage "follow.jpg" "image" 72 55 iw ih 50 "https://budnitzbicycles.com"

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
		//dchart -left=10 -right=45 -top=50 -bottom=30 -fulldeck=f -textsize=0.7 -color=tan  -barwidth=1.5 AAPL.d  
	slide end
deck end
