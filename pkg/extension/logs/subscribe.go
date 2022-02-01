package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	schemaVersion                         = "2021-03-18"
	lambdaAgentIdentifierHeaderKey string = "Lambda-Extension-Identifier"
	actionPath                            = "2020-08-15/logs"
)

// EventType represents the type of logs in Lambda
type EventType string

const (
	Platform  EventType = "platform"  // Platform is to receive logs emitted by the platform
	Function  EventType = "function"  // Function is to receive logs emitted by the function
	Extension EventType = "extension" // Extension is to receive logs emitted by the extension
)

// BufferingCfg is the configuration set for receiving logs from Logs API. Whichever of the conditions below is met first, the logs will be sent
type BufferingCfg struct {
	// MaxItems is the maximum number of events to be buffered in memory. (default: 10000, minimum: 1000, maximum: 10000)
	MaxItems uint32 `json:"maxItems"`
	// MaxBytes is the maximum size in bytes of the logs to be buffered in memory. (default: 262144, minimum: 262144, maximum: 1048576)
	MaxBytes uint32 `json:"maxBytes"`
	// TimeoutMS is the maximum time (in milliseconds) for a batch to be buffered. (default: 1000, minimum: 100, maximum: 30000)
	TimeoutMS uint32 `json:"timeoutMs"`
}

// SubscribeRequest is the request body that is sent to Logs API on subscribe
type SubscribeRequest struct {
	SchemaVersion string       `json:"schemaVersion"`
	EventTypes    []EventType  `json:"types"`
	BufferingCfg  BufferingCfg `json:"buffering"`
	Destination   Destination  `json:"destination"`
}

// Destination is the configuration for listeners who would like to receive logs with HTTP
type Destination struct {
	Protocol   string `json:"protocol"`
	URI        string `json:"URI"`
	HttpMethod string `json:"method"`
	Encoding   string `json:"encoding"`
}

// Subscribe calls the Logs API to subscribe for the log events.
func Subscribe(types []EventType, bufferingCfg BufferingCfg, destinationURL string, extensionID string) error {
	data, err := json.Marshal(
		&SubscribeRequest{
			SchemaVersion: schemaVersion,
			EventTypes:    types,
			BufferingCfg:  bufferingCfg,
			Destination: Destination{
				Protocol:   "HTTP",
				URI:        destinationURL,
				HttpMethod: "POST",
				Encoding:   "JSON",
			},
		})
	if err != nil {
		return err
	}

	statusCode, err := httpSubscribe(data, extensionID)
	if err != nil {
		return err
	}
	if statusCode == http.StatusAccepted {
		fmt.Println("WARNING!!! Logs API is not supported! Is this extension running in a local sandbox?")
	} else if statusCode != http.StatusOK {
		return fmt.Errorf("subscribe failed with status %d", statusCode)
	}
	return nil
}

func httpSubscribe(data []byte, extensionID string) (int, error) {
	awsLambdaRuntimeAPI := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	url := fmt.Sprintf("http://%s/2020-08-15/logs", awsLambdaRuntimeAPI)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	contentType := "application/json"
	req.Header.Set("Content-Type", contentType)
	req.Header.Set(lambdaAgentIdentifierHeaderKey, extensionID)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
