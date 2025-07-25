#!/bin/env bash
set -e

curl -ks -X 'DELETE' \
  'https://localhost:8643/sessions/current' \
  -H 'accept: application/json' \
  -H "Accept-Version: ${VERSION}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
