#!/bin/bash

usage() {
    echo "usage: `basename $0` -h hash -d date [-o out]"
}

out="anwork"
while getopts h:d:o: o;
do
    case "$o" in
        h)
            hash="$OPTARG"
            ;;
        d)
            date="$OPTARG"
            ;;
        o)
            out="$OPTARG"
            ;;
        [?])
            usage
            exit 1
    esac
done

if [[ -z "$hash" ]] || [[ -z "$date" ]]; then
    usage
    exit 1
fi

command="go build -ldflags '-X main.buildHash=$hash -X main.buildDate=$date' -o $out github.com/ankeesler/anwork/cmd/anwork"
#echo "command: $command"
eval "$command"
