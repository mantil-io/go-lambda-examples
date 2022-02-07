package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const (
	envFunctionName  = "AWS_LAMBDA_FUNCTION_NAME"
	envMemorySize    = "AWS_LAMBDA_FUNCTION_MEMORY_SIZE"
	envLogStreamName = "AWS_LAMBDA_LOG_STREAM_NAME"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("request path: %s body: %s", req.RawPath, req.Body)

	// read information from environment variables
	// ref: https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime
	functionName := os.Getenv(envFunctionName)
	memorySize := os.Getenv(envMemorySize)
	log.Printf("max memory size: %s", memorySize)
	body := fmt.Sprintf("Hello from %s", functionName)
	// // show all environment vairables in log
	// for _, l := range os.Environ() {
	// 	log.Printf("%s", l)
	// }

	// use context to get execution deadline
	if deadline, ok := ctx.Deadline(); ok {
		log.Printf("execution deadline: %v max run duration: %v", deadline, deadline.Sub(time.Now()))
	}

	// get runtime request ID
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		log.Printf("aws request id: %s", lc.AwsRequestID)
	}

	// build response
	logStreamName := os.Getenv(envLogStreamName)
	rsp := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		// set some http response headers
		Headers: map[string]string{
			"LogStreamName": logStreamName,
		},
	}
	return rsp, nil
}

func main() {
	lambda.Start(handler)
}
