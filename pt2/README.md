In the first part of this guide we saw how to create simple Lambda function in Go. Here we will exapand that and make our function callable from the internet. We will integrate API Gateway  with a Lambda function on the backend. When a client calls our API, API Gateway sends the request to the Lambda function and returns the function's response to the client.

There are three flawors of API Gateway. First one was REST API it still has most features, HTTP API overlaps with REST in many features. It is more 'modern' implementation. AWS [calims](https://aws.amazon.com/about-aws/whats-new/2019/12/amazon-api-gateway-offers-faster-cheaper-simpler-apis-using-http-apis-preview/) that "HTTP APIs are up to 71% cheaper compared to REST APIs". It is little simplier than REST API. In this example we will use HTTP API. The last API Gateway flawor is WebSocket which enables bidirectional clinet to backend communication. I'll save that for some future example.

For running example you will need access to an aws account. If you already walk through [first part](https://github.com/mantil-io/go-lambda-examples/tree/master/guide#readme) you are all set. If not take look into [aws credential](https://github.com/mantil-io/go-lambda-examples/tree/master/guide#aws-credentials) chapter.

<!--
https://github.com/mantil-io/go-lambda-examples/tree/master/guide#view-lambda-function-logs

What is stage
What is integration
What is deployment
-->

## Running example

Let's get something working. Than we will explore Go handler code and Terrafrom configuration. 

From the folder where this readme file is located step into handler folder. It contains simple Go Lambda function prepared for HTTP API Gateway integration. *build.sh* will prepare Lambda function deployment package. It is expalined in the first part. 

``` sh
cd handler
../../scripts/build.sh
```

After this we will have *function.zip* in the *handler* folder.


For the Terraform I suggest that you set global [plugin-cache](https://www.terraform.io/cli/config/config-file#provider-plugin-cache) folder. That will save you time and disk space if you are working with different Terraform projects. With this configuration you will reuse plugins between projects: 
``` sh
echo 'plugin_cache_dir="$HOME/.terraform.d/plugin-cache"' > $HOME/.terraformrc
mkdir -p $HOME/.terraform.d/plugin-cache
```

Now move to the terraform folder where we will spend rest of the time. *terrafrom init* will prepare plugins and download them into *plugin-cache*.

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
*endpoint* is location of our API Gateway, *url* is location on which you can reach our Lambda function.
We can use this *terraform* to list some of this ouputs whenever we need that. For example this will show Lambda function URL: 
``` sh
echo $(terraform output -raw url)
```

we can use that URL to execute function:

``` sh
curl $(terraform output -raw url)
```
expected output is something like:

``` sh
Hello from arn:aws:lambda:eu-central-1:052548195718:function:api-example-handler
```

Now we can play with changing Go code. Sending input and building response of the function. If you change Go code run the *build.sh** step then *terraform apply* to update infrastructure. 

When you want to remove all created resources run: 
``` sh
terraform destroy --auto-approve
```
You can of course return than to the apply step and create them again. 


v1:
``` json
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
```

v2: 
``` json
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
    "rawPath": "/handler/pero",
    "rawQueryString": "",
    "requestContext": {
        "accountId": "052548195718",
        "apiId": "9pyofn5yi9",
        "domainName": "9pyofn5yi9.execute-api.eu-central-1.amazonaws.com",
        "domainPrefix": "9pyofn5yi9",
        "http": {
            "method": "POST",
            "path": "/handler/pero",
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
    "routeKey": "ANY /handler/{proxy+}",
    "version": "2.0"
}
```

