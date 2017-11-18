#!/bin/sh

echo "`basename $0`: updating file..."
git config --global user.email "travis@travis-ci.org"
git config --global user.name "Travis CI"

git checkout master

cat "another line at `date`" > file.txt
git add file.txt
git commit -m "Added another line to file.txt at `date`"

git push
