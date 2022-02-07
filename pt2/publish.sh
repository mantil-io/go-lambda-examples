#!/usr/bin/env bash -e

#echo "=> create deployment package"
cd handler
../../scripts/build.sh

echo "=> update infrastructure"
cd ../terraform
terraform apply --auto-approve

echo "=> execute function"
curl -i $(terraform output -raw url)
echo ""

echo "=> wait for logs and show last log stream"
sleep 10
../../scripts/logs.sh $(terraform output -raw function_name)
