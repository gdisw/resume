package view

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/sprig/v3"
)

var (
	basePath           string
	funcMap            template.FuncMap
	pageTemplate       *template.Template
	partialTemplate    *template.Template
	refreshViewEnabled bool
	lastLoadBase       time.Time
)

func LoadBase(path string) {
	basePath = path

	funcMap = template.FuncMap{
		"parseDate":        parseDate,
		"generateDomId":    GenerateDomId,
		"getStaticConfig":  getStaticConfig,
		"jsEscapeString":   template.JSEscapeString,
		"base64URLEncode":  base64URLEncode,
		"safeHTML":         safeHTML,
		"route":            route,
		"formatInLocation": formatInLocation,
		"arr":              arr,
		"formatDuration":   formatDuration,
		"js":               JavaScript,
		"formatDateTime":   formatDateTime,
		"md":               markdown,
		"formatSize":       formatSize,
	}

	pageDeps := []string{
		"templates/helpers/csrf.partial.tmpl",
		"templates/layouts/app.layout.tmpl",
		"templates/layouts/base.layout.tmpl",
	}
	for i, name := range pageDeps {
		pageDeps[i] = filepath.Join(path, name)
	}

	pageTemplate = template.Must(
		template.New("app").
			Funcs(sprig.HtmlFuncMap()).
			Funcs(funcMap).
			ParseFiles(pageDeps...),
	)

	if partialDeps := []string{}; len(partialDeps) > 0 {
		for i, name := range partialDeps {
			partialDeps[i] = filepath.Join(path, name)
		}

		partialTemplate = template.Must(
			template.New("partial").
				Funcs(sprig.HtmlFuncMap()).
				Funcs(funcMap).
				ParseFiles(partialDeps...),
		)
	}

	lastLoadBase = time.Now()
}

func SetRefreshViewEnabled() {
	refreshViewEnabled = true
}

func lastViewChange() time.Time {
	lastModificationTime := time.Time{}
	filepath.Walk(filepath.Join(basePath, "templates"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if info.ModTime().After(lastModificationTime) {
			lastModificationTime = info.ModTime()
		}

		return nil
	})

	return lastModificationTime
}
