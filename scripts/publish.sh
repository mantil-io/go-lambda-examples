#!/usr/bin/env bash -e

# read function name from first argument or use default
function_name="${1:-go-handler-example}"

# get folder of the this script
scripts=$(dirname "$0")

# run build script
$scripts/build.sh ${@:2}

# check if the function already exists
if $(aws lambda get-function --function-name $function_name > /dev/null 2>&1); then
    echo "=> update existing function"
    aws lambda update-function-code \
        --no-cli-pager \
        --function-name "$function_name" \
        --zip-file fileb://function.zip
else
    # create new function
    $scripts/create_function.sh $function_name
fi

# delete artifacts
rm function.zip bootstrap
