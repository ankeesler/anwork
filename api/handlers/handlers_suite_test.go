package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

func handleGet(h http.Handler, target string) *httptest.ResponseRecorder {
	return handle(h, http.MethodGet, target, nil)
}

func handle(h http.Handler, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)
	rsp := httptest.NewRecorder()
	h.ServeHTTP(rsp, r)
	return rsp
}
