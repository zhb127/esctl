#!/bin/sh
#set -ex

currPath=$(cd `dirname $0`; pwd)
rootPath=$(cd ${currPath}/../; pwd)
scriptsUtilPath=$(cd ${rootPath}/scripts/util ; pwd)

dotenvPath=$1
if [ "$dotenvPath" == "" ]; then
  dotenvPath="${currPath}/.env"
fi

. ${scriptsUtilPath}/load-dotenv.sh "${dotenvPath}"

set -o pipefail

go test -v -cover -count=1 -failfast ${rootPath}/pkg/... | { grep -v 'no test files'; true; }
go test -v -cover -count=1 -failfast ${rootPath}/internal/... | { grep -v 'no test files'; true; }

set +o pipefail