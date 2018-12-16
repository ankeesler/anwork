package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api2 Suite")
}
