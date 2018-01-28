#!/usr/bin/env bash

set -e -x -u

echo "$PWD"

export GOPATH="$PWD"

go version
go env

mkdir -p src/github.com/ankeesler/
cp -r ./anwork src/github.com/ankeesler/
cd src/github.com/ankeesler/anwork
go get -t -v ./...
go vet ./...
ci/check-fmt.sh
go test ./...

