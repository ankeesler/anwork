#!/bin/bash

set -e

HERE=`dirname $0`

usage() {
    echo "usage: `basename $0` <version-number>"
}

if [[ -z "$1" ]]; then
    usage
    exit 1
fi

version_number=$1
out_dir=$HERE/../integration/data
context=v$version_number-context

$HERE/build.sh -h generate-test-data-hash -d generate-test-data-date -o ./anwork
./anwork -o $out_dir -c $context create task-a
./anwork -o $out_dir -c $context create task-b
./anwork -o $out_dir -c $context create task-c
./anwork -o $out_dir -c $context set-running task-a
./anwork -o $out_dir -c $context set-blocked task-b
./anwork -o $out_dir -c $context set-finished task-a
./anwork -o $out_dir -c $context set-running task-c
rm anwork
