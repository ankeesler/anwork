#!/bin/sh

# This script will package the anwork executable (found in cmd/anwork) in the zip file format
# specified in the anwork_testing repo.

ME=`basename $0`

note() {
    echo ">>> $ME: NOTE: $1"
}

error() {
    echo ">>> $ME: ERROR: $1"
    exit 1
}

build() {
    exec="anwork_$1_$2"
    GOOS="$1" GOARCH="$2" go build -o "$exec" github.com/ankeesler/anwork/cmd/anwork
    if [ $? -ne 0 ]; then
        error "failed to build anwork executable"
    fi
    echo "$exec"
}

HERE=`basename $PWD`
if [ "$HERE" != "anwork" ]; then
    error "this script must be run from the root of the anwork directory"
fi

version=$(awk '/const Version =/ {print $NF}' cmd/anwork/command/command.go)
note "using version $version"

root="anwork-$version"
mkdir $root
mkdir $root/bin
cp ci/anwork $root/bin/
cp $(build darwin amd64) $root/bin/
cp $(build linux amd64) $root/bin/
cp -R doc $root/
note "created staging directory $root"

zipfile="$root.zip"
zip -r $zipfile $root
note "created zip file $zipfile"

rm -r $root
note "removed staging directory $root"
