#!/bin/bash

set -e

ME=`basename $0`
HERE=`dirname $0`
VERSION_FILE="$HERE/../runner/command.go"

log() {
    echo "[$ME]: $@"
}

last_version=$(awk '/const Version = / {print $NF}' $VERSION_FILE)
next_version=$(($last_version + 1))

log "Upgrading from v$last_version to v$next_version"

log "Updating version from $VERSION_FILE..."
sed -i .bak -e "s/const Version = $last_version/const Version = $next_version/" $VERSION_FILE
log "...done"

log "Regenerating CLI documentation..."
go generate $VERSION_FILE
log "...done"

log "Updating README with latest version..."
sed -i .bak -e "s/v$last_version/v$next_version/g" $HERE/../README.md
log "...done"

log "Updating backwards compat test data..."
$HERE/generate-test-data.sh $next_version
log "...done"

log "Running tests..."
./ci/test.sh
log "...done"

log "Committing changes..."
git add -AN
git commit -m "Release v$next_version."
log "...done"

log "Adding tag..."
git tag v$next_version
log "...done"
