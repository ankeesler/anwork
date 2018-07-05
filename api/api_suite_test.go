package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	address = "localhost:12345"
)

func get(path string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s%s", address, path)
	return http.Get(url)
}

func handleGet(h http.Handler, target string) *httptest.ResponseRecorder {
	return handle(h, http.MethodGet, target, nil)
}

func handlePost(h http.Handler, target string, body io.Reader) *httptest.ResponseRecorder {
	return handle(h, http.MethodPost, target, body)
}

func handlePut(h http.Handler, target string, body io.Reader) *httptest.ResponseRecorder {
	return handle(h, http.MethodPut, target, body)
}

func handleDelete(h http.Handler, target string) *httptest.ResponseRecorder {
	return handle(h, http.MethodDelete, target, nil)
}

func handle(h http.Handler, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)
	rsp := httptest.NewRecorder()
	h.ServeHTTP(rsp, r)
	return rsp
}

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
