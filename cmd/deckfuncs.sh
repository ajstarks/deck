#
# Shell functions for deck markup
#
function deck {
	case $1 in
		begin) echo "<deck>";;
		end) echo "</deck>";;
		*) echo "<deck>";;
	esac
}

function slide {
	case $1 in
		begin) echo "<slide bg=\"$2\" fg=\"$3\">";;
		end) echo "</slide>";;
		*) echo "<slide>";;
	esac
}

function canvas {
	echo "<canvas width=\"$1\" height=\"$2\"/>"
}

function text {
	echo "<text xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" color=\"$6\" opacity=\"$7\">$1</text>"
}

function ctext {
	echo "<text xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" color=\"$6\" opacity=\"$7\" align=\"c\">$1</text>"
}

function etext {
	echo "<text xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" color=\"$6\" opacity=\"$7\" align=\"e\">$1</text>"
}

function textfile {
	echo "<text file=\"$1\" xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" color=\"$6\" opacity=\"$7\"/>"
}

function textblock {
	echo "<text xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" wp=\"$6\" color=\"$7\" opacity=\"$8\" type=\"block\">$1</text>"
}

function textcode {
	echo "<text xp=\"$2\" yp=\"$3\" sp=\"$4\" font=\"$5\" wp=\"$6\" color=\"$7\" opacity=\"$8\" type=\"code\">$1</text>"
}

function list {
	echo "<list xp=\"$1\" yp=\"$2\" sp=\"$3\" color=\"$4\" opacity=\"$5\">"
}

function blist {
	echo "<list xp=\"$1\" yp=\"$2\" sp=\"$3\" color=\"$4\" opacity=\"$5\" type=\"bullet\">"
}

function nlist {
	echo "<list xp=\"$1\" yp=\"$2\" sp=\"$3\" color=\"$4\" opacity=\"$5\" type=\"number\">"
}

function li {
	echo "<li>$*</li>"
}

function elist {
	echo "</list>"
}

function listend {
	echo "</list>"
}

function image {
	echo "<image name=\"$1\" xp=\"$2\" yp=\"$3\" width=\"$4\" height=\"$5\" scale=\"$6\"/>"
}

function cimage {
	echo "<image name=\"$1\" caption=\"$2\" xp=\"$3\" yp=\"$4\" width=\"$5\" height=\"$6\" scale=\"$7\"/>"
}

function line {
	echo "<line xp1=\"$1\" yp1=\"$2\" xp2=\"$3\" yp2=\"$4\" sp=\"$5\" color=\"$6\" opacity=\"$7\"/>"
}

function rect {
	echo "<rect xp=\"$1\" yp=\"$2\" wp=\"$3\" hp=\"$4\" color=\"$5\" opacity=\"$6\"/>" 
}

function square {
	echo "<rect xp=\"$1\" yp=\"$2\" wp=\"$3\" hr=\"100\" color=\"$4\" opacity=\"$5\"/>" 
}

function ellipse {
	echo "<ellipse xp=\"$1\" yp=\"$2\" wp=\"$3\" hp=\"$4\" color=\"$5\" opacity=\"$6\"/>" 
}

function circle {
	echo "<ellipse xp=\"$1\" yp=\"$2\" wp=\"$3\" hr=\"100\" color=\"$4\" opacity=\"$5\"/>" 
}

function arc {
	echo "<arc xp=\"$1\" yp=\"$2\" wp=\"$3\" hp=\"$4\" a1=\"$5\" a2=\"$6\" color=\"$7\" opacity=\"$8\"/>"
}

function curve {
	echo "<curve xp1=\"$1\" yp1=\"$2\" xp2=\"$3\" yp2=\"$4\" xp3=\"$5\" yp3=\"$6\" color=\"$7\" opacity=\"$8\"/>"
}

function polygon {
	echo "<polygon xc=\"$1\" yc=\"$2\" color=\"$3\" opacity=\"$4\"/>"
}

function legend {
		text "$1" $(($2 + 2)) $3 $4 $5
		circle $2 $3".5" $4 "$6"
}
