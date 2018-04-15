package runner_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAnworkrunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Anworkrunner Suite")
}
