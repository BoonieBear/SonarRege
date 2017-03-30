#!/bin/bash
echo 'Build regener for windows....'
increasebuild.sh
FILENAME="buildnumber.h"
while read LINE
do
let BUILD=$LINE
done < $FILENAME

export GOARCH="amd64"
export GOOS="windows"
APPNAME="regener_v"$BUILD".exe"
cd ..
go build -o regener.exe
if [ -f "regener.exe" ];then
mv regener.exe $APPNAME 
echo "New build completed:$APPNAME"
else
echo "New build Failed:$APPNAME" 
fi