package main

import (
	"encoding/json"
	"log"

	"github.com/mantil-io/go-lambda-examples/pkg/extension"
)

func main() {
	h := &handler{}
	extension.Run(h)
}

type handler struct{}

func (h *handler) Invoke(evt *extension.NextEventResponse) error {
	log.Printf("extension invoke %s", pp(evt))
	return nil
}

func (h *handler) Shutdown(evt *extension.NextEventResponse) error {
	log.Printf("extension shutdown %s", pp(evt))
	return nil
}

func pp(o interface{}) []byte {
	buf, _ := json.Marshal(o)
	return buf
}
