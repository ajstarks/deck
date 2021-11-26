#!/bin/bash
base="github.com/ajstarks/deck/cmd"
for b in pdfdeck pngdeck svgdeck
do
	echo -n "$b - "
	for o in "linux/amd64" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/386"
	do
		echo -n "$o "
		goos=$(echo $o|cut -f1 -d /)
		goarch=$(echo $o|cut -f2 -d /)
		if test "$goos" = "windows"
		then
			exe="${goos}-${goarch}-${b}.exe"
		else
			exe="${goos}-${goarch}-${b}"
		fi
		CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build  -o $exe $base/$b
	done
	echo
done
