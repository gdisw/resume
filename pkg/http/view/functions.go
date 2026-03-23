package view

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gdisw/resume/pkg/identifier"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func parseDate(in string) time.Time {
	t, err := time.Parse(time.RFC3339, in)
	if err == nil {
		return t
	}

	t, _ = time.Parse("2006-01-02", in)
	return t
}

func GenerateDomId() string {
	return identifier.Mint("hdom")
}

func base64URLEncode(in string) string {
	return base64.
		URLEncoding.
		WithPadding(base64.NoPadding).
		EncodeToString([]byte(in))
}

func safeHTML(value string) template.HTML {
	return template.HTML(value)
}

func formatInLocation(t time.Time, layout string, l *time.Location) string {
	return t.In(l).Format(layout)
}

func arr(els ...any) []any {
	return els
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithXHTML()),
)

func markdown(input string) string {
	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return input
	}
	return buf.String()
}

func formatSize(size int64) string {
	return humanize.Bytes(uint64(size))
}
