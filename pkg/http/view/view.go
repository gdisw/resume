package view

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"syscall"

	"maps"
)

type ViewOption func(*View)

type View struct {
	template  *template.Template
	layout    string
	defaults  ViewData
	reference *templateReference
}

type templateReference struct {
	path string
	glob []string
}

type ViewData map[string]any

func (vd ViewData) PutAll(other ViewData) ViewData {
	maps.Copy(vd, other)

	return vd
}

func NewPage(layout, path string, glob ...string) *View {
	var t *template.Template
	t = template.Must(pageTemplate.Clone())
	for _, g := range glob {
		t = template.Must(t.ParseGlob(prefixTemplate(g)))
	}
	t = template.Must(t.ParseFiles(prefixTemplate(path)))

	view := &View{template: t, layout: layout}
	if refreshViewEnabled {
		view.reference = &templateReference{path: path, glob: glob}
	}

	return view
}

func NewPartialGlob(glob string) *View {
	view := &View{
		template: template.Must(
			template.Must(partialTemplate.Clone()).
				ParseGlob(prefixTemplate(glob)),
		),
	}

	if refreshViewEnabled {
		view.reference = &templateReference{glob: []string{glob}}
	}

	return view
}

func (v *View) WithDefaults(vd ViewData) *View {
	v.defaults = vd
	return v
}

func (v *View) Debug() *View {
	for _, t := range v.template.Templates() {
		fmt.Println(t.Name())
	}
	return v
}

func (v *View) Render(w http.ResponseWriter, ctx context.Context, data ViewData) error {
	return v.RenderTemplate(w, v.layout, data, WithContext(ctx))
}

func (v *View) RenderTemplate(w http.ResponseWriter, p string, in ViewData, options ...ViewOption) error {
	for _, option := range options {
		option(v)
	}

	data := make(ViewData)
	maps.Copy(data, v.defaults)
	maps.Copy(data, in)

	if v.reference != nil {
		if lastViewChange().After(lastLoadBase) {
			LoadBase(basePath)
		}

		if v.reference.path == "" {
			for _, glob := range v.reference.glob {
				v = NewPartialGlob(glob)
				break
			}
		} else {
			v = NewPage(v.layout, v.reference.path, v.reference.glob...)
		}
	}

	err := v.template.ExecuteTemplate(w, p, data)
	if err != nil &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		slog.Log(context.Background(), slog.LevelError, "error rendering template", "error", err)
	}

	return err
}

func WithContext(ctx context.Context) func(*View) {
	return func(v *View) {
		defaults := make(ViewData)
		maps.Copy(defaults, v.defaults)

		fromCtx := viewDataFromContext(ctx)
		maps.Copy(defaults, fromCtx)

		v.defaults = defaults
	}
}

func prefixTemplate(path string) string {
	return filepath.Join(basePath, "templates", "pages", path)
}
