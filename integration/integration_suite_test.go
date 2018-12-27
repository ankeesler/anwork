package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankeesler/anwork/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	anworkBin string
	outputDir string

	runningOnTravis bool

	version = runner.Version

	runWithApi                 bool
	apiSession                 *gexec.Session
	testsRunning, waitingOnApi chan struct{}
	apiOut, apiErr             *gbytes.Buffer
)

func init() {
	_, runWithApi = os.LookupEnv("ANWORK_TEST_RUN_WITH_API")
}

func run(outBuf, errBuf *gbytes.Buffer, args ...string) {
	reallyRunWithStatus(2, 0, outBuf, errBuf, args...)
}

func runWithStatus(exitCode int, outBuf, errBuf *gbytes.Buffer, args ...string) {
	reallyRunWithStatus(2, exitCode, outBuf, errBuf, args...)
}

func reallyRunWithStatus(offset, exitCode int, outBuf, errBuf *gbytes.Buffer, args ...string) {
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
	ExpectWithOffset(offset, err).To(Succeed())

	timer := time.NewTimer(time.Second * 3)
	select {
	case <-s.Exited:
		fmt.Fprintln(GinkgoWriter, "[out]:", string(outBuf.Contents()))
		fmt.Fprintln(GinkgoWriter, "[err]:", string(errBuf.Contents()))
		ExpectWithOffset(offset, s.ExitCode()).To(Equal(exitCode))
	case <-timer.C:
		Fail(fmt.Sprintf("The session %+v failed to exit within 3 seconds", s))
	}
}

func getBuildHash() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.CombinedOutput()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("output: %s", string(out)))
	return strings.Trim(string(out), "\n")
}

func getBuildDate() string {
	return time.Now().Format(time.RFC3339)
}

func runOfficialBuildScript(hash, date string) string {
	cwd, err := os.Getwd()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	script := filepath.Join(cwd, "..", "ci", "build.sh")
	anworkBin := filepath.Join(cwd, "..", "anwork")
	out, err := exec.Command(script, "-h", hash, "-d", date, "-o", anworkBin).CombinedOutput()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("output: %s", string(out)))

	return anworkBin
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

			privateKey, publicKey, secret := generateAPICreds()

			cmd := exec.Command(apiBin)
			cmd.Env = []string{"PORT=12346"}
			cmd.Env = append(cmd.Env, fmt.Sprintf("ANWORK_API_PUBLIC_KEY=%s", publicKey))
			cmd.Env = append(cmd.Env, fmt.Sprintf("ANWORK_API_SECRET=%s", secret))

			if _, ok := os.LookupEnv("ANWORK_API_ADDRESS"); !ok {
				Expect(os.Setenv("ANWORK_API_ADDRESS", "127.0.0.1:12346")).To(Succeed())
				Expect(os.Setenv("ANWORK_API_PRIVATE_KEY", privateKey)).To(Succeed())
				Expect(os.Setenv("ANWORK_API_SECRET", secret)).To(Succeed())

				apiOut, apiErr = gbytes.NewBuffer(), gbytes.NewBuffer()
				apiSession, err = gexec.Start(cmd, apiOut, apiErr)
				Expect(err).ToNot(HaveOccurred())

				testsRunning = make(chan struct{})
				waitingOnApi = make(chan struct{})

				go func() {
					select {
					case <-apiSession.Exited:
						panic(fmt.Sprintf("API exited with exit code %d: stdout='%s', stderr='%s'", apiSession.ExitCode(),
							string(apiOut.Contents()), string(apiErr.Contents())))

					case <-testsRunning:
					}
					close(waitingOnApi)
				}()
			}
		}
	})
	AfterSuite(func() {
		if apiSession != nil {
			close(testsRunning)
			<-waitingOnApi
			apiSession.Kill()
			fmt.Fprintln(GinkgoWriter, "\nAPI OUT:", string(apiOut.Contents()))
			fmt.Fprintln(GinkgoWriter, "\nAPI ERR:", string(apiErr.Contents()))
		}

		Expect(os.RemoveAll(outputDir)).To(Succeed())

		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}

func generateAPICreds() (privateKey, publicKey, secret string) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKey = string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	Expect(err).NotTo(HaveOccurred())
	publicKey = string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	))

	r := make([]byte, 8)
	n, err := rand.Read(r)
	Expect(err).NotTo(HaveOccurred())
	Expect(n).To(Equal(8))
	secret = hex.EncodeToString(r)

	return
}
