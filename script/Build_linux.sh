#!/bin/bash
echo 'Build regener for mac....'
export GOARCH="386"
export GOOS="linux"
go build -o regener -v ../