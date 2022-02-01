package logs

import (
	"encoding/json"

	"github.com/mantil-io/go-lambda-examples/pkg/extension"
)

type Handler interface {
	Init(*extension.Register) error
	Line([]byte) error
	BatchEnd()
	Shutdown(*extension.Event) error
}

func Run(handler Handler, events ...EventType) {
	a := &agent{events: events, handler: handler}
	extension.Run(a, extension.Shutdown)
}

type agent struct {
	events  []EventType
	cli     *Client
	handler Handler
}

func (h *agent) Init(evt *extension.Register) error {
	if err := h.handler.Init(evt); err != nil {
		return err
	}
	h.cli = NewClient(h.handler)
	h.cli.Init(evt.ExtensionID, h.events...)
	return nil
}

// not used, just to satisfy interface
func (h *agent) Invoke(evt *extension.Event) error {
	return nil
}

func (h *agent) Shutdown(evt *extension.Event) error {
	h.cli.Shutdown()
	return h.handler.Shutdown(evt)
}

func pp(o interface{}) []byte {
	buf, _ := json.Marshal(o)
	return buf
}
