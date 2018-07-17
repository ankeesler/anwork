#!/bin/sh

set -e

ME=`basename $0`
HERE=`dirname $0`

log() {
    echo "[$ME]: $@"
}

DIR=$(mktemp -d)
log "Using temporary directory $DIR..."

GOOS=linux GOARCH=amd64 go build -o $DIR/main ./cmd/service/main.go
log "Built $DIR/main binary..."

cf push anwork_service -p $DIR -f $HERE/manifest.yml
log "Pushed anwork_service..."

rm -rf $DIR
log "Removed temp directory..."
