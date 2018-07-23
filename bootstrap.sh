#!/bin/bash

set -ex

bootstrap_linux() {
  apt-get update
  apt-get install git golang -y
}

bootstrap_common() {
  go get -t github.com/ankeesler/anwork/...
  go install github.com/onsi/ginkgo/ginkgo
}

case `uname` in
  Linux)
    bootstrap_linux
    bootstrap_common
    ;;
  *)
    echo 'Your OS is not supported!'
    exit 1
    ;;
esac
