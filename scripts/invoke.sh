#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"
payload=${2:-'{"name":"Foo"}'}

aws lambda invoke \
  --function-name "$function_name" \
  --no-cli-pager \
  --cli-binary-format raw-in-base64-out \
  --payload "${payload}" \
  response.json && cat response.json

rm response.json
