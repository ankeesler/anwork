package fs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FS Task Repo", func() {
	var dir, file string

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "fs-task-repo-test")
		file = filepath.Join(dir, "test-context")

		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	task.RunRepoTests(func() task.Repo {
		return fs.New(file)
	})

	Context("when file is invalid", func() {
		It("fails to run operations", func() {
			repo := fs.New("/this/file/totally/does/not/exist")
			Expect(repo.CreateTask(&task.Task{Name: "task-a"})).NotTo(Succeed())
		})
	})
})
