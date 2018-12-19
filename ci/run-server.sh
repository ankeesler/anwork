#!/bin/bash

set -euo pipefail

tmp_dir=/tmp
log_file="$tmp_dir/service.log"
private_key_file="$tmp_dir/service.key"
public_key_file="$tmp_dir/service.key.pub"
service_binary="$tmp_dir/service"
port=12347

openssl genrsa -out "$private_key_file" 2048
openssl rsa -in "$private_key_file" -pubout -outform PEM -out "$public_key_file"
public_key="$(cat $public_key_file)"
secret="$(openssl rand -base64 32)"

echo
echo "starting server on :$port"
echo "private key stored in $private_key_file"
echo "secret is $secret"
echo "logging to file $log_file"
go build -o "$service_binary" ./cmd/service/main.go
ANWORK_API_PUBLIC_KEY="$public_key" \
  ANWORK_API_SECRET="$secret" \
  PORT="$port" \
  "$service_binary" > "$log_file"
