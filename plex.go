package plex

import (
	"log"
	"net/http"
	"regexp"
)

type HandlerFunc func(*http.Request, map[string]string) Response

func (h HandlerFunc) Serve(r *http.Request, params map[string]string) Response {
	return h(r, params)
}

type Handler interface {
	Serve(*http.Request, map[string]string) Response
}

type Route struct {
	Method  string
	Pattern *regexp.Regexp
	Keys    []string
	Handler Handler
}

type Routes []Route

var re = regexp.MustCompile(`:(\w+)`)

func (rs *Routes) Add(method, path string, handler Handler) *Routes {
	// Step 1: Set method
	r := Route{Method: method}

	// Step 2:
	result := re.FindAllStringSubmatch(path, -1)
	for _, tuple := range result {
		r.Keys = append(r.Keys, tuple[1])
	}

	// Step 3:
	r.Pattern = regexp.MustCompile("^" + re.ReplaceAllLiteralString(path, `([^\\/]+)`) + "$")

	// Step 4:
	r.Handler = handler

	// Step 5:
	*rs = append(*rs, r)

	return rs
}

func (rs *Routes) Get(path string, handler Handler) {
	rs.Add("GET", path, handler)
}

func (rs *Routes) Post(path string, handler Handler) {
	rs.Add("POST", path, handler)
}

func (rs *Routes) Put(path string, handler Handler) {
	rs.Add("PUT", path, handler)
}

func (rs *Routes) Delete(path string, handler Handler) {
	rs.Add("DELETE", path, handler)
}

func (rs Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rs {
		// Step 1:
		if r.Method != route.Method {
			continue
		}

		// Step 2:
		matches := route.Pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			continue
		}

		// Step 3:
		params := map[string]string{}
		for i, value := range matches[1:] {
			key := route.Keys[i]
			params[key] = value
		}

		// Step 4:
		response := route.Handler.Serve(r, params)

		// Step 5:
		for k, vs := range response.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(response.StatusCode)
		_, err := w.Write(response.Body)
		if err != nil {
			log.Println("PLEX WRITE ERROR:", err)
		}

		// Step 6: Return since correct handler was found
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func Redirect(url string) Response {
	return Response{
		http.StatusFound,
		nil,
		http.Header{"Location": []string{url}},
	}
}
