#!/bin/bash
echo 'Build regener for windows....'
export GOARCH="amd64"
export GOOS="windows"
cd ..
go build -o regener.exe
if [ -f "regener.exe" ];then
VERSION=$(git describe --abbrev=4 --dirty --always --tags)
APPNAME="regener_v"$VERSION".exe"
mv regener.exe $APPNAME 
echo "New build completed:$APPNAME"
else
echo "New build Failed:$APPNAME" 
fi