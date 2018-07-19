package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	anworkBin string
	outputDir string

	runningOnTravis bool

	version = 4

	runWithApi     bool
	apiSession     *gexec.Session
	apiOut, apiErr *gbytes.Buffer
)

func init() {
	_, runWithApi = os.LookupEnv("ANWORK_TEST_RUN_WITH_API")
}

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

	needOutput := true
	for _, a := range args {
		if a == "-o" {
			needOutput = false
		}
	}

	if needOutput {
		args = append([]string{"-o", outputDir}, args...)
	}

	fmt.Fprintln(GinkgoWriter, "\n[running]:", anworkBin, strings.Join(args, " "))
	s, err := gexec.Start(exec.Command(anworkBin, args...), outBuf, errBuf)
	ExpectWithOffset(2, err).To(Succeed())

	EventuallyWithOffset(2, s).Should(gexec.Exit(exitCode), "STDOUT: %s\nSTDERR: %s\n",
		string(outBuf.Contents()), string(errBuf.Contents()))
	fmt.Fprintln(GinkgoWriter, "[out]:", string(outBuf.Contents()))
	fmt.Fprintln(GinkgoWriter, "[err]:", string(errBuf.Contents()))
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		anworkBin, err = gexec.Build("github.com/ankeesler/anwork/cmd/anwork")
		Expect(err).ToNot(HaveOccurred())

		outputDir, err = ioutil.TempDir("", "anwork.integration.test")
		Expect(err).ToNot(HaveOccurred())

		_, runningOnTravis = os.LookupEnv("TRAVIS")

		Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "y")).To(Succeed())

		if runWithApi {
			var apiBin string
			apiBin, err = gexec.Build("github.com/ankeesler/anwork/cmd/service")
			Expect(err).ToNot(HaveOccurred())

			Expect(os.Setenv("ANWORK_API_ADDRESS", "127.0.0.1:12345")).To(Succeed())

			apiOut, apiErr = gbytes.NewBuffer(), gbytes.NewBuffer()
			apiSession, err = gexec.Start(exec.Command(apiBin), apiOut, apiErr)
			Expect(err).ToNot(HaveOccurred())
		}
	})
	AfterSuite(func() {
		if apiSession != nil {
			apiSession.Kill()
			fmt.Fprintln(GinkgoWriter, "\nAPI OUT:", string(apiOut.Contents()))
			fmt.Fprintln(GinkgoWriter, "\nAPI ERR:", string(apiErr.Contents()))
		}

		Expect(os.RemoveAll(outputDir)).To(Succeed())

		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}
