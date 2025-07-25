#!/bin/env bash

curl -ks -X 'GET' \
  'https://localhost:8643/health' \
  -H 'accept: */*'
