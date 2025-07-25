#!/bin/env bash
set -e

DATA=`
curl -ks -X 'PUT' \
  'https://localhost:8643/sessions/current' \
  -H 'accept: application/json' \
  -H 'Accept-Version: 0.0.0' \
  -H 'Content-Type: application/json' \
  -d "{\"token\": { \"access\": \"${ACCESS_TOKEN}\", \"refresh\": \"${REFRESH_TOKEN}\" }}"
`

echo export ACCESS_TOKEN=`echo "$DATA" | jq .tokens.access -r`
echo export REFRESH_TOKEN=`echo "$DATA" | jq .tokens.refresh -r`
