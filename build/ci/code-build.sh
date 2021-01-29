#!/bin/sh
set -ex

ciPath=$(cd `dirname $0`; pwd)
rootPath=$(cd $ciPath/../../; pwd)

export GOPROXY=https://goproxy.cn,direct
export GO111MODULE="auto"
go mod download

export CGO_ENABLED=0

export GOOS=linux
export GOARCH=amd64
go build -o "${ciPath}/output/linux-amd64/esctl" ${rootPath}/main.go

export GOOS=darwin
export GOARCH=amd64
go build -o "${ciPath}/output/darwin-amd64/esctl" ${rootPath}/main.go