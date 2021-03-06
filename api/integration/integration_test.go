package integration_test

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/auth"
	"github.com/ankeesler/anwork/api/client"
	"github.com/ankeesler/anwork/api/client/cache"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Repo", func() {
	var (
		dir       string
		cacheFile string

		logger     lager.Logger
		privateKey *rsa.PrivateKey
		secret     []byte

		process ifrit.Process
	)

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "anwork-api-integration")
		Expect(err).NotTo(HaveOccurred())

		repo := fs.New(filepath.Join(dir, "test-context"))

		privateKey = generatePrivateKey()
		secret = generateSecret()
		auth := auth.NewServer(
			clock.NewClock(),
			rand.Reader,
			&privateKey.PublicKey,
			secret,
		)

		logger = lagertest.NewTestLogger("api")
		a := api.New(logger, repo, auth)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)

		cacheFile = filepath.Join(dir, "cache")
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())

		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	task.RunRepoTests(func() task.Repo {
		return client.New(
			logger,
			"127.0.0.1:12345",
			auth.NewClient(clock.NewClock(), privateKey, secret),
			cache.New(cacheFile),
		)
	})
})
