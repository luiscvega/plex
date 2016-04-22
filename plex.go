package plex

import (
	"net/http"
	"regexp"
)

type route struct {
	method  string
	pattern *regexp.Regexp
	keys    []string
	handler func(*http.Request, map[string]string) Response
}

type Mux struct {
	routes []route
}

var re = regexp.MustCompile(`:(\w+)`)

func (m *Mux) Add(method, path string, handler func(*http.Request, map[string]string) Response) {
	// Step 1: Set method
	r := route{method: method}

	// Step 2:
	result := re.FindAllStringSubmatch(path, -1)
	for _, tuple := range result {
		r.keys = append(r.keys, tuple[1])
	}

	// Step 3:
	r.pattern = regexp.MustCompile("^" + re.ReplaceAllLiteralString(path, `([^\\/]+)`) + "$")

	// Step 4:
	r.handler = handler

	// Step 5:
	m.routes = append(m.routes, r)
}

func (m *Mux) Get(path string, handler func(*http.Request, map[string]string) Response) {
	m.Add("GET", path, handler)
}

func (m *Mux) Post(path string, handler func(*http.Request, map[string]string) Response) {
	m.Add("POST", path, handler)
}

func (m *Mux) Put(path string, handler func(*http.Request, map[string]string) Response) {
	m.Add("PUT", path, handler)
}

func (m *Mux) Delete(path string, handler func(*http.Request, map[string]string) Response) {
	m.Add("DELETE", path, handler)
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range m.routes {
		// Step 1:
		if r.Method != route.method {
			continue
		}

		// Step 2:
		matches := route.pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			continue
		}

		// Step 3:
		params := map[string]string{}
		for i, value := range matches[1:] {
			key := route.keys[i]
			params[key] = value
		}

		// Step 4:
		response := route.handler(r, params)

		// Step 5:
		for k, vs := range response.Headers {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)

		// Step 6: Return since correct handler was found
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}
