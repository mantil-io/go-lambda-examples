package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	requestIDHeader = "Lambda-Runtime-Aws-Request-Id"
	runtimeAPIEnv   = "AWS_LAMBDA_RUNTIME_API"
)

func main() {
	loop(handler)
}

func handler(req []byte) []byte {
	log.Printf("handler: %s", req)
	return req
}

func loop(callback func([]byte) []byte) {
	apiRoot := os.Getenv(runtimeAPIEnv)
	nextURL := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", apiRoot)

	for {
		req, header := next(nextURL)
		requestID := header.Get(requestIDHeader)

		rsp := callback(req)

		responseURL := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/%s/response", apiRoot, requestID)
		response(responseURL, rsp)
	}
}

func next(url string) ([]byte, http.Header) {
	rsp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rsp.StatusCode != http.StatusOK {
		log.Fatalf("%s unexpected status code %d", url, rsp.StatusCode)
	}
	rsp.Body.Close()
	return body, rsp.Header
}

func response(url string, data []byte) {
	rsp, err := http.Post(url, "", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	if rsp.StatusCode != http.StatusAccepted {
		log.Fatalf("%s unexpected status code %d", url, rsp.StatusCode)
	}
}
