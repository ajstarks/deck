// Example deck
midx=50
midy=50
iw=640
ih=480

imfile="follow.jpg"
imlink="https://budnitzbicycles.com"
imscale=58
dtop=87

opts="-fulldeck=f -textsize 1  -xlabel=2  -barwidth 1.5"
deck
	slide "white" "black"
		ctext "Deck elements" midx dtop 5
		cimage "follow.jpg" "Dreams" 72 midy iw ih imscale imlink
		textblock "Budnitz #1, Plainfield, NJ, May 10, 2015" 55 35 10 1 "serif" "white"

		// List
		blist 10 75 3
			li "text, image, list"
			li "rect, ellipse, polygon"
			li "line, arc, curve"
		elist

		// Graphics
		gy=10
		c1="red"
		c2="blue"
		c3="green"
		rect    15 gy 8 6              c1
		ellipse 27.5 gy 8 6            c2
		polygon "37 37 45" "7 13 10"   c3
		line    50 gy 60 gy 0.25       c1
		arc     70 gy 10 8 0 180 0.25  c2
		curve   80 gy 95 25 90 gy 0.25 c3


		// Chart
		chartleft=10
		chartright=45
		charttop=42
		chartbottom=28
		dchart -left chartleft -right chartright -top charttop -bottom chartbottom opts AAPL.d 
	eslide
edeck