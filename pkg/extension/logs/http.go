package logs

// Client has the listener that receives the logs and the logger that handles the received logs
type Client struct {
	listener *httpListener
}

// NewClient returns an agent to listen and handle logs coming from Logs API for HTTP
// Make sure the agent is initialized by calling Init(agentId) before subscription for the Logs API.
func NewClient(cb Handler) *Client {
	return &Client{
		listener: newHttpListener(cb),
	}
}

// Init initializes the configuration for the Logs API and subscribes to the Logs API for HTTP
func (a Client) Init(extensionID string, events ...EventType) error {
	a.listener.Start()
	if len(events) == 0 {
		events = []EventType{Platform, Function} //, logsapi.Extension}
	}
	//reference: https://docs.aws.amazon.com/lambda/latest/dg/runtimes-logs-api.html#runtimes-logs-api-buffering
	bufferingCfg := BufferingCfg{
		MaxItems:  1000,
		MaxBytes:  262144,
		TimeoutMS: 25,
	}
	return Subscribe(events, bufferingCfg, a.listener.URL(), extensionID)
}

// Shutdown finalizes the logging and terminates the listener
func (a *Client) Shutdown() {
	a.listener.shutdown()
}
