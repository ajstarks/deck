#!/bin/bash
. $HOME/Library/deckfuncs.sh

opts="-vol=f -volop=30 -title=f -min=0 -max=800 -dot=f -val=f -bar=f -line -xlabel=0 -fulldeck=f"
deck
	slide
		dchart $opts -color=blue -yrange=0,750,150 -yaxis -xlabel=1 -grid  AAPL.d
		dchart $opts -color=red GOOG.d
		dchart $opts -color=green MSFT.d
		dchart $opts -color=orange AMZN.d
		dchart $opts -color=gray FB.d
		legend "AAPL" 10 20 1 "sans" "blue" 
		legend "GOOG" 20 20 1 "sans" "red" 
		legend "MSFT" 30 20 1 "sans" "green" 
		legend "AMZN" 40 20 1 "sans" "orange" 
		legend "FB"   50 20 1 "sans" "gray" 
	eslide
edeck
