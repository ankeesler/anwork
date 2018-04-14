package local_test

import (
	task "github.com/ankeesler/anwork/tasknew"
	"github.com/ankeesler/anwork/tasknew/local"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Journal", func() {
	task.RunJournalTests(local.NewJournal)
})
