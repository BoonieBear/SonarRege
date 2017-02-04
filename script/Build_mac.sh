#!/bin/bash
echo 'Build regener for mac....'
export GOARCH="amd64"
export GOOS="darwin"
go build -o regener -v ../