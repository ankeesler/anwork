package local_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/local"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManagerFactory", func() {
	AfterEach(func() {
		os.RemoveAll(filepath.Join("testdata", "non-existent-context"))
	})

	task.RunFactoryTests(local.NewManagerFactory("testdata", "non-existent-context"))
})

var _ = Describe("ManagerFactory custom tests", func() {
	var (
		factory task.ManagerFactory

		outputDir string
	)

	BeforeEach(func() {
		var err error
		outputDir, err = ioutil.TempDir("", "anwork.task.local.test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(outputDir)).To(Succeed())
	})

	Context("when an invalid outputDir is provided", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory("this directory does not exist", "empty-context")
		})
		It("Create() returns an error", func() {
			_, err := factory.Create()
			Expect(err).To(HaveOccurred())
		})
		It("Save() returns an error", func() {
			manager := createEmptyManager(outputDir)
			Expect(factory.Save(manager)).To(HaveOccurred())
		})
	})

	Context("when an non-existent context is provided", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory(outputDir, "non-existent-context")
		})
		It("successfully creates an empty manager", func() {
			manager, err := factory.Create()
			Expect(err).NotTo(HaveOccurred())
			Expect(manager.Tasks()).To(BeEmpty())
			Expect(manager.Events()).To(BeEmpty())
		})
		It("successfully saves a manager", func() {
			manager := createEmptyManager(outputDir)
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

	Context("when a context is passed that has multiple path segments", func() {
		BeforeEach(func() {
			factory = local.NewManagerFactory(outputDir, "this/has/multiple/path/segments")
		})
		It("errors when trying to create the manager", func() {
			_, err := factory.Create()
			Expect(err).To(HaveOccurred())
		})
	})
})

func createEmptyManager(outputDir string) task.Manager {
	factory := local.NewManagerFactory(outputDir, "non-existent-context")
	manager, err := factory.Create()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return manager
}
