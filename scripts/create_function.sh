#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"

echo "=> create new role"
role_name="$function_name-role"
aws iam create-role \
    --role-name "$role_name" \
    --no-cli-pager \
    --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}'

# read role arn
role_arn=$(aws iam get-role --role-name "$role_name" | jq .Role.Arn -r)
aws iam attach-role-policy \
    --no-cli-pager \
    --role-name "$role_name" \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

aws iam wait role-exists --role-name "$role_name"

echo "=> create Lambda function"
# run with retries of few seconds to give time role to become visible
for i in 5 1 1 1 1 1; do
    sleep "$i" # waiting for role to be available
    aws lambda create-function \
        --function-name "$function_name" \
        --runtime provided.al2 \
        --zip-file fileb://function.zip \
        --role "$role_arn" \
        --handler provided \
        --architectures "arm64" \
        --no-cli-pager && break
done
