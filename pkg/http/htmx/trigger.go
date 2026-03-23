package htmx

import (
	"net/http"
)

const (
	triggerModeDefault = triggerMode(iota)
	triggerModeAfterSwap
	triggerModeAfterSettle
)

type triggerMode uint8

func (m triggerMode) headerName() string {
	switch m {
	case triggerModeDefault:
		return "HX-Trigger"
	case triggerModeAfterSwap:
		return "HX-Trigger-After-Swap"
	case triggerModeAfterSettle:
		return "HX-Trigger-After-Settle"
	}

	return ""
}

type triggerEvent struct {
	Name    string
	Mode    triggerMode
	Payload any
}

type triggerOption func(*triggerEvent)

func AfterSwap() triggerOption {
	return func(te *triggerEvent) {
		te.Mode = triggerModeAfterSwap
	}
}

func AfterSettle() triggerOption {
	return func(te *triggerEvent) {
		te.Mode = triggerModeAfterSettle
	}
}

func WithPayload(payload any) triggerOption {
	return func(te *triggerEvent) {
		te.Payload = payload
	}
}

func (b *eventBuilder) Trigger(event string, options ...triggerOption) *eventBuilder {
	te := triggerEvent{
		Name: event,
		Mode: triggerModeDefault,
	}

	for _, opt := range options {
		opt(&te)
	}

	b.triggers[te.Mode] = append(b.triggers[te.Mode], te)

	return b
}

func Trigger(w http.ResponseWriter, event string, options ...triggerOption) {
	NewEventBuilder(w).
		Trigger(event, options...).
		Flush()
}
