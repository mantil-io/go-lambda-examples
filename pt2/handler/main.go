package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request path: %s body: %s", req.Path, req.Body)

	lc, _ := lambdacontext.FromContext(ctx)
	body := fmt.Sprintf("Hello from, %s", lc.InvokedFunctionArn)

	var rsp events.APIGatewayProxyResponse
	rsp.StatusCode = http.StatusOK
	rsp.Body = body
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	rsp.Headers = headers
	return rsp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
