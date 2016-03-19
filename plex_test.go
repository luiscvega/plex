package plex

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test1(t *testing.T) {
	m := Mux{}

	m.Get("/", func(req *http.Request, params map[string]string) Response {
		return Response{200, []byte("OK")}
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)

	body := w.Body.String()
	statusCode := w.Code

	if body != "OK" {
		t.Error(body, "!= \"OK\"")
	}

	if statusCode != http.StatusOK {
		t.Error(statusCode, "!=", http.StatusOK)
	}
}
