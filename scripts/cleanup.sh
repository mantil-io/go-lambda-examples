#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"
role_name="$function_name-role"

echo "=> delete Lambda function"
aws lambda delete-function --function-name "$function_name"

echo "=> delete role"
aws iam detach-role-policy \
    --role-name "$role_name" \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
aws iam delete-role --role-name "$role_name"

echo "=> remove cloudwatch logs"
aws logs delete-log-group --log-group-name /aws/lambda/$function_name
