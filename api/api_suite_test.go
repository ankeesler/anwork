package api_test

import (
	"fmt"
	"net/http"
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

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
