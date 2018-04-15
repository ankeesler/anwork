package local_test

import (
	task "github.com/ankeesler/anwork/tasknew"
	"github.com/ankeesler/anwork/tasknew/local"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Manager", func() {
	task.RunManagerTests(local.NewManagerFactory("testdata", "non-existent-context"))
})
