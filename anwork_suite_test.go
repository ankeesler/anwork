package anwork_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAnwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Anwork Suite")
}
