package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		log.Printf("requestID: %s", lc.AwsRequestID)
	}
	for _, e := range os.Environ() {
		log.Printf("%s", e)
	}

	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		log.Printf("waiting for the SIGTERM ")
		s := <-sigs
		log.Printf("received signal %s", s)
		time.Sleep(500 * time.Millisecond)
		log.Printf("done")
	}()

	lambda.Start(HandleRequest)
}
