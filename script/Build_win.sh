#!/bin/bash
echo 'Build regener for windows....'
export GOARCH="amd64"
export GOOS="windows"
go build -o regener -v ../src/regener.go