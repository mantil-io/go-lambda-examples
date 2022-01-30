// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package extension

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Register is the body of the response for /register
type Register struct {
	FunctionName    string `json:"functionName"`
	FunctionVersion string `json:"functionVersion"`
	Handler         string `json:"handler"`
}

// Event is the response for /event/next
type Event struct {
	EventType          EventType `json:"eventType"`
	DeadlineMs         int64     `json:"deadlineMs"`
	RequestID          string    `json:"requestId"`
	InvokedFunctionArn string    `json:"invokedFunctionArn"`
	Tracing            Tracing   `json:"tracing"`
}

func (r *Event) Deadline() time.Time {
	return time.UnixMilli(r.DeadlineMs)
}

func (r *Event) Timeout() time.Duration {
	return r.Deadline().Sub(time.Now())
}

// Tracing is part of the response for /event/next
type Tracing struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// StatusResponse is the body of the response for /init/error and /exit/error
type StatusResponse struct {
	Status string `json:"status"`
}

// EventType represents the type of events recieved from /event/next
type EventType string

const (
	Invoke   EventType = "INVOKE"   // Invoke is a lambda invoke
	Shutdown EventType = "SHUTDOWN" // Shutdown is a shutdown event for the environment

	extensionNameHeader      = "Lambda-Extension-Name"
	extensionIdentiferHeader = "Lambda-Extension-Identifier"
	extensionErrorType       = "Lambda-Extension-Function-Error-Type"
)

type Handler interface {
	Init(*Register) error
	Invoke(*Event) error
	Shutdown(*Event) error
}

func Run(handler Handler, events ...EventType) error {
	cli := NewClient()
	rr, err := cli.Register(events...)
	if err != nil {
		return err
	}
	if err := handler.Init(rr); err != nil {
		cli.InitError(context.TODO(), err.Error())
	}
	return cli.loop(handler)
}

func interputContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		cancel()
	}()
	return ctx
}

func (cli *Client) loop(handler Handler) error {
	ctx := interputContext()
	for {
		res, err := cli.NextEvent(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			return err
		}

		switch res.EventType {
		case Invoke:
			err := handler.Invoke(res)
			if err != nil {
				_, _ = cli.ExitError(ctx, err.Error())
				return err
			}
			continue
		case Shutdown:
			return handler.Shutdown(res)
		}
	}
}

// Client is a simple client for the Lambda Extensions API
type Client struct {
	baseURL       string
	httpClient    *http.Client
	extensionID   string
	extensionName string
}

// NewClient returns a Lambda Extensions API client
func NewClient() *Client {
	extensionName := filepath.Base(os.Args[0]) // extension name has to match the filename
	awsLambdaRuntimeAPI := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	baseURL := fmt.Sprintf("http://%s/2020-01-01/extension", awsLambdaRuntimeAPI)
	return &Client{
		extensionName: extensionName,
		baseURL:       baseURL,
		httpClient:    &http.Client{},
	}
}

func (c *Client) ExtensionID() string { return c.extensionID }
func (e *Client) Name() string        { return e.extensionName }

// Register will register the extension with the Extensions API
func (e *Client) Register(events ...EventType) (*Register, error) {
	const action = "/register"
	url := e.baseURL + action
	if len(events) == 0 {
		events = []EventType{Invoke, Shutdown}
	}
	reqBody, err := json.Marshal(map[string]interface{}{
		"events": events,
	})
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(extensionNameHeader, e.extensionName)
	httpRes, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed with status %s", action, httpRes.Status)
	}
	defer httpRes.Body.Close()
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	var res Register
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	e.extensionID = httpRes.Header.Get(extensionIdentiferHeader)
	return &res, nil
}

// NextEvent blocks while long polling for the next lambda invoke or shutdown
func (e *Client) NextEvent(ctx context.Context) (*Event, error) {
	const action = "/event/next"
	url := e.baseURL + action

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(extensionIdentiferHeader, e.extensionID)
	httpRes, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed with status %s", action, httpRes.Status)
	}
	defer httpRes.Body.Close()
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	res := Event{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// InitError reports an initialization error to the platform. Call it when you registered but failed to initialize
func (e *Client) InitError(ctx context.Context, errorType string) (*StatusResponse, error) {
	const action = "/init/error"
	url := e.baseURL + action

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(extensionIdentiferHeader, e.extensionID)
	httpReq.Header.Set(extensionErrorType, errorType)
	httpRes, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed with status %s", action, httpRes.Status)
	}
	defer httpRes.Body.Close()
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	res := StatusResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ExitError reports an error to the platform before exiting. Call it when you encounter an unexpected failure
func (e *Client) ExitError(ctx context.Context, errorType string) (*StatusResponse, error) {
	const action = "/exit/error"
	url := e.baseURL + action

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(extensionIdentiferHeader, e.extensionID)
	httpReq.Header.Set(extensionErrorType, errorType)
	httpRes, err := e.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s failed with status %s", action, httpRes.Status)
	}
	defer httpRes.Body.Close()
	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	res := StatusResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
