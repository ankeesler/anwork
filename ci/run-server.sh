#!/bin/bash

set -euo pipefail

tmp_dir=/tmp
stdout_file="$tmp_dir/service.stdout.log"
stderr_file="$tmp_dir/service.stderr.log"
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
echo "logging to file $stdout_file and $stderr_file"

creds="/tmp/creds"
private_key="$(cat $private_key_file)"
echo "export ANWORK_API_ADDRESS=localhost:$port" > "$creds"
echo "export ANWORK_API_PRIVATE_KEY='$private_key'" >> "$creds"
echo "export ANWORK_API_SECRET=$secret" >> "$creds"
echo "creds written to $creds"

go build -o "$service_binary" ./cmd/service/main.go
ANWORK_API_PUBLIC_KEY="$public_key" \
  ANWORK_API_SECRET="$secret" \
  PORT="$port" \
  "$service_binary" 1>"$stdout_file" 2>"$stderr_file"
