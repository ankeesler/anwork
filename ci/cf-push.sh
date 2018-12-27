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

private_key="$(openssl genrsa 2048)"
public_key="$(echo "$private_key" | openssl rsa -pubout -outform PEM -inform PEM)"
secret="$(openssl rand -base64 32)"
log "Generated RSA key and secret..."

manifest="$DIR/manifest.yml"
go run ./cmd/genmanifest/main.go "$public_key" "$secret" > "$manifest"
log "Wrote manifest file $manifest..."

cf push anwork_service -p $DIR -f "$manifest"
log "Pushed anwork_service..."

rm -rf $DIR
log "Removed temp directory..."

ANWORK_API_PRIVATE_KEY="$private_key" \
  ANWORK_API_SECRET="$secret" \
  ANWORK_API_ADDRESS=$(cf app anwork_service | awk '/routes/ {print $2}') \
  go run ./cmd/anwork/main.go show
log "Passed canary test..."

echo "export ANWORK_API_PRIVATE_KEY=$private_key"
echo "export ANWORK_API_SECRET=$secret"
