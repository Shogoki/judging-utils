#!/bin/sh
REPO=$1
for f in *.md; do
	if [[ "$f" != "comment.md" ]]
	then
		NUM=$(echo $f | cut -f 1 -d .) # cutoff md
		COMMAND="gh issue comment $NUM -R $REPO -F comment.md"
		echo $COMMAND
	fi
done
