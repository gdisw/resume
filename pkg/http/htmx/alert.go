package htmx

import (
	"net/http"
)

const (
	AlertLevelDefault = alertLevel("primary")
	AlertLevelSuccess = alertLevel("success")
	AlertLevelWarning = alertLevel("warning")
	AlertLevelDanger  = alertLevel("danger")
)

type alertLevel string

func (l alertLevel) title() string {
	switch l {
	case AlertLevelDefault:
		return "New message"
	case AlertLevelSuccess:
		return "Operation successful"
	case AlertLevelWarning:
		return "A small problem occurred"
	case AlertLevelDanger:
		return "An unexpected error occurred"
	}

	return ""
}

func (l alertLevel) statusCode() int {
	switch l {
	case AlertLevelWarning:
		fallthrough
	case AlertLevelDanger:
		return http.StatusBadRequest
	default:
		return httpStatusUnset
	}
}

type alertEvent struct {
	level          alertLevel
	title          string
	message        string
	skipStatusCode bool
}

type alertOption func(*alertEvent)

func WithTitle(title string) alertOption {
	return func(fe *alertEvent) {
		fe.title = title
	}
}

func SkipStatusCode() alertOption {
	return func(fe *alertEvent) {
		fe.skipStatusCode = true
	}
}

func (b *eventBuilder) Alert(level alertLevel, message string, options ...alertOption) *eventBuilder {
	ft := alertEvent{
		level:   level,
		title:   level.title(),
		message: message,
	}

	for _, opt := range options {
		opt(&ft)
	}

	if !ft.skipStatusCode {
		b.statusCode = level.statusCode()
	}

	b.Trigger("alert", WithPayload(map[string]any{
		"Level":   ft.level,
		"Title":   ft.title,
		"Message": ft.message,
	}))

	return b
}

func Alert(w http.ResponseWriter, level alertLevel, message string, options ...alertOption) {
	NewEventBuilder(w).
		Alert(level, message, options...).
		Flush()
}
