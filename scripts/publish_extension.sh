#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"
extension_name="${2:-go-extension-example}"

echo "=> build"
GOOS=linux GOARCH=arm64 go build -trimpath -o build/extensions/$extension_name main.go

cd build
zip -r extension.zip extensions/

echo "=> publish layer"
aws lambda publish-layer-version \
 --layer-name $extension_name \
 --zip-file  "fileb://extension.zip" \
 --no-cli-pager > publish_layer_output

layer_arn=$(cat publish_layer_output | jq  ".LayerVersionArn" -r)

echo "=> update function"
aws lambda update-function-configuration \
    --function-name $function_name \
    --layers $layer_arn \
    --no-cli-pager > update_function_output


rm publish_layer_output
rm update_function_output
