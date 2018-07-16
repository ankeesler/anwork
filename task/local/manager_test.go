package local_test

import (
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/local"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	task.RunManagerTests(local.NewManagerFactory("testdata", "non-existent-context"))

	It("does not return nil when Tasks() is called immediately!", func() {
		mf := local.NewManagerFactory("testdata", "non-existent-context")
		m, err := mf.Create()
		Expect(err).NotTo(HaveOccurred())
		Expect(m.Tasks()).NotTo(BeNil())
	})
})
