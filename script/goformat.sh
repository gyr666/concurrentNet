#!/bin/bash
set -e
set -o pipefail
cd $(dirname $0)/..
SOURCE_LIST_FILE=/tmp/source-files.list
find . -not \( -path './vendor' -prune \) -name '*.go' -not -name '*.gw.go' -not -name '*.pb.go' > $SOURCE_LIST_FILE
goimports $@ -local gunplan.top -d $(cat $SOURCE_LIST_FILE)
