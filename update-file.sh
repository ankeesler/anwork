#!/bin/sh

echo "`basename $0`: updating file..."
sed -i .bak -e 's/version=\(.*\)/version=7/' file.txt
git add file.txt
git commit -m 'Update file to version 7.'
git push
