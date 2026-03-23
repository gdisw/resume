package view

import (
	"encoding/json"
	"html/template"
)

func JavaScript(value any) template.JS {
	raw, _ := json.Marshal(value)
	return template.JS(raw)
}
