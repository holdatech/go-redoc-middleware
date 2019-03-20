package redocmiddle // import "git.holda.io/golang/redoc-middleware"

import (
	"bytes"
	"net/http"
	"path"
	"text/template"
)

// Opts configures the Redoc middlewares
type Opts struct {
	// BasePath for the UI path, defaults to: /
	BasePath string
	// Path combines with BasePath for the full UI path, defaults to: docs
	Path string
	// SpecPath the path to find the spec for
	SpecPath string
	// RedocURL for the js that generates the redoc site, defaults to: https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js
	RedocURL string
	// Title for the documentation site, default to: API documentation
	Title string
}

// EnsureDefaults in case some options are missing
func (r *Opts) EnsureDefaults() {
	if r.BasePath == "" {
		r.BasePath = "/"
	}
	if r.Path == "" {
		r.Path = "docs"
	}
	if r.SpecPath == "" {
		r.SpecPath = "./swagger.json"
	}
	if r.RedocURL == "" {
		r.RedocURL = redocLatest
	}
	if r.Title == "" {
		r.Title = "API documentation"
	}
}

// NewHandler creates a middleware to serve a documentation site for a swagger spec.
// This allows for altering the spec before starting the http listener.
func NewHandler(opts Opts) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		opts.EnsureDefaults()

		pth := path.Join(opts.BasePath, opts.Path)
		tmpl := template.Must(template.New("redoc").Parse(redocTemplate))

		buf := bytes.NewBuffer(nil)
		_ = tmpl.Execute(buf, opts)
		b := buf.Bytes()

		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path == pth {
				rw.Header().Set("Content-Type", "text/html; charset=utf-8")
				rw.WriteHeader(http.StatusOK)

				_, _ = rw.Write(b)
				return
			}

			// serve spec
			if r.URL.Path == pth+"/swagger.json" {
				http.ServeFile(rw, r, opts.SpecPath)
				return
			}
			h.ServeHTTP(rw, r)
		})
	}
}

const (
	redocLatest   = "https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js"
	redocTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
    <!-- needed for adaptive design -->
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!--
    ReDoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='{{ .Path }}/swagger.json'></redoc>
    <script src="{{ .RedocURL }}"> </script>
  </body>
</html>
`
)
