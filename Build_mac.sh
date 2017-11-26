#!/bin/bash
echo 'Build regener for mac....'
chmod 755 *
export GOARCH="amd64"
export GOOS="darwin"

go build -o $GOPATH/src/regener/regener
if [ -f "regener" ];then
VERSION=$(git describe --abbrev=4 --dirty --always --tags)
APPNAME="regener_v"$VERSION""
mv regener $APPNAME 
echo "New build completed:$APPNAME"
else
echo "New build Failed:$APPNAME" 
fi