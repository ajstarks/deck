// comprehensive tests
l1=20
l2=50
l3=80
op1=70
op2=50
op3=30
ts1=4
deck begin
	canvas 1200 900
	slide begin "rgb(240,240,240)"
		text "one" 		l1 80 ts1
		text "two" 		l1 70 ts1 "serif"
		text "three" 	l1 60 ts1 "mono" "red"
		text "four"  	l1 50 ts1 "sans" "blue" op1


		ctext "one" 	l2 80 ts1
		ctext "two" 	l2 70 ts1 "serif"
		ctext "three" 	l2 60 ts1 "mono" "red"
		ctext "four"  	l2 50 ts1 "sans" "blue" op2
		

		etext "one"		l3 80 ts1
		etext "two"		l3 70 ts1 "serif"
		etext "three"	l3 60 ts1 "mono" "red"
		etext "four"	l3 50 ts1 "sans" "blue" op3
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
	
	lsize=2
	lx1=20
	lx2=40
	lx3=60
	slide begin
		list lx1 90 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		list lx2 90 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 90 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 70 lsize "sans"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx2 70 lsize "serif"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 70 lsize "mono"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 50 lsize "sans" "red"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx2 50 lsize "serif" "green"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 50 lsize "mono" "blue"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 30 lsize "sans" "red" op1
			li "one"
			li "two"
			li "three"
		elist
		
		list lx2 30 lsize "serif" "green" op2
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 30 lsize "mono" "blue" op3
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
	
	// Example deck
	notecolor="maroon"
	notesize=1.8
	notefont="mono"
	imscale=55
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
		chartleft=10
		chartright=45
		top=50
		bottom=35
		dchart -fulldeck=f  -left chartleft -right chartright -top top -bottom bottom -textsize 1 -color tan -xlabel=2  -barwidth 1.5 AAPL.d 
	slide end
deck end
