package local_test

import (
	"os"
	"path/filepath"

	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/local"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManagerFactory", func() {
	var (
		factory task.ManagerFactory
	)

	Context("when an invalid outputDir is provided", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory("this directory does not exist", "empty-context")
		})
		It("Create() returns an error", func() {
			_, err := factory.Create()
			Expect(err).To(HaveOccurred())
		})
		It("Save() returns an error", func() {
			manager := createEmptyManager()
			Expect(factory.Save(manager)).To(HaveOccurred())
		})
	})

	Context("when an non-existent context is provided", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory("testdata", "non-existent-context")
		})
		AfterEach(func() {
			os.RemoveAll(filepath.Join("testdata", "non-existent-context"))
		})
		It("successfully creates an empty manager", func() {
			manager, err := factory.Create()
			Expect(err).NotTo(HaveOccurred())
			Expect(manager.Tasks()).To(BeEmpty())
			Expect(manager.Events()).To(BeEmpty())
		})
		It("successfully saves a manager", func() {
			manager := createEmptyManager()
			Expect(factory.Save(manager)).To(Succeed())
		})
	})

	Context("when the outputDir/context data is invalid", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory("testdata", "bad-context")
		})
		It("returns an error", func() {
			_, err := factory.Create()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the outputDir/context data is valid", func() {
		AfterEach(func() {
			os.RemoveAll(filepath.Join("testdata", "good-tmp-context"))
		})
		It("successfully creates the manager", func() {
			factory := local.NewManagerFactory("testdata", "good-context")
			manager, err := factory.Create()
			Expect(err).NotTo(HaveOccurred())
			tasks := manager.Tasks()
			Expect(tasks).To(HaveLen(3))
			events := manager.Events()
			Expect(events).To(HaveLen(6))
		})
		It("successfully saves the manager", func() {
			factory := local.NewManagerFactory("testdata", "good-tmp-context")
			manager := createEmptyManager()
			manager.Create("1")
			manager.Create("2")
			manager.Create("3")
			manager.SetPriority("1", 5)
			manager.SetState("2", task.StateRunning)
			manager.Note("3", "tuna")
			Expect(factory.Save(manager)).To(Succeed())
		})
	})

	Context("when a context is passed that has multiple path segments", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory("testdata", "this/has/multiple/path/segments")
		})
		It("errors when trying to create the manager", func() {
			_, err := factory.Create()
			Expect(err).To(HaveOccurred())
		})
	})
})

func createEmptyManager() task.Manager {
	factory := local.NewManagerFactory("testdata", "non-existent-context")
	manager, err := factory.Create()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return manager
}
