package plex

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test1(t *testing.T) {
	m := Mux{}

	m.Get("/", func(res http.ResponseWriter, req *http.Request, params map[string]string) {
		res.Write([]byte("OK"))
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	body := w.Body.String()

	if body != "OK" {
		t.Error(body, "!= \"OK\"")
	}
}
