package local_test

import (
	"os"
	"path/filepath"

	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/local"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	AfterEach(func() {
		Expect(os.RemoveAll(filepath.Join("testdata", "non-existent-context"))).To(Succeed())
	})

	task.RunManagerTests(local.NewManagerFactory("testdata", "non-existent-context"))
})
