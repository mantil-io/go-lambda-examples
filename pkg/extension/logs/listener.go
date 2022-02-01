package logs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
)

// DefaultHttpListenerPort is used to set the URL where the logs will be sent by Logs API
const DefaultHttpListenerPort = "3000"

// httpListener is used to listen to the Logs API using HTTP
type httpListener struct {
	httpServer *http.Server
	address    string
	cb         Handler
}

// newHttpListener returns a LogsApiHttpListener with the given log queue
func newHttpListener(cb Handler) *httpListener {
	return &httpListener{
		httpServer: nil,
		cb:         cb,
		address:    "sandbox:" + DefaultHttpListenerPort,
	}
}

func (l *httpListener) URL() string {
	return fmt.Sprintf("http://%s", l.address)
}

// Start initiates the server in a goroutine where the logs will be sent
func (l *httpListener) Start() {
	l.httpServer = &http.Server{Addr: l.address}
	http.HandleFunc("/", l.handler)
	go func() {
		//logger.Infof("Serving agent on %s", address)
		err := l.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("error unexpected stop on Http Server: %v", err)
			l.shutdown()
		}
	}()
}

// handler handles the requests coming from the Logs API.
// Everytime Logs API sends logs, this function will read the logs from the response body
// and put them into a synchronous queue to be read by the main goroutine.
// Logging or printing besides the error cases below is not recommended if you have subscribed to receive extension logs.
// Otherwise, logging here will cause Logs API to send new logs for the printed lines which will create an infinite loop.
func (l *httpListener) handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %+v", err)
		return
	}
	defer r.Body.Close()
	jsonparser.ArrayEach(body, func(line []byte, dataType jsonparser.ValueType, offset int, err error) {
		if line == nil {
			return
		}
		if err := l.cb.Line(line); err != nil {
			log.Printf("failed to send log line: %s", line)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	l.cb.BatchEnd()
}

// shutdown terminates the HTTP server listening for logs
func (l *httpListener) shutdown() {
	if l.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := l.httpServer.Shutdown(ctx)
		if err != nil {
			log.Printf("failed to shutdown http server gracefully %s", err)
		} else {
			l.httpServer = nil
		}
	}
}
