#!/bin/bash
FILENAME="buildnumber.h"
BUILD=0
while read LINE
do
echo -e "Old build number:$LINE"
let BUILD=$LINE+1
done  < $FILENAME

echo $BUILD>buildnumber_tmp.h
rm -f buildnumber.h
mv buildnumber_tmp.h buildnumber.h
echo -e "New build number:$BUILD"
