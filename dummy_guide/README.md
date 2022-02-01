## Dummy guide to AWS Lambda for Go developers - Part 1

This multi-part tutorial series aims to give you a feeling of Lambda function programming in Go. It assumes that you don't have previous knowledge of AWS or Lambda. 

This is the first part where we will make the few small steps to start you in this new world. For the difference of many other examples, we will not use AWS Console or a higher-level serverless tool. We will create one dummy Lambda function and then execute it. We will do that in two ways: one using AWS command-line interface, and in the other, we will use Terraform to create AWS resources. It's interesting to see two different approaches. One is imperative (AWS CLI), where we specify each step, and the other is declarative (Terraform), where we define the desired end state of the infrastructure.

## Toolset

For those who are on macOS and are using Homebrew getting started required tools is a one-liner:
``` sh
brew bundle
```
in the root of this repo. Of course you first need clone this repo.

For other OS-es, you will need to install [Go](https://go.dev/doc/install), [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html), [jq](https://stedolan.github.io/jq/) and [terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli).

## AWS Credentials

You will need an [AWS account](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/) and [access keys](https://aws.amazon.com/premiumsupport/knowledge-center/create-access-key/) for a user in that account.


After you have the access key, set it as [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html) in the shell:
``` sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=us-east-1
```
Change these demo values with your access key, secret key and AWS
region closest to you. We will use Graviton2 (ARM) powered Lambda
functions, so you need to chose one of the regions where it is
[supported](https://aws.amazon.com/blogs/aws/aws-lambda-functions-powered-by-aws-graviton2-processor-run-your-functions-on-arm-and-get-up-to-34-better-price-performance/):

> US East (N. Virginia), US East (Ohio), US West (Oregon), Europe (Frankfurt), Europe (Ireland), EU (London), Asia Pacific (Mumbai), Asia Pacific (Singapore), Asia Pacific (Sydney), Asia Pacific (Tokyo).


To test connectivity to the AWS account with CLI run:
``` sh
aws sts get-caller-identity --output table --no-cli-pager
```
This should return Account, Arn and, UserId of the AWS user for which you set
access credentials. If this command succeeds, you are ready to go.

## About Go code 

In the handler folder is a dummy Lambda function handler. That is an unmodified copy
of the code from AWS
[docs](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html). Lambda
package provides all the plumbing with Lambda execution runtime. It uses [reflection](https://github.com/aws/aws-lambda-go/blob/2e104a66b60ac51aa6d7e494981203da7628426f/lambda/handler.go#L87) to analyze provided handler, performs JSON  [deserialization](https://github.com/aws/aws-lambda-go/blob/2e104a66b60ac51aa6d7e494981203da7628426f/lambda/handler.go#L115) of the payload and [serialization](https://github.com/aws/aws-lambda-go/blob/2e104a66b60ac51aa6d7e494981203da7628426f/lambda/handler.go#L29) of the response.

Implemented handler (HandleRequest function in the example) must satisfy these [rules](https://github.com/aws/aws-lambda-go/blob/0462b0000e7468bdc8a9c456273c1551fab284aa/lambda/entry.go#L16).


## AWS CLI

### Create Lambda function

Lets create first create a lambda function, and then we will look into the
process. Position yourself into the _handler_ folder and there run _publish.sh_
from the scripts folder:

``` sh
cd handler
../../scripts/publish.sh
```

The expected output is something like this:
``` sh
=> build
=> create deployment package
  adding: bootstrap (deflated 49%)
=> create new role
{
    "Role": {
        "Path": "/",
        "RoleName": "go-handler-example-role",
        "RoleId": "AROAQYPA52WDGY247IQCE",
        "Arn": "arn:aws:iam::052548195718:role/go-handler-example-role",
        "CreateDate": "2022-01-19T14:43:33+00:00",
        "AssumeRolePolicyDocument": {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    },
                    "Action": "sts:AssumeRole"
                }
            ]
        }
    }
}
=> create Lambda function
{
    "FunctionName": "go-handler-example",
    "FunctionArn": "arn:aws:lambda:eu-central-1:052548195718:function:go-handler-example",
    "Runtime": "provided.al2",
    "Role": "arn:aws:iam::052548195718:role/go-handler-example-role",
    "Handler": "provided",
    "CodeSize": 4087476,
    "Description": "",
    "Timeout": 3,
    "MemorySize": 128,
    "LastModified": "2022-01-19T14:43:43.568+0000",
    "CodeSha256": "PRgB5sSH1C+B9YrsAquFvpyWgSfHvwBaOK33564ZZ6k=",
    "Version": "$LATEST",
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "d7e38f6b-ff8a-4873-a259-b6350b149b3d",
    "State": "Pending",
    "StateReason": "The function is being created.",
    "StateReasonCode": "Creating",
    "PackageType": "Zip",
    "Architectures": [
        "arm64"
    ]
}
```

Let's look inside the script and expain what is happening. This is _publish.sh_ script:

``` shell
 1 #!/usr/bin/env bash -e
 2 
 3 # read function name from first argument or use default
 4 function_name="${1:-go-handler-example}"
 5 
 6 # get folder of the this script
 7 scripts=$(dirname "$0")
 8 
 9 # run build script
10 $scripts/build.sh ${@:2}
11 
12 # check if the function already exists
13 if $(aws lambda get-function --function-name $function_name > /dev/null 2>&1); then
14     echo "=> update existing function"
15     aws lambda update-function-code \
16         --no-cli-pager \
17         --function-name "$function_name" \
18         --zip-file fileb://function.zip
19 else
20     # create new function
21     $scripts/create_function.sh $function_name
22 fi
23 
24 # delete artifacts
25 rm function.zip bootstrap
```

I don't assume any previous knowledge of the shell scripting, so we will look
into this script line by line. The script accepts the Lambda function name as the first
parameter; if not supplied, _go-handler-example_ will be used as default. Line 4
is where this happens; it uses the first argument or default for setting variable
*function_name*. If you don't want default function name run the script like
`../../scripts/publish.sh my-lambda-function-name`.  
Line 7 gets the folder where _publish.sh_ is located. We will call the other
scripts (deploy.sh, create_function.sh) from this one so we grab that path and
store it into _scripts_ variable.  
In line 10, we call _build.sh_, which will prepare Lambda function
deployment package. Lets look into _build.sh_:

``` shell
1 #!/usr/bin/env bash -e
2 
3 echo "=> build"
4 GOOS=linux GOARCH=arm64 go build -o bootstrap
5 
6 echo "=> create deployment package"
7 zip function.zip bootstrap $@
```

Line 4 is go build command. We are building for Linux arm64 platform. Lambda
functions can be run on either Intel on AWS Graviton2 processors. Use Graviton
to get [lower price and better
performance](https://aws.amazon.com/blogs/aws/aws-lambda-functions-powered-by-aws-graviton2-processor-run-your-functions-on-arm-and-get-up-to-34-better-price-performance/)
unless some requirements pull you back.  
The resulting binary is named bootstrap. That is a requirement of the Lambda runtime
provided.al2, which we will use for building the Lambda function. That runtime is
a tiny Linux instance based on Amazon Linux 2; it will execute the bootstrap binary when started. Again, that bootstrap name is a requirement of the provided.al2 runtime.

Line 7 creates _function.zip_ file with bootstrap file in it. That zip file is
the Lambda deployment package accepted by CLI commands, AWS Console or any other
tool from which you can create the Lambda function. `$@` at the end of the zip
command is here to enable you to add any other files to the package. So you can, for
example, with `../../scripts/build.sh config.yml` add a config file to the
package. That file will be available when running the Lambda function in the same
folder as the binary.

After the build phase, we have _function.zip_ file in the _handler_ folder. So
let's return to the publish script.

Line 13 checks whether the Lambda function with the name from variable
*function_name* already exists. `aws lambda get-function --function-name
$function_name` returns function configuration, but here we are checking only
the result. Whether it was successful or not. If it was successful that function
already exists, and we will just update the function code. If not, we will call another
script *create_function.sh*. Let's examine the process of creating a new Lambda
function:

``` shell
 1 #!/usr/bin/env bash -e
 2 
 3 function_name="${1:-go-handler-example}"
 4 
 5 echo "=> create new role"
 6 role_name="$function_name-role"
 7 aws iam create-role \
 8     --role-name "$role_name" \
 9     --no-cli-pager \
10     --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}'
11 
12 # read role arn
13 role_arn=$(aws iam get-role --role-name "$role_name" | jq .Role.Arn -r)
14 aws iam attach-role-policy \
15     --no-cli-pager \
16     --role-name "$role_name" \
17     --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
18 
19 aws iam wait role-exists --role-name "$role_name"
20 
21 echo "=> create Lambda function"
22 # run with retries of few seconds to give time role to become visible
23 for i in 5 1 1 1 1 1; do
24     sleep "$i" # waiting for role to be available
25     aws lambda create-function \
26         --function-name "$function_name" \
27         --runtime provided.al2 \
28         --zip-file fileb://function.zip \
29         --role "$role_arn" \
30         --handler provided \
31         --architectures "arm64" \
32         --no-cli-pager && break
33 done
```

From line 5 to line 19 is the process of creating the [Lambda execution
role](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html).
We need that role in line 29 for actually creating the Lambda function. I will
skip details of the IAM, policy, role story... It is necessary to give your function
permission on other AWS resources, but that is a separate topic. Just note that in
line 17, we give our Lambda function AWSLambdaBasicExecutionRole, which provides that
function with permission to upload logs to Cloudwatch and nothing more.

Loop in lines 23, 24 and `&& break` part at the end of the create-function
command (line 25) gives us a few retries for *create-function* command.
When we create a new role, it is not immediately visible by new Lambda functions, so
*create-function* command can fail with an error: `An error occurred
(InvalidParameterValueException) when calling the CreateFunction operation: The
role defined for the function cannot be assumed by Lambda.` The script waits 5
seconds before the first try and then makes a few more tries each after 1 second.

Line 25 is the actual
[create-function](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/create-function.html)
AWS CLI command. We provide the function name (line 26), runtime on which
Lambda function will be build (line 27). This is a Go application, so we use
provided.al2 runtime. There are
[runtimes](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html)
for other languages. In line 28, we give our _function.zip_ as content for the
new Lambda. Other useful option is to instead of using local file to specify S3
location where the package is located. In line 31, we specify architecture for
the function (arm64 or x86_64).
 
### Invoke Lambda function

Run the _invoke.sh_ script from the _scripts_ folder:
``` sh
../../scripts/invoke.sh
```

The expected output is:
``` shell
{
    "StatusCode": 200,
    "ExecutedVersion": "$LATEST"
}
"Hello Foo!"
```
And the _invoke.sh_ is:
``` shell
 1 #!/usr/bin/env bash -e
 2 
 3 function_name="${1:-go-handler-example}"
 4 
 5 aws lambda invoke \
 6   --function-name "$function_name" \
 7   --no-cli-pager \
 8   --cli-binary-format raw-in-base64-out \
 9   --payload '{"name":"Foo"}' \
10   response.json && cat response.json
11 
12 rm response.json
```

Script uses [lambda
invoke](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/lambda/invoke.html)
CLI command. Line 6 specifies a function and line 9 JSON payload. In this case we
are sending `{"name":"Foo"}` JSON. This command writes a response to the file. So
we provide a file, show response content `cat response.json` and remove that file
at the end of the script.

You can play by changing payload attribute to get different results.

### View Lambda function logs

Run: 
``` shell
../../scripts/logs.sh
```

The expected output is something like:
``` shell
last stream name: 2022/01/20/[$LATEST]7b576140275c4b4d9aee7288717766c3
1642696161536   START RequestId: 65da5fff-aea9-4d67-8366-499d3942adf7 Version: $LATEST\n
1642696161537   END RequestId: 65da5fff-aea9-4d67-8366-499d3942adf7\n
1642696161537   REPORT RequestId: 65da5fff-aea9-4d67-8366-499d3942adf7\tDuration: 1.14 ms\tBilled Duration: 46 ms\tMemory Size: 128 MB\tMax Memory Used: 17 MB\tInit Duration: 44.06 ms\t\n
```

_logs.sh_ script:
``` shell
 1 #!/usr/bin/env bash -e
 2 
 3 function_name="${1:-go-handler-example}"
 4 
 5 # get the name of the last log stream
 6 stream_name=$(aws logs describe-log-streams --log-group-name /aws/lambda/$function_name | jq ".logStreams[].logStreamName" -r | tail -n 1)
 7 
 8 echo "last stream name: $stream_name"
 9 # show logs as table
10 aws logs get-log-events \
11     --log-group-name /aws/lambda/$function_name \
12     --log-stream-name "$stream_name" \
13     | jq ".events[] | [.timestamp, .message] | @tsv" -r
```

Here we are showing lambda function logs from the AWS Cloudwatch service. By
default, the Lambda function sends logs to the Cloudwatch service. Cloudwatch is
organized into log groups and log streams. Each lambda function gets a log group
named `/aws/lambda/[function-name]`. Into that group, each Lambda initialization
creates a new log stream. Function initialization happens on
first invoke after that execution environment lives for some time.  
This scripts finds the last stream name for our function log group and then lists
logs in that stream. Line 6 executes *describe-log-streams* which list all log
streams in the log group in JSON array. We use jq tool here to select only
*logStreamName* attribute, `tail -n 1` returns last line from the list of
all streams. Now when we have *stream_name* we can call get-log-events for that
stream in line 10. Again we use jq to reformat JSON into the table.

These logs show only Lambda execution environment stats. Put some
`log.Printf(...)` lines into the handler Go code, and you will find them into the
logs. Any output from the handler binary will be available in Cloudwatch logs.

### Cleanup

To remove Lambda function and other created resources (role and logs) in the AWS
account, run the cleanup script:

``` shell
../../scripts/cleanup.sh
```

## Terraform

### Create infrastructure

Again position yourself into the _handler_ folder. Use _build.sh_ to create Lambda
deployment package _function.zip_ there:
``` sh
../../scripts/build.sh
```

Then switch to the _terraform_ folder:
``` shell
cd ..
cd terraform
```

Be sure to have set `AWS_DEFAULT_REGION` environment variable before running terraform. For example:
``` shell
export AWS_REGION=eu-central-1
```

Then execute terrafrom init and apply commands:
``` sh
terraform init
terraform apply --auto-approve
```

### Execute function

``` sh
../../scripts/invoke.sh $(terraform output --raw function_name)
```

Here, the `$(terraform output --raw function_name)` part is to read *function_name*
from the terraform state.

### Explore terraform configuration

This guide is not intended to be a terraform manual. We will just explore terraform
configuration to get a sense of this declarative approach to building
infrastructure.

``` hcl
 1 terraform {
 2   backend "local" {
 3     path = "./.state/terraform.tfstate"
 4   }
 5 }
 6 
 7 variable "function_name" {
 8   type    = string
 9   default = "go-tf-handler-example"
10 }
11 
12 provider "aws" {}
13 
14 resource "aws_iam_role" "fn" {
15   name = "${var.function_name}-role"
16 
17   assume_role_policy = jsonencode({
18     Version = "2012-10-17"
19     Statement = [
20       {
21         Effect = "Allow"
22         Action = "sts:AssumeRole"
23         Principal = {
24           Service = "lambda.amazonaws.com"
25         }
26       }
27     ]
28   })
29 }
30 
31 resource "aws_lambda_function" "fn" {
32   role          = aws_iam_role.fn.arn
33   function_name = var.function_name
34   filename      = "../handler/function.zip"
35   runtime       = "provided.al2"
36   handler       = "bootstrap"
37   architectures = ["arm64"]
38 }
39 
40 output "function_name" {
41   value = var.function_name
42 }
43 
44 output "function_arn" {
45   value = aws_lambda_function.fn.arn
46 }
```

The first five lines define where terraform will save its state. The simplest
method is to use a local filesystem. State will be saved into the _.state_ folder
into _terraform_ folder (where main.tf is located).  
In 7-10 we define function name as variable. You can [change
that](https://www.terraform.io/language/values/variables#variables-on-the-command-line)
in apply for example: `terraform apply --var="function_name=my-function-name"
--auto-approve`.  
Lines 14-29 will create IAM role for the Lambda function. It is referenced in
line 32 when creating Lambda function. With this reference terraform knows that it
first needs to create role resource because lambda function resource depends on the
role.  
Lines 31-38 are actual function creation. Again we provide the same information like
in CLI; function_name, a zip file with Lambda deployment package (filename),
runtime on which function is based (runtime) and architecture. Handler is fixed
to the _bootstrap_ for provided.al2 runtime.  
Lines 40-46 define output variables. We can view them with `terraform output` command. 

### Cleanup

To remove all resources created in apply run: 
``` shell
terraform destroy --auto-approve
```

## Conclusion
Congratulations! Now you know how to build the AWS Lambda function in Go, both via AWS CLI and Terraform. Thanks for following this tutorial.
In the next part of this serverless series, I will break down how to set up an API Gateway and enable public access to created Lambda function.

