#!/usr/bin/env bash -e

echo "=> build"
GOOS=linux GOARCH=arm64 go build -o bootstrap

echo "=> create deployment package"
zip function.zip bootstrap
