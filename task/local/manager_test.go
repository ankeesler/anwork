package local_test

import (
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/local"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Manager", func() {
	task.RunManagerTests(local.NewManagerFactory("testdata", "non-existent-context"))
})
