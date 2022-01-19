#!/usr/bin/env bash -e

function_name="${1:-go-handler-example}"

# get the name of the last log stream
stream_name=$(aws logs describe-log-streams --log-group-name /aws/lambda/$function_name | jq ".logStreams[].logStreamName" -r | tail -n 1)

echo "last stream name: $stream_name"
# show logs as table
aws logs get-log-events \
    --log-group-name /aws/lambda/$function_name \
    --log-stream-name "$stream_name" \
    | jq ".events[] | [.timestamp, .message] | @tsv" -r
