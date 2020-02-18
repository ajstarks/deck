#!/bin/sh
mopts="-fulldeck=f -bar=f -vol -top=90 -bottom 60 -left=20 -right=80 -val=f -title=f"
mfunc -f sine   | dchart $mopts -xlabel=10
mfunc -f cosine | dchart $mopts -xlabel=0 -color=orange