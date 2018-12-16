package integration_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	api "github.com/ankeesler/anwork/api2"
	"github.com/ankeesler/anwork/api2/authenticator"
	"github.com/ankeesler/anwork/api2/client"
	"github.com/ankeesler/anwork/task2"
	"github.com/ankeesler/anwork/task2/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Repo", func() {
	var (
		dir string

		process ifrit.Process
	)

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "anwork-api-integration")
		Expect(err).NotTo(HaveOccurred())

		repo := fs.New(filepath.Join(dir, "test-context"))
		authenticator := authenticator.New()

		a := api.New(log.New(GinkgoWriter, "api-test: ", 0), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())

		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	task2.RunRepoTests(func() task2.Repo {
		return client.New("127.0.0.1:12345")
	})
})
