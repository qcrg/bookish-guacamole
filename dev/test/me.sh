#!/bin/env bash
set -e

DATA=`
curl -ks -X 'GET' \
  'https://localhost:8643/me' \
  -H 'accept: application/json' \
  -H "Accept-Version: ${VERSION}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
`

RES_USER_ID=`echo "$DATA" | jq .user_id -r`
if [ "$RES_USER_ID" != "$USER_ID" ]; then
  echo $DATA
  exit 1
else
  echo $RES_USER_ID
fi
