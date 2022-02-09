## A guide to AWS Lambda for Go developers  -  Part 2

In the first [part](https://github.com/mantil-io/go-lambda-examples/tree/master/guide) of this guide, we saw how to create a simple Lambda function in Go. Here we will expand that and make our function callable from the internet. We will integrate API Gateway with a Lambda function on the backend. When a client calls our function URL, API Gateway sends the request to the Lambda function and returns the function's response to the client.

For running example, you will need access to an AWS account. You are all set if you have already walked through the first part. If not, take a look into [aws credentials](https://github.com/mantil-io/go-lambda-examples/tree/master/guide#aws-credentials) chapter.

<!--
https://github.com/mantil-io/go-lambda-examples/tree/master/guide#view-lambda-function-logs

What is stage
What is integration
What is deployment

opisati kakvi su ovo primjeri

kod je ogoljen do esencije
kod je pripremljen za igranje
za pocetak experimentiranja 
za nekog tko bi htio zagnjuriti ali ne zna od kuda krenuti

radi se o primjerima, ne o clancima
clanak je tu da potkrijepi kod, ali mogao bi ga napisati i kao komentare u kodu

tako da nisam dobar fit onome sto vi zamisljate kao clanak
taj clanak bez koda ne vrijedi nista
on je potpora kodu

da bi pisao o tome moram vec imati nekog iskustva
ne pisaem o temama o kojima nema iskustva
osim u nekim slucajevima kada zelim nauciti, onda si postavljam pitanja sto bi me sve mogli pitati oni kojima to budem objasnjavao
zelim nauciti da bi mogao objasniti drugima

destiliram koncepte i kod dok ne ostane samo no sto je esencijalno 
micem accidental complexiti, zbog projekta, zbog library
-->

## Running example

Let's get something working. Then we will explore Go code and terraform configuration. 

Step into the *handler* folder from the folder where this readme file is located. It contains a simple Go Lambda function prepared for HTTP API Gateway integration. *build.sh* will create a Lambda function deployment package. It is explained in the first part. 

``` sh
cd handler
../../scripts/build.sh
```

After this, we will have *function.zip* in the *handler* folder.

For the Terraform, I suggest that you set global [plugin-cache](https://www.terraform.io/cli/config/config-file#provider-plugin-cache) folder. That will save you time and disk space if you work with different Terraform projects. With this configuration, you will reuse plugins between projects: 
``` sh
echo 'plugin_cache_dir="$HOME/.terraform.d/plugin-cache"' > $HOME/.terraformrc
mkdir -p $HOME/.terraform.d/plugin-cache
```

Now move to the terraform folder, where we will spend the rest of the time. *terraform init* will prepare plugins and download them into *plugin-cache*.

``` sh
cd ../terraform
terraform init
``` 

Execute Terraform configuration and create infrastructure: 
``` sh
terraform apply --auto-approve
```
after ~20 seconds expected output is something like:

``` sh
Apply complete! Resources: 10 added, 0 changed, 0 destroyed.

Outputs:

endpoint = "https://in2keb62qf.execute-api.eu-central-1.amazonaws.com"
function_arn = "arn:aws:lambda:eu-central-1:052548195718:function:api-example-handler"
function_name = "api-example-handler"
url = "https://in2keb62qf.execute-api.eu-central-1.amazonaws.com/handler"
```
*endpoint* is the location of our API Gateway, *url* is the location on which you can reach our Lambda function.
We can use this *terraform* to list some of these outputs whenever needed. For example, this will show Lambda function URL: 
``` sh
echo $(terraform output -raw url)
```

we can use that URL to execute the function:

``` sh
curl $(terraform output -raw url)
```

the expected output is something like:

``` sh
Hello from arn:aws:lambda:eu-central-1:052548195718:function:api-example-handler
```

Now we can play with changing the Go code and sending input and building response of the function. When you change the Go code, run the *build.sh* step, then *terraform apply* to update infrastructure. 

When you want to remove all created resources, run: 
``` sh
terraform destroy --auto-approve
```
You can, of course, return to the apply step and create them again. 

## Go handler code

Handler code shows how to extract useful information from various available sources and return a response to the caller. The caller makes an HTTP request to the API Gateway endpoint. API Gateway packs that request and pushes a payload for the function invocation.

We are using HTTP API Gateway [proxy payload](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html) format 2.0 integration type. When we make an HTTP request to our endpoint, for example:

``` sh
curl https://in2keb62qf.execute-api.eu-central-1.amazonaws.com/handler -d "request body"
```
<!--
[handler/main.go](handler/main.go) is a simple Lambda function. We are passing [handler](handler/main.go#56) to the lambda package. It will run our handler on each Lambda function invocation. In this case when function is invoked through HTTP API Gatweay integration we expect *APIGatewayV2HTTPRequest* in the request and we are using *APIGatewayV2HTTPResponse* for response. 

When we invoke our function through HTTP API Gateway with [proxy payload](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html) format 2.0 the payload with wich the function is invoked look like this:
-->

payload which API Gateway passes to our function looks like this:
``` json
{
    "version": "2.0",
    "routeKey": "ANY /handler",
    "requestContext": {
        "timeEpoch": 1644265482216,
        "time": "07/Feb/2022:20:24:42 +0000",
        "stage": "$default",
        "routeKey": "ANY /handler",
        "requestId": "NMDxogOhliAEJTg=",
        "http": {
            "userAgent": "curl/7.77.0",
            "sourceIp": "93.140.84.169",
            "protocol": "HTTP/1.1",
            "path": "/handler",
            "method": "POST"
        },
        "domainPrefix": "in2keb62qf",
        "domainName": "in2keb62qf.execute-api.eu-central-1.amazonaws.com",
        "apiId": "in2keb62qf",
        "accountId": "123456789012"
    },
    "rawQueryString": "",
    "rawPath": "/handler",
    "isBase64Encoded": true,
    "headers": {
        "x-forwarded-proto": "https",
        "x-forwarded-port": "443",
        "x-forwarded-for": "93.140.84.169",
        "x-amzn-trace-id": "Root=1-6201800a-351c97731d1f143b5094ee4c",
        "user-agent": "curl/7.77.0",
        "host": "in2keb62qf.execute-api.eu-central-1.amazonaws.com",
        "content-type": "application/x-www-form-urlencoded",
        "content-length": "12",
        "accept": "*/*"
    },
    "body": "cmVxdWVzdCBib2R5"
}
```

*aws/aws-lambda-go* package provides Go structs for unpacking this payload. For the request that is [APIGatewayV2HTTPRequest](https://github.com/aws/aws-lambda-go/blob/main/events/apigw.go#L51-L64) and the response that API Gateway expects is defined in [APIGatewayV2HTTPResponse](https://github.com/aws/aws-lambda-go/blob/main/events/apigw.go#L123-L130). We are using these two types in the signature of our [handler](handler/main.go#L27) function, and the *lambda* package will handle unmarshal of the request and marshaling of the response.

Code in handler shows how to get request body. But, first, we need to [decode](handler/main.go#L64-L73) it from base64. 

Then we show how to get information from the Lambda [environment](handler/main.go#L36). A full list of the environment variables can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime). 

The context provided to the function caries execution [deadline](handler/main.go#L42). The function should complete before the deadline. 

The runtime request information can be found in the [lambdacontext](handler/main.go#L47). 

In the end, we show how to create [response](handler/main.go#L53) for the API Gateway. Status code, body and headers will be returned to the caller who made an HTTP request to the API Gateway.  

## Terraform configuration

Terraform configuration consists of three files:

* *main.tf* defines input and output variables 
* *function.gf* defines lambda function with the supporting resources
* *api.tf* defines API Gateway

### function.tf

In *function.tf*, we first create IAM role and attach [AWSLambdaBasicExecutionRole](terraform/function.tf#L24), which gives function permission to upload logs to CloudWatch. Other common Lambda roles can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html). 

After that, we define the Cloudwatch [log group](terraform/function.tf#L30-L33) for the function. The function can create a log group on its own if it doesn't exist. We create it upfront to make it part of the terraform-managed resources. So it will be deleted by terraform on infrastructure destroy.

We use the deployment package we prepared in the *handler* folder for the building [function](terraform/function.tf#L36-L48). [source_code_hash](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function#source_code_hash) directive will trigger function code update whenever file hash changes.

### api.tf

Here we create an API Gateway. Then define a mapping between API Gateway and our function, route on which function will be exposed. With all that in place, function will be reachable on the internet. 

We start with the definition of the [API Gateway resource](terraform/api.tf#L3-L9). 

There are three flavors of API Gateway. The first one was REST API; it still has the most features. HTTP API overlaps with REST in many segments. It is a more 'modern' implementation. AWS [claims](https://aws.amazon.com/about-aws/whats-new/2019/12/amazon-api-gateway-offers-faster-cheaper-simpler-apis-using-http-apis-preview/) that "HTTP APIs are up to 71% cheaper compared to REST APIs". It is a little simpler than REST API. The last API Gateway flavor is WebSocket, enabling the bidirectional client to backend communication. I'll save that for some future example.

Here we are using HTTP API Gateway. Terraform resource type *aws_apigatewayv2_api* will create HTTP API Gateway type if protocol_type is HTTP. The other option for *protocol_type* is WEBSOCKET for creating WebSocket API Gateway. [CORS](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-cors.html) configuration enables browsers to access API while served from different domains. 

[CloudWatch log group](terraform/api.tf#L13-L16) is the place where Gateway access logs will be stored. */aws/vendedlogs* is a required prefix for services that create a vast number of [log groups](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AWS-logs-and-resource-policy.html).

API Gateway can have multiple [stages](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-stages.html) with different configurations (for example, dev beta prod...). Here we will use just the *\$default* stage. It is a reserved name for the stage which is served from the base of our API's URL. Stages and stage deployments can be powerful concepts but reserve them for some complicated scenarios. Until then, stick to the *\$default* stage and [automatic deployment](terraform/api.tf#L23).

In [*access_log_settings*](terraform/api.tf#L24-L38) we are configuring where to send access logs and how they will look.

Integration resource, [*aws_apigatewayv2_integration*](terraform/api.tf#L43-L49), is the place where we connect function and API Gateway. [*aws_apigatewayv2_route*](terraform/api.tf#L53-L57) sets path in HTTP request where function will be exposed. Route key "ANY /\${var.route}" when route [variable](terraform/main.tf#L16) is set to "handler" exposes function on /handler path for all types of HTPP request (GET, POST, ...).

In the end, we need to allow our API Gateway to [invoke function](terraform/api.tf#L61-L66). By default, in AWS, every resource is created without explicit permissions, so we need to set them for each resource to resource access. 

<!--
stages... ima ih vise $default automatic deployment
-->


## The path to the Mantil

<!--
znamo da je ovo komplicirano
all the code you write is only business logic


v1:

{
    "body": null,
    "headers": {
        "Content-Length": "0",
        "Host": "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com",
        "User-Agent": "curl/7.77.0",
        "X-Amzn-Trace-Id": "Root=1-61fff17e-5b8496e96e81bed7494730a7",
        "X-Forwarded-For": "93.136.72.29",
        "X-Forwarded-Port": "443",
        "X-Forwarded-Proto": "https",
        "accept": "*/*"
    },
    "httpMethod": "POST",
    "isBase64Encoded": false,
    "multiValueHeaders": {
        "Content-Length": [
            "0"
        ],
        "Host": [
            "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com"
        ],
        "User-Agent": [
            "curl/7.77.0"
        ],
        "X-Amzn-Trace-Id": [
            "Root=1-61fff17e-5b8496e96e81bed7494730a7"
        ],
        "X-Forwarded-For": [
            "93.136.72.29"
        ],
        "X-Forwarded-Port": [
            "443"
        ],
        "X-Forwarded-Proto": [
            "https"
        ],
        "accept": [
            "*/*"
        ]
    },
    "multiValueQueryStringParameters": null,
    "path": "/handler/",
    "pathParameters": {
        "proxy": ""
    },
    "queryStringParameters": null,
    "requestContext": {
        "accountId": "052548195718",
        "apiId": "9pyofn5yi9",
        "domainName": "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com",
        "domainPrefix": "9pyofn5yi9",
        "extendedRequestId": "NIKr1jmIFiAEJmw=",
        "httpMethod": "POST",
        "identity": {
            "accessKey": null,
            "accountId": null,
            "caller": null,
            "cognitoAmr": null,
            "cognitoAuthenticationProvider": null,
            "cognitoAuthenticationType": null,
            "cognitoIdentityId": null,
            "cognitoIdentityPoolId": null,
            "principalOrgId": null,
            "sourceIp": "93.136.72.29",
            "user": null,
            "userAgent": "curl/7.77.0",
            "userArn": null
        },
        "path": "/handler/",
        "protocol": "HTTP/1.1",
        "requestId": "NIKr1jmIFiAEJmw=",
        "requestTime": "06/Feb/2022:16:04:14 +0000",
        "requestTimeEpoch": 1644163454741,
        "resourceId": "ANY /handler/{proxy+}",
        "resourcePath": "/handler/{proxy+}",
        "stage": "$default"
    },
    "resource": "/handler/{proxy+}",
    "stageVariables": null,
    "version": "1.0"
}


v2: 

{
    "headers": {
        "accept": "*/*",
        "content-length": "0",
        "host": "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com",
        "user-agent": "curl/7.77.0",
        "x-amzn-trace-id": "Root=1-61fff2e3-6c0b0b63585d1a33199498a6",
        "x-forwarded-for": "93.136.72.29",
        "x-forwarded-port": "443",
        "x-forwarded-proto": "https"
    },
    "isBase64Encoded": false,
    "pathParameters": {
        "proxy": "pero"
    },
    "rawPath": "/handler",
    "rawQueryString": "",
    "requestContext": {
        "accountId": "123456789012",
        "apiId": "9pyofn5yi9",
        "domainName": "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com",
        "domainPrefix": "9pyofn5yi9",
        "http": {
            "method": "POST",
            "path": "/handler",
            "protocol": "HTTP/1.1",
            "sourceIp": "93.136.72.29",
            "userAgent": "curl/7.77.0"
        },
        "requestId": "NILjoi6FFiAEJ6A=",
        "routeKey": "ANY /handler/{proxy+}",
        "stage": "$default",
        "time": "06/Feb/2022:16:10:11 +0000",
        "timeEpoch": 1644163811849
    },
    "routeKey": "ANY /handler",
    "version": "2.0"
}


-->
