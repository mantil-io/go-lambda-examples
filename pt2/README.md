
cd handler
../../scripts/build.sh

cd ../terraform
terraform init

export AWS_REGION=eu-central-1
export AWS_PROFILE=org5

terraform apply

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
