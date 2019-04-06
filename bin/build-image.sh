#!/usr/bin/env bash

set -x
set -e

if (( "$#" != 1 ))
then
    echo "Tag has to be provided"
    exit 1
fi

docker build --rm --pull --force-rm --tag "formapro/crony:$1" .

docker push "formapro/crony:$1"
