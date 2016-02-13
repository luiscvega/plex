package plex

import (
	"net/http"
	"regexp"
)

type route struct {
	method  string
	pattern *regexp.Regexp
	keys    []string
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

type Mux struct {
	routes []route
}

func (m *Mux) Add(method, path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	// Step 1: Set method
	r := route{method: method}

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

func (m *Mux) Get(path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	m.Add("GET", path, handler)
}

func (m *Mux) Post(path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	m.Add("POST", path, handler)
}

func (m *Mux) Put(path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	m.Add("PUT", path, handler)
}

func (m *Mux) Delete(path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	m.Add("DELETE", path, handler)
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

		// Step 5:
		break
	}
}
