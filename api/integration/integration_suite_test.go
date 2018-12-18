package integration_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Integration Suite")
}

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return privateKey
}

func generateSecret() []byte {
	bytes := make([]byte, 32)
	n, err := rand.Read(bytes)
	ExpectWithOffset(1, n).To(Equal(32))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return bytes
}
