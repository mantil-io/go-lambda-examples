package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func v1Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

func v2Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("request path: %s body: %s", req.RawPath, req.Body)

	lc, _ := lambdacontext.FromContext(ctx)
	body := fmt.Sprintf("Hello from v2, %s", lc.InvokedFunctionArn)

	var rsp events.APIGatewayV2HTTPResponse
	rsp.StatusCode = http.StatusOK
	rsp.Body = body
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	rsp.Headers = headers
	return rsp, nil
}

func main() {
	lambda.Start(v2Handler)
	//lambda.Start(raw)
}

func raw(ctx context.Context, req map[string]interface{}) error {
	buf, _ := json.Marshal(req)
	fmt.Printf("%s", buf)
	return nil
}
