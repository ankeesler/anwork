package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = XDescribe("TaskIDHandler", func() {
	var (
		manager *taskfakes.FakeManager

		logWriter *gbytes.Buffer

		handler http.Handler
	)

	BeforeEach(func() {
		factory := &taskfakes.FakeManagerFactory{}
		manager = &taskfakes.FakeManager{}
		factory.CreateReturnsOnCall(0, manager, nil)

		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", log.Ldate|log.Ltime|log.Lshortfile)

		handler = api.NewTaskIDHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/tasks/10")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/tasks.."))
		Eventually(logWriter).Should(gbytes.Say("Getting taskID 10"))
	})

	Context("when the last path segment is bunk", func() {
		It("logs an error", func() {
			handleGet(handler, "/api/v1/tasks/tuna")
			Eventually(logWriter).Should(gbytes.Say("Unable to parse last path segment"))
		})

		It("returns internal server error", func() {
			rsp := handleGet(handler, "/api/v1/tasks/tuna")
			Expect(rsp.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("GET", func() {
		var t *task.Task
		BeforeEach(func() {
			t = &task.Task{Name: "task-a", ID: 5}
			manager.FindByIDReturnsOnCall(0, t)
		})

		It("returns a JSON object representing the task", func() {
			rsp := handleGet(handler, "/api/v1/tasks/5")

			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedTaskJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())
			actualTaskJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualTaskJson).To(Equal(expectedTaskJson))
		})

		It("logs the response", func() {
			handleGet(handler, "/api/v1/tasks/5")

			expectedTaskJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning task: %s", expectedTaskJson)))
		})

		Context("when there is no task associated with the provided ID", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("returns a not found", func() {
				rsp := handleGet(handler, "/api/v1/tasks/10")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})

			It("logs that it was not able to find the task", func() {
				handleGet(handler, "/api/v1/tasks/10")
				Eventually(logWriter).Should(gbytes.Say("No task with ID 10"))
			})
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/tasks/5", nil)
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		BeforeEach(func() {
			task := &task.Task{Name: "task-a", ID: 5}
			manager.FindByIDReturnsOnCall(0, task)
		})

		It("only updates the state when the state parameter is set", func() {
			req := api.UpdateTaskRequest{State: task.StateRunning}
			reqBytes, err := json.Marshal(req)
			Expect(err).NotTo(HaveOccurred())
			reqBody := bytes.NewBuffer(reqBytes)

			rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(manager.SetStateCallCount()).To(Equal(1))
			Expect(manager.SetPriorityCallCount()).To(Equal(0))

			name, state := manager.SetStateArgsForCall(0)
			Expect(name).To(Equal("task-a"))
			Expect(state).To(Equal(task.StateRunning))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("only updates the priority when the priority parameter is set", func() {
			req := api.UpdateTaskRequest{Priority: 10}
			reqBytes, err := json.Marshal(req)
			Expect(err).NotTo(HaveOccurred())
			reqBody := bytes.NewBuffer(reqBytes)

			rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(manager.SetStateCallCount()).To(Equal(0))
			Expect(manager.SetPriorityCallCount()).To(Equal(1))

			name, priority := manager.SetPriorityArgsForCall(0)
			Expect(name).To(Equal("task-a"))
			Expect(priority).To(Equal(10))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("updates the state and priority when both parameters are set", func() {
			req := api.UpdateTaskRequest{State: task.StateRunning, Priority: 10}
			reqBytes, err := json.Marshal(req)
			Expect(err).NotTo(HaveOccurred())
			reqBody := bytes.NewBuffer(reqBytes)

			rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(manager.SetStateCallCount()).To(Equal(1))
			Expect(manager.SetPriorityCallCount()).To(Equal(1))

			name, state := manager.SetStateArgsForCall(0)
			Expect(name).To(Equal("task-a"))
			Expect(state).To(Equal(task.StateRunning))

			name, priority := manager.SetPriorityArgsForCall(0)
			Expect(name).To(Equal("task-a"))
			Expect(priority).To(Equal(10))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("logs a bunch of stuff", func() {
			req := api.UpdateTaskRequest{State: task.StateRunning, Priority: 10}
			reqBytes, err := json.Marshal(req)
			Expect(err).NotTo(HaveOccurred())
			reqBody := bytes.NewBuffer(reqBytes)

			handlePut(handler, "/api/v1/tasks/5", reqBody)

			Eventually(logWriter).Should(gbytes.Say("updating task task-a"))
			Eventually(logWriter).Should(gbytes.Say(fmt.Sprintf("handling request %s", string(reqBytes))))
			Eventually(logWriter).Should(gbytes.Say("set state Running"))
			Eventually(logWriter).Should(gbytes.Say("set priority 10"))
		})

		Context("when the state is invalid", func() {
			It("returns bad request", func() {
				req := api.UpdateTaskRequest{State: -1}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				handlePut(handler, "/api/v1/tasks/5", reqBody)

				Expect(manager.SetStateCallCount()).To(Equal(0))

				Eventually(logWriter).Should(gbytes.Say("updating task task-a"))
			})

			It("logs the error", func() {
				req := api.UpdateTaskRequest{State: -1}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				handlePut(handler, "/api/v1/tasks/5", reqBody)

				Expect(manager.SetStateCallCount()).To(Equal(0))

				Eventually(logWriter).Should(gbytes.Say("invalid state -1"))
			})
		})

		Context("when the task id is invalid", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("returns not found", func() {
				req := api.UpdateTaskRequest{State: 5}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})

			It("logs the error", func() {
				req := api.UpdateTaskRequest{State: 5}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				handlePut(handler, "/api/v1/tasks/5", reqBody)

				Eventually(logWriter).Should(gbytes.Say("invalid task id 5"))
			})
		})

		Context("when we fail to set the state", func() {
			BeforeEach(func() {
				manager.SetStateReturnsOnCall(0, errors.New("some state error"))
			})

			It("returns internal server error", func() {
				req := api.UpdateTaskRequest{State: task.StateRunning}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

				Expect(rsp.Code).To(Equal(http.StatusInternalServerError))
				Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

				expectedErrJson, err := json.Marshal(api.ErrorResponse{Message: "failed to set state: some state error"})
				Expect(err).NotTo(HaveOccurred())
				errJson, err := ioutil.ReadAll(rsp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(errJson).To(Equal(expectedErrJson))
			})

			It("logs the error", func() {
				req := api.UpdateTaskRequest{State: task.StateRunning}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				handlePut(handler, "/api/v1/tasks/5", reqBody)

				Eventually(logWriter).Should(gbytes.Say("failed to set state: some state error"))
			})
		})

		Context("when we fail to set the priority", func() {
			BeforeEach(func() {
				manager.SetPriorityReturnsOnCall(0, errors.New("some priority error"))
			})

			It("returns internal server error", func() {
				req := api.UpdateTaskRequest{Priority: 10}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				rsp := handlePut(handler, "/api/v1/tasks/5", reqBody)

				Expect(rsp.Code).To(Equal(http.StatusInternalServerError))
				Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

				expectedErrJson, err := json.Marshal(api.ErrorResponse{Message: "failed to set priority: some priority error"})
				Expect(err).NotTo(HaveOccurred())
				errJson, err := ioutil.ReadAll(rsp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(errJson).To(Equal(expectedErrJson))
			})

			It("logs the error", func() {
				req := api.UpdateTaskRequest{Priority: 10}
				reqBytes, err := json.Marshal(req)
				Expect(err).NotTo(HaveOccurred())
				reqBody := bytes.NewBuffer(reqBytes)

				handlePut(handler, "/api/v1/tasks/5", reqBody)

				Eventually(logWriter).Should(gbytes.Say("failed to set priority: some priority error"))
			})
		})
	})

	Describe("DELETE", func() {
		BeforeEach(func() {
			manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"})
			manager.DeleteReturnsOnCall(0, nil)
		})

		It("responds with no content", func() {
			rsp := handleDelete(handler, "/api/v1/tasks/5")

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))
			Expect(manager.DeleteCallCount()).To(Equal(1))
			Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("responds with not found", func() {
				rsp := handleDelete(handler, "/api/v1/tasks/5")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when the delete operation fails", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"})
				manager.DeleteReturnsOnCall(0, errors.New("some delete error"))
			})

			It("responds with a server error", func() {
				rsp := handleDelete(handler, "/api/v1/tasks/5")

				Expect(rsp.Code).To(Equal(http.StatusInternalServerError))
				Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

				expectedErrJson, err := json.Marshal(api.ErrorResponse{Message: "some delete error"})
				Expect(err).NotTo(HaveOccurred())
				errJson, err := ioutil.ReadAll(rsp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(errJson).To(Equal(expectedErrJson))
			})
		})
	})
})
