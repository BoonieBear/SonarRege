#!/bin/bash
echo 'Build regener for mac....'
chmod 755 *
./increasebuild.sh
FILENAME="buildnumber.h"
while read LINE
do
let BUILD=$LINE
done < $FILENAME
export GOARCH="amd64"
export GOOS="darwin"
APPNAME="regener_v"$BUILD""
cd ..
go build -o regener
if [ -f "regener" ];then
mv regener $APPNAME 
echo "New build completed:$APPNAME"
else
echo "New build Failed:$APPNAME" 
fi