<!--
most trivial lambda runtime

What is Lambda runtime?

Whay is that ineteresting?

How it works?
- next is blocking call, your Lambd is frozen on the next call

https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html



salje i druge headere
https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html
najvazniji koji bi trebalo obraditi 
Lambda-Runtime-Deadline-Ms 


za vjezbu moze napisati hendlanje error a 
kada callback vrati error
-->


## Go Lambda Runtime example

This is example of most trivial Lambda function build in Go. It is not build around [AWS Lambda Go](https://github.com/aws/aws-lambda-go/tree/0462b0000e7468bdc8a9c456273c1551fab284aa) package which provides integration between function and Lambda execution environment. Here is just plain Go code which demonstrates how Lambda function, runtime and execution environment interacts. The goal of this example is to explain interaction between this parts. In pure Go code without using any libraries. 

First lets explain parts of the system. What is function, runtime and execution environment.
[Here](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html) is visualization of interaction between this parts.

Function is the code that we write when building Lambda function. It accepts JSON payload and returns response and optionally an error. Payload needs to be JSON encoded. Response Go []byte. So function signature in Go code is `func handler(payload []byte) ([]byte, error)`.

Runtime is wrapper around function code which connects it to the Lambda execution environment. When using [Lambda custom runtime](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-custom.html) we need to provide our own wrapper. Custom runtime is build around Amazon Linux 2 expects bootstrap executable in the /var/task folder. It will start that executable which then needs to connect to the execution environment.  

Execution environment is container which Lambda service starts. It provides API endpoints for interaction with runtime. Here we will use only Runtime API. There are also Extensions and Logs API endpoints. Runtime API is reachable at the `http://127.0.0.1:9001` address inside container. There are two methods next and response. Next is used to get invocation request. Response is for sending response after handling that request. URL-s are:
* next: http://127.0.0.1:9001/2018-06-01/runtime/invocation/next
* response: http://127.0.0.1:9001/2018-06-01/runtime/invocation/requestID/response
RequestID changes with every invocation. It is provided in http header on in next http response. 

Runtime works in the endless loop. It makes HTTP GET request on the next API endpoint. That request is blocked until Lambda is invoked. During that blocking phase function is frozen. Process is not running, any coroutines are frozen. We are not charged for the time while the runtime is waiting for the next response. When Lambda is invoked next HTTP finishes and in the body we get invocation payload. That HTTP response has useful [headers](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html). Bare minimum that runtime needs to handle is to read `Lambda-Runtime-Aws-Request-Id` header which is needed for making API response call. Runtime executes function with the payload gets response payload which uses to make HTTP POST to the API response endpoint. After that runtime enters into new loop cycle makes next request on which is again frozen until next invocation occurs.   

## Code walk-trough




