package api_test

import (
	"log"
	"net/http"
	"os"

	api "github.com/ankeesler/anwork/api2"
	"github.com/ankeesler/anwork/api2/api2fakes"
	"github.com/ankeesler/anwork/task2"
	"github.com/ankeesler/anwork/task2/task2fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Tasks", func() {
	var (
		repo          *task2fakes.FakeRepo
		authenticator *api2fakes.FakeAuthenticator

		process ifrit.Process
	)

	BeforeEach(func() {
		repo = &task2fakes.FakeRepo{}
		authenticator = &api2fakes.FakeAuthenticator{}

		a := api.New(log.New(GinkgoWriter, "api-test: ", 0), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Describe("Get", func() {
		var tasks []*task2.Task
		BeforeEach(func() {
			tasks = []*task2.Task{
				&task2.Task{Name: "task-a", ID: 1},
				&task2.Task{Name: "task-b", ID: 2},
				&task2.Task{Name: "task-c", ID: 3},
			}
		})

		It("responds with the tasks that the repo returns", func() {
			rsp, err := get("/api/v1/tasks")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			assertTasks(rsp, tasks)
		})
	})
})
