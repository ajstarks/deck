// comprehensive tests
deck begin
	canvas 1200 900
	tx1=5
	tx2=35
	tx3=65
	ty=92
	tb="Now is the time for all good men to come to the aid of the party & 'do it now'"
	// Text Functions
	slide begin
		textblock tb tx1 ty 20 2
		textblock tb tx2 ty 15 2 "serif"
		textblock tb tx3 ty 10 2 "mono" "red"
		
		textfile "AAPL.d" tx1 50 2
		textfile "AAPL.d" tx2 50 2 "serif"
		textfile "AAPL.d" tx3 50 2 "mono" "red"
	
		textcode "code/hw.go" tx1 75 20 1
		textcode "code/hw.go" tx2 75 20 1 "red"
	slide end

	l1=20
	l2=50
	l3=80
	op1=70
	op2=50
	op3=30
	ts1=4
	slide begin "rgb(240,240,240)"
		line l1 0 l1 100 0.1
		line l2 0 l2 100 0.1
		line l3 0 l3 100 0.1
		text "one"   l1 80 ts1
		text "two"   l1 70 ts1 "serif"
		text "three" l1 60 ts1 "mono" "red"
		text "four"  l1 50 ts1 "sans" "blue" op1


		ctext "one"   l2 80 ts1
		ctext "two"   l2 70 ts1 "serif"
		ctext "three" l2 60 ts1 "mono" "red"
		ctext "four"  l2 50 ts1 "sans" "blue" op2
		

		etext "one"   l3 80 ts1
		etext "two"   l3 70 ts1 "serif"
		etext "three" l3 60 ts1 "mono" "red"
		etext "four"  l3 50 ts1 "sans" "blue" op3
	slide end
	
	midx=50
	midy=50
	iw=640
	ih=480
	s1=50
	s2=20
	imfile="follow.jpg"
	imlink="https://budnitzbicycles.com"
	// Images
	slide begin
		image imfile midx midy iw ih
		image imfile midx midy iw ih s1
		image imfile midx midy iw ih s2 imlink
	slide end
	
	slide begin "black" "white"
		cimage imfile "LARGE" midx midy iw ih
		cimage imfile "MEDIUM" midx midy iw ih s1
		cimage imfile "SMALL" midx midy iw ih s2 imlink
	slide end
	
	lsize=2
	lx1=20
	lx2=40
	lx3=60
	ly1=90
	ly2=70
	ly3=50
	ly4=30
	// Lists
	slide begin
		list lx1 ly1 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		blist lx2 ly1 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 ly1 lsize
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 ly2 lsize "sans"
			li "one"
			li "two"
			li "three"
		elist
		
		blist lx2 ly2 lsize "serif"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 ly2 lsize "mono"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 ly3 lsize "sans" "red"
			li "one"
			li "two"
			li "three"
		elist
		
		blist lx2 ly3 lsize "serif" "green"
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 ly3 lsize "mono" "blue"
			li "one"
			li "two"
			li "three"
		elist
		
		list lx1 ly4 lsize "sans" "red" op1
			li "one"
			li "two"
			li "three"
		elist
		
		blist lx2 ly4 lsize "serif" "green" op2
			li "one"
			li "two"
			li "three"
		elist
		
		nlist lx3 ly4 lsize "mono" "blue" op3
			li "one"
			li "two"
			li "three"
		elist
	slide end
	
	c1="red"
	c2="blue"
	c3="green"
	shapeop=30
	// Shapes
	slide begin
		polygon	   "15 20 25" "90 95 90"
		polygon	   "35 40 45" "90 95 90" c1
		polygon	   "55 60 65" "90 95 90" c2 shapeop 
		
		rect	   l1 80 10 5
		rect	   40 80 10 5 c1
		rect	   60 80 10 5 c2 shapeop 
		
		square	   l1 70 5
		square	   40 70 5 c1
		square	   60 70 5 c2 shapeop 
		
		ellipse	   l1 60 10 5
		ellipse	   40 60 10 5 c1
		ellipse	   60 60 10 5 c2 shapeop 
		
		circle	   l1 50 5
		circle	   40 50 5 c1
		circle	   60 50 5 c2 shapeop 
		
		line	   15 35 25 40
		line	   35 35 45 40 1 c1
		line	   55 35 65 40 1 c2
		line	   75 35 85 40 1 c3 shapeop 
		
		arc        20 25 10 5 0 180
		arc        40 25 10 5 0 180 1 c1
		arc        60 25 10 5 0 180 1 c2
		arc        80 25 10 5 0 180 1 c3 shapeop 
		
		curve	   15 15 10 25 25 15
		curve	   35 15 30 25 45 15 1
		curve	   55 15 45 25 65 15 1 c2
		curve	   75 15 65 25 85 15 1 c3 shapeop 
	slide end
	
	// Example deck
	imscale=58
	dtop=87
	chartleft=10
	chartright=42
	chartop=42
	chartbottom=28
	imy=50
	opts="-fulldeck=f -textsize 1  -xlabel=2  -barwidth 1.5"

	slide begin "white" "black"
		ctext     "Deck elements" 50 dtop 5
		cimage    "follow.jpg" "Dreams" 72 imy iw ih imscale imlink
		textblock "Budnitz #1, Plainfield, NJ, May 10, 2015" 55 35 10 1 "serif" "white"

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

		// Chart
		dchart -left chartleft -right chartright -top chartop -bottom chartbottom opts AAPL.d 
	slide end
	
	
	slide begin "white" "black"
		ctext     "Deck elements" 50 dtop 5
		cimage    "follow.jpg" "Dreams" 72 imy iw ih imscale imlink
		textblock "Budnitz #1, Plainfield, NJ, May 10, 2015" 55 35 10 1 "serif" "white"

		// List
		blist 10 75 3
			li "text, image, list"
			li "rect, ellipse, polygon"
			li "line, arc, curve"
		elist

		// Graphics
		gy=10
		rect    15 gy 8 6              c1
		ellipse 27.5 gy 8 6            c2
		polygon "37 37 45" "7 13 10"   c3
		line    50 gy 60 gy 0.25       c1
		arc     70 gy 10 8 0 180 0.25  c2
		curve   80 gy 95 25 90 gy 0.25 c3

		// Annotations
		ns=5
		nc="gray"
		nf="serif"
		nop=30
		ctext "text"	50 95		ns nf nc nop
		ctext "image"	72 80		ns nf nc nop
		ctext "list"	25 80		ns nf nc nop
		ctext "chart"	25 50		ns nf nc nop

		ns=2
		notey=17
		ctext "rect"	15 notey	ns nf nc
		ctext "ellipse"	27.5 notey	ns nf nc
		ctext "polygon"	40 notey	ns nf nc
		ctext "line"	55 notey	ns nf nc
		ctext "arc"		70 notey	ns nf nc
		ctext "curve"	85 notey	ns nf nc

		// Chart
		dchart -left chartleft -right chartright -top chartop -bottom chartbottom opts AAPL.d 
	slide end
deck end
