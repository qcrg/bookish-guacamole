#!/bin/env bash
set -e

DATA=`
curl -ks -X 'POST' \
  'https://localhost:8643/auth' \
  -H 'accept: application/json' \
  -H "Accept-Version: ${VERSION}" \
  -H 'Content-Type: application/json' \
  -d "{\"id\":\"${USER_ID}\"}"
`

echo export ACCESS_TOKEN=`echo "$DATA" | jq .tokens.access -r`
echo export REFRESH_TOKEN=`echo "$DATA" | jq .tokens.refresh -r`
