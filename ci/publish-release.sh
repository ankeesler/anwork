#!/bin/bash

set -e

usage() {
    echo "usage: `basename $0` -t tag -g github-token -a artifact-path -d release-doc"
}

tag=
token=
artifact=
doc=
while getopts t:g:a:d: o
do
    case "$o" in
        t)
            tag="$OPTARG"
            ;;
        g)
            token="$OPTARG"
            ;;
        a)
            artifact="$OPTARG"
            ;;
        d)
            doc="$OPTARG"
            ;;
        [?])
            usage
            exit 1
    esac
done

if [[ -z "$tag" ]] || [[ -z "$token" ]] || [[ -z "$artifact" ]] || [[ -z "$doc" ]]; then
    usage
    exit 1
fi

commit="$(git rev-parse $tag)"
body="$(awk '{printf "%s\\n", $0}' $doc | tr '"' "'")"
data=$(cat <<EOF
{
  "tag_name": "$tag",
  "target_commitish": "$commit",
  "name": "$tag",
  "body": "$body",
  "draft": true
}
EOF
)

set -x
response=$(curl -X POST https://api.github.com/repos/ankeesler/anwork/releases \
     -H "Authorization: token $token" \
     -d "$data")
upload_url="$(echo $response | jq -r .upload_url)"
set +x

name="$tag"_anwork_darwin_amd64
label="anwork%20binary%20%28Mac%20OS%20X%29"
upload_url="${upload_url%\{*}"
set -x
curl -X POST "$upload_url?name=$name&label=$label" \
     -H "Authorization: token $token" \
     -H "Content-Type: application/octet-stream" \
     --data-binary "@$artifact"
set +x
