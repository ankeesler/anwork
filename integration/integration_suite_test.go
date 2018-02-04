package integration

import (
	"io"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	anworkBin string
)

func run(outBuf, errBuf io.Writer, args ...string) {
	runWithStatus(0, outBuf, errBuf, args...)
}

func runWithStatus(exitCode int, outBuf, errBuf io.Writer, args ...string) {
	s, err := gexec.Start(exec.Command(anworkBin, args...), outBuf, errBuf)
	ExpectWithOffset(1, err).To(Succeed())
	EventuallyWithOffset(1, s).Should(gexec.Exit(exitCode))
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		anworkBin, err = gexec.Build("github.com/ankeesler/anwork/cmd/anwork")
		Expect(err).ToNot(HaveOccurred())
	})
	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}
