package client_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

func makeLogger() *log.Logger {
	return log.New(GinkgoWriter, "client-test: ", 0)
}
