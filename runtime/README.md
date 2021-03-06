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

<!--
This is example of most trivial Lambda build in Go. It is not build around [AWS Lambda Go](https://github.com/aws/aws-lambda-go/tree/0462b0000e7468bdc8a9c456273c1551fab284aa) package which provides integration between function and Lambda execution environment. Here is just few lines Go code which demonstrates how Lambda function, runtime and execution environment interacts. The goal of this example is to explain interaction between this parts. 

First lets describe parts of the system. What is function, runtime and execution environment.
[Here](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html) is visualization of connections between this parts.
-->
I will start with describing parts of the system; what is Lambda, execution environment, runtime and function. After that we will concentrate on the runtime part. 

First two are provided by the AWS. When build Lambda in Go we need to provide other two; runtime and function. I think that is worth exploring runtime to get deeper understanding of how things are functioning under the hood. So instead of using standard [AWS Lambda Go](https://github.com/aws/aws-lambda-go/tree/0462b0000e7468bdc8a9c456273c1551fab284aa) runtime package we will build our own from scratch.

**Lambda** is whole system. We are pushing our code to the Lambda, setting configuration, sending invocations. Around our code Lambda starts one or many **execution environments**. The execution environment provides a **runtime API** interface for getting invocation events and sending responses. 

**Runtime** is wrapper around our code which connects it to the execution environment. For Go we are starting from [custom runtime](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-custom.html). That is tiny Linux container. We need to put an executable into that container. That executable consist of the runtime and our function. Runtime is responsible for the integration with the execution environment runtime API. It is running the function's setup code, reading invocation events from the runtime API. The runtime passes the event data to the function, and posts the response from the function back to runtime API. Integration is established by calling HTTP API endpoint from runtime. That endpoint is injected into container by execution environment. The rest of the example will explore that integration between runtime API and runtime. 

**Function** is the code that we write when building Lambda. It accepts JSON invocation event payload, returns response and optionally an error. Payload needs to be JSON encoded. Response is of `[]byte` type. Function signature in Go is `func handler(payload []byte) ([]byte, error)`.

My whiteboard visualization of this parts:
![whiteboard](runtime.png)
 [Here](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-context.html) you can find one by AWS. 

<!--
When using [Lambda custom runtime](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-custom.html) we need to provide our own wrapper. Custom runtime is build by Amazon Linux 2, expects bootstrap executable in the /var/task folder. It will start that executable which then needs to connect to the execution environment. Conection is established by calling HTTP API endpoint from runtime.    

Runtime is wrapper around function which connects it to the Lambda execution environment. When using [Lambda custom runtime](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-custom.html) we need to provide our own wrapper. Custom runtime is build by Amazon Linux 2, expects bootstrap executable in the /var/task folder. It will start that executable which then needs to connect to the execution environment. Conection is established by calling HTTP API endpoint from runtime.   
-->

### Execution environment runtime API

Our runtime is communicating with execution environment runtime API by making HTTP requests to the runtime API endpoint. Runtime API is reachable at the `http://127.0.0.1:9001` address inside container. There are two methods in runtime API *next* and *response*. *Next* is used to get invocation event. *Response* is for sending result after handling invocation request. Inside execution environment those endpoints are reachable at this URL-s:

* next: http://127.0.0.1:9001/2018-06-01/runtime/invocation/next
* response: http://127.0.0.1:9001/2018-06-01/runtime/invocation/requestID/response

RequestID changes with every invocation. It is provided in HTTP header of the *next* response. 
 
<!--
Execution environment is container which Lambda service starts. It provides API endpoints for runtime. Here we will use only runtime API endpoint. There are also extensions and logs API endpoints. Runtime API is reachable at the `http://127.0.0.1:9001` address inside container. There are two methods in Runtime API *next* and *response*. *Next* is used to get invocation request. *Response* is for sending result after handling invocation request. Inside execution environment those endpoints are reachable at this URL-s:
* next: http://127.0.0.1:9001/2018-06-01/runtime/invocation/next
* response: http://127.0.0.1:9001/2018-06-01/runtime/invocation/requestID/response

RequestID changes with every invocation. It is provided in HTTP header of the *next* response. 
-->


Runtime works in the endless loop. It makes HTTP GET request on the *next* API endpoint. That HTTP request is blocked until Lambda is invoked. During that blocking phase whole execution environment is frozen. Process are not running, any goroutines are frozen. We are not charged for the time while the runtime is waiting for the *next* response. When Lambda is invoked *next* HTTP finishes and in the body we get invocation event.

*Next* HTTP response has useful [headers](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html). Bare minimum that runtime needs to read is `Lambda-Runtime-Aws-Request-Id` header which is needed for making API *response* call. Next one usefull is `Lambda-Runtime-Deadline-Ms` which deadline for the function to finish execution, it will be killed after that point.

Runtime executes function with the invocation event gets response payload which is used to make HTTP POST to the API *response* endpoint. After that API call it enters into new loop cycle; makes *next* request on which is again frozen until next Lambda invocation occurs.

## Running example

Position yourself into _runtime_ folder of this repo. To create new Lambda named 'go-runtime-example' run the *publish.sh* script: 
``` sh
../scripts/publish.sh go-runtime-example
```

We can invoke our new Lambda with _invoke.sh_:

``` sh
../scripts/invoke.sh go-runtime-example '"my payload"'
``` 
Payload has to be JSON encoded. That's the reason that we have those single and double quotes in the shell. If you want to send JSON object write something like `'{"name":"My name"}'` for the second argument of the *invoke.sh*.

The output is something like this:

``` 
{
    "StatusCode": 200,
    "ExecutedVersion": "$LATEST"
}
"my payload"
``` 
First we see invocation response headers and after that response payload. Our function echoes request into response so we got what we send.

To explore function logs run *logs.sh*:

``` sh
../scripts/logs.sh go-runtime-example
```

``` sh
last stream name: 2022/01/29/[$LATEST]6e009ed08b454a9aa504aa8b445cfd68
1643454793031	START RequestId: f90d095b-ab88-4577-80eb-17977106d24c Version: $LATEST\n
1643454793032	2022/01/29 11:13:13 handler: "my payload"\n
1643454793033	END RequestId: f90d095b-ab88-4577-80eb-17977106d24c\n
1643454793033	REPORT RequestId: f90d095b-ab88-4577-80eb-17977106d24c\tDuration: 1.22 ms\tBilled Duration: 37 ms\tMemory Size: 128 MB\tMax Memory Used: 15 MB\tInit Duration: 34.80 ms\t\n
``` 

First line is the name of the last found log stream in the CloudWatch Lambda logs. The other lines are actual invocation log lines. 


To remove all created resources from your AWS account run *cleanup.sh*:

``` sh
../scripts/cleanup.sh go-runtime-example
``` 

## Code walk-through

Let's explore Go code to see how is it working in practice. [main.go](main.go) in this project is all the code needed for building basic Lambda function and runtime. You can see from imports [L3-L10](main.go#L3-L10) that we are depending only on the code from Go standard library.

Like any another Go executable we have `func main` [L17-L19](main.go#L17-L19) which is really simple. It starts *runtime* and passes *function* to the runtime. I'm using here *runtime* and *function* as names for Go funcs to describe meaning. Runtime is glue between Lambda execution environment and our code. Runtime is same for all Lambdas we will write. Function is code specific to this Lambda. 

In this example *function* [L21-L24](main.go#L21-L24) is trivial; just returns what it receives. 

*Runtime* func [L26-L39](main.go#L26-L39) is the most interesting part. Shows how to do integration with Lambda runtime API. First we read environment variable `AWS_LAMBDA_RUNTIME_API` [L27](main.go#L27) in which is address of the runtime API HTTP endpoint. Currently value is `127.0.0.1:9001`. In _nextURL_ [L28](main.go#L28) we prepare address for the *next* request. And then we enter endless loop. 

This loop is the heart of the runtime. We first make get request to the *nextURL* [L31](main.go#L31). Func *next* [L41-L55](main.go#L41-L55) makes HTTP GET request, reads response body, and returns body and the HTTP headers. That call is **blocked and our code is frozen in [L42](main.go#L42) until Lambda is invoked**. 

This is critical part to understand. It is main difference from running Go (or any other) code in non-Lambda environment. When Lambda execution environment gets *next* request it freezes process until it has invocation to push into the runtime. When Lambda is invoked execution environment unfreezes process and responds to *next* request. 

After the invocation *next* completes and we have req and headers [L31](main.go#L31). From the headers we extract requestID [L32](main.go#L32) which is needed for making *response* API call. Runtime than executes handler function [L34](main.go#L34), gets rsp from handler and uses it for the body of HTTP POST to the API *response* endpoint [L57-L65](main.go#L57-L65). 

Then we go to the next loop cycle. Runtime makes *next* request and blocks there until next invocation. 

## Extending example runtime

Here are some ideas of how to extend this simple example to be more complete Lambda runtime. 

Runtime API has endpoint where it accepts execution error. Runtime can finish the loop by calling `/runtime/invocation/requestId/response` or `/runtime/invocation/requestId/error`. We should allow our *function* to return error and if that is the case pass that to the *error* endpoint. On [this](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html) page under "Invocation error" is described JSON object for packing error data.

On the same [page](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-api.html) under "Next invocation" are described all headers which we get in the *next* API response. Currently we are using only one but more complete example should allow function to get values of that headers and call the function with context which has deadline from `Lambda-Runtime-Deadline-Ms` header.

If we add some code before runtime func call in the main that will be executed only one during the cold start of the execution environment. That is place for making initializations for the things used in function.
