#!/bin/sh

# Check that all of the go sources at the root of the anwork repo are formatted properly via "gofmt."

ME=`basename $0`

HERE=`basename $PWD`
if [ "$HERE" != "anwork" ]; then
    echo "$ME: ERROR: this script must be run from the root of the anwork directory"
    exit 1
fi

o="$(go fmt ./...)"
if [ -z "$o" ]; then
    echo "$ME: PASS"
else
    echo "$ME: FAILURE: the following files are not formatted properly"
    echo "$o"
    exit 1
fi
