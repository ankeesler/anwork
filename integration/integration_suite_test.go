package integration

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	anworkBin string

	runningOnTravis bool

	version = 4
)

func run(outBuf, errBuf *gbytes.Buffer, args ...string) {
	runWithStatus(0, outBuf, errBuf, args...)
}

func runWithStatus(exitCode int, outBuf, errBuf *gbytes.Buffer, args ...string) {
	if outBuf == nil {
		outBuf = gbytes.NewBuffer()
	}
	if errBuf == nil {
		errBuf = gbytes.NewBuffer()
	}
	s, err := gexec.Start(exec.Command(anworkBin, args...), outBuf, errBuf)
	ExpectWithOffset(1, err).To(Succeed())
	EventuallyWithOffset(1, s).Should(gexec.Exit(exitCode), "STDOUT: %s\nSTDERR: %s\n",
		string(outBuf.Contents()), string(errBuf.Contents()))
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		anworkBin, err = gexec.Build("github.com/ankeesler/anwork/cmd/anwork")
		Expect(err).ToNot(HaveOccurred())

		_, runningOnTravis = os.LookupEnv("TRAVIS")

		Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "y")).To(Succeed())
	})
	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}
