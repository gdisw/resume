package htmx

import (
	"encoding/json"
	"net/http"
)

const httpStatusUnset = 0

type eventBuilder struct {
	writer     http.ResponseWriter
	statusCode int
	triggers   map[triggerMode][]triggerEvent
}

func NewEventBuilder(writer http.ResponseWriter) *eventBuilder {
	return &eventBuilder{
		writer:   writer,
		triggers: make(map[triggerMode][]triggerEvent),
	}
}

func (b *eventBuilder) Flush() {
	for mode, triggers := range b.triggers {
		if len(triggers) == 1 &&
			triggers[0].Payload == nil {

			b.writer.Header().Set(mode.headerName(), triggers[0].Name)
			continue
		}

		eventsByName := make(map[string]any)
		for _, trigger := range triggers {
			eventsByName[trigger.Name] = trigger.Payload
		}

		raw, _ := json.Marshal(eventsByName)
		b.writer.Header().Set(mode.headerName(), string(raw))
	}

	if b.statusCode > httpStatusUnset {
		b.writer.WriteHeader(b.statusCode)
	}
}
