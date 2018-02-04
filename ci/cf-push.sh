#!/bin/sh

set -e -x

GOOS=linux GOARCH=amd64 go build cmd/service/main.go
cf push anwork_service -c './main' -b binary_buildpack
rm main
