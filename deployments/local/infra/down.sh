#!/bin/bash
#set -ex

currPath=$(cd `dirname $0`; pwd)

docker-compose -p esctl-infra -f ${currPath}/docker-compose.yaml down --remove-orphans