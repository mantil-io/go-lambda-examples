package main

import (
	"encoding/json"
	"log"

	"github.com/mantil-io/go-lambda-examples/pkg/extension"
	"github.com/mantil-io/go-lambda-examples/pkg/extension/logs"
)

func main() {
	h := &handler{}
	logs.Run(h)
}

type handler struct{}

func (h *handler) Init(evt *extension.Register) error {
	log.Printf("extension init %s", pp(evt))
	return nil
}

func (h *handler) Line(l []byte) error {
	log.Printf("logs_extension line %s", l)
	return nil
}

func (h *handler) BatchEnd() {
	log.Printf("logs_extension batch end")
}

func (h *handler) Shutdown(evt *extension.Event) error {
	log.Printf("extension shutdown %s", pp(evt))
	return nil
}

func pp(o interface{}) []byte {
	buf, _ := json.Marshal(o)
	return buf
}
