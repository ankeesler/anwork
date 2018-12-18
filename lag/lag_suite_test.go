package lag_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLag(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lag Suite")
}
