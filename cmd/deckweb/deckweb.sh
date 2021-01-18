#!/bin/sh
# deck -- command line interface to the deck API
if [[ $DECKS == "" ]]
then
	DECKS="http://localhost:1958"
fi
srv=$DECKS
duration="1s"
textsize=1.4
case $1 in
	list)
		gttp -raw GET $srv/deck/?filter=$2
		;;
	start|play)
		if [[ $3 != "" ]]
		then
			duration=$3
		fi
		gttp -raw POST $srv/deck/$2?cmd=$duration
		;;
	stop)
		gttp -raw POST $srv/deck/f.xml?cmd=stop
		;;
	up|upload)
		shift
		for f in $*
		do
			deck=`basename $f`
			gttp -raw POST $srv/upload/ Deck:$deck -@$f
		done
		;;
	del|delete|remove)
		shift
		for f in $*
		do
			gttp -raw DELETE $srv/deck/$f Deck:$f
		done
		;;
	video)
		gttp -raw POST $srv/media/ Media:$2
		;;
	table)
		f=`basename $2`
		if [[ $3 != "" ]]
		then
			textsize=$3
		fi
		gttp -raw POST $srv/table/?textsize=$textsize Deck:$f -@$2
		;;
	*)
		echo "Play:     deck [start|play] file [duration]" 1>&2
		echo "Stop:     deck stop" 1>&2
		echo "List:     deck list [image|deck|video]" 1>&2
		echo "Upload:   deck upload file..." 1>&2
		echo "Remove:   deck [del|delete|remove] file..." 1>&2
		echo "Video:    deck video file" 1>&2
		echo "Table:    deck table file [textsize]" 1>&2

		exit 1
		;;
esac
