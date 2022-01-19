#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"

aws lambda invoke \
  --function-name "$function_name" \
  --no-cli-pager \
  --cli-binary-format raw-in-base64-out \
  --payload '{"name":"Foo"}' \
  response.json && cat response.json

rm response.json
