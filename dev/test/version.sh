#!/bin/env bash

set -e

DATA=`curl -ks -X 'GET' \
  'https://localhost:8643/version' \
  -H 'accept: application/json'
`
echo "$DATA" | jq .version -r
if [ `echo $DATA | jq -r .version` != "0.0.0" ]; then
  exit 1
fi
