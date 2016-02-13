package plex

import (
	"net/http"
	"regexp"
)

type Handler func(http.ResponseWriter, *http.Request, map[string]string)

type Route struct {
	method  string
	pattern *regexp.Regexp
	keys    []string
	handler Handler
}

type Mux struct {
	routes []Route
}

func (m *Mux) Add(method, path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	// Step 1: Set method
	r := Route{method: method}

	// Step 2:
	re := regexp.MustCompile(`:(\w+)`)

	result := re.FindAllStringSubmatch(path, -1)
	for _, tuple := range result {
		r.keys = append(r.keys, tuple[1])
	}

	// Step 3:
	r.pattern = regexp.MustCompile(re.ReplaceAllLiteralString(path, `([^\\/]+)`))

	// Step 4:
	r.handler = handler

	// Step 5:
	m.routes = append(m.routes, r)
}

func (m Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, r := range m.routes {
		// Step 1:
		if req.Method != r.method {
			continue
		}

		// Step 2:
		matches := r.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) == 0 {
			continue
		}

		// Step 3:
		params := map[string]string{}
		for i, value := range matches[1:] {
			key := r.keys[i]
			params[key] = value
		}

		// Step 4:
		r.handler(res, req, params)
	}
}
