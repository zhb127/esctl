#!/bin/sh
set -ex

imageName=$1
if [ "$imageName" == "" ]; then
    echo "imageName is empty"
    exit 1
fi

ciPath=$(cd `dirname $0`; pwd)
buildPath=$(cd $ciPath/../; pwd)
rootPath=$(cd $ciPath/../../; pwd)

cd $ciPath

docker build -t "${imageName}" -f  "$buildPath/docker/ci.Dockerfile" $rootPath