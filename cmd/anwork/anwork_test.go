package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ankeesler/anwork/cmd/anwork/command"
	"github.com/ankeesler/anwork/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	root                = "test-data"
	context             = "test-context"
	journalPrefixRegexp = "\\[.*\\]"
)

var _ = Describe("anwork", func() {
	var output *bytes.Buffer
	var ret int
	callRun := func(args ...string) {
		output = new(bytes.Buffer)
		ret = run(append([]string{"anwork", "-context", context, "-root", root}, args...), output)
	}

	expectSuccess := func() {
		Expect(ret).To(Equal(0))
	}
	expectFailure := func() {
		Expect(ret).ToNot(Equal(0))
	}

	expectUsagePrinted := func() {
		Expect(output.String()).To(ContainSubstring("Usage of anwork"))
		for _, c := range command.Commands {
			Expect(output.String()).To(ContainSubstring(c.Name))
			Expect(output.String()).To(ContainSubstring(c.Description))
			for _, a := range c.Args {
				Expect(output.String()).To(ContainSubstring(a))
			}
		}
	}

	AfterEach(func() {
		p := storage.Persister{Root: root}
		Expect(p.Delete(context)).To(Succeed())
	})

	Context("when no args are passed", func() {
		BeforeEach(func() {
			callRun()
		})
		It("succeeds", expectSuccess)
		It("prints usage", expectUsagePrinted)
		It("does not print error", func() {
			Expect(output.String()).ToNot(ContainSubstring("Error!"))
		})
	})
	Context("when help is requested", func() {
		BeforeEach(func() {
			callRun("-h")
		})
		It("succeeds", expectSuccess)
		It("prints usage", expectUsagePrinted)
		It("prints usage only once!", func() {
			firstIndex := strings.Index(output.String(), "Usage of anwork")
			Expect(firstIndex).ToNot(Equal(-1))

			secondIndex := strings.Index(output.String()[firstIndex+1:], "Usage of anwork")
			Expect(secondIndex).To(Equal(-1))
		})
	})
	Context("when a bad flag is passed", func() {
		BeforeEach(func() {
			callRun("-tuna")
		})
		It("fails", expectFailure)
		It("prints usage", expectUsagePrinted)
	})
	Context("when a bad command is passed", func() {
		BeforeEach(func() {
			callRun("fish")
		})
		It("fails", expectFailure)
		It("prints usage", expectUsagePrinted)
		It("prints something about the unknown command", func() {
			Expect(output.String()).To(ContainSubstring("Error! Unknown command: fish"))
		})
	})
	Context("when a bad context is passed", func() {
		BeforeEach(func() {
			output = new(bytes.Buffer)
			ret = run([]string{"anwork", "-context", "/i/really/hope/this/file/does/not/exist", "show"},
				output)
		})
		It("fails", expectFailure)
	})
	Context("when the context is corrupt", func() {
		BeforeEach(func() {
			output = new(bytes.Buffer)
			ret = run([]string{"anwork", "-context", "bad-context", "-root", "test-data", "show"}, output)
		})
		It("fails", expectFailure)
		It("prints something about the context being bad", func() {
			Expect(output.String()).To(ContainSubstring("Could not read manager from file"))
		})
	})

	Context("when the version command is passed", func() {
		BeforeEach(func() {
			callRun("version")
		})
		It("succeeds", expectSuccess)
		It("prints the version", func() {
			msg := fmt.Sprintf("ANWORK Version = %d", command.Version)
			Expect(output.String()).To(ContainSubstring(msg))
		})
	})

	Context("when a task is created", func() {
		BeforeEach(func() {
			callRun("create", "task-a")
		})
		It("succeeds", expectSuccess)
		Context("when show is called", func() {
			BeforeEach(func() {
				callRun("show")
			})
			It("succeeds", expectSuccess)
			It("shows the created task", func() {
				Expect(output.String()).To(ContainSubstring("WAITING tasks:\n  task-a ("))
			})
		})
	})

	Context("when multiple tasks are created", func() {
		BeforeEach(func() {
			callRun("create", "task-a")
			callRun("create", "task-b")
			callRun("create", "task-c")
		})
		It("succeeds", expectSuccess)
		Context("when show is called", func() {
			BeforeEach(func() {
				callRun("show")
			})
			It("shows the created tasks in order of creation", func() {
				regexp := "WAITING tasks:\n  task-a.*\n  task-b.*\n  task-c.*"
				Expect(output.String()).To(MatchRegexp(regexp))
			})
		})
		Context("when journal is called", func() {
			BeforeEach(func() {
				callRun("journal")
			})
			It("shows the 3 creation entries", func() {
				regexp := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s", journalPrefixRegexp, ".*Created.*task-c.*",
					journalPrefixRegexp, ".*Created.*task-b.*", journalPrefixRegexp, ".*Created.*task-a.*")
				Expect(output.String()).To(MatchRegexp(regexp))
			})
		})
		Context("when the priority is changed on the tasks", func() {
			BeforeEach(func() {
				callRun("set-priority", "task-a", "20")
				callRun("set-priority", "task-b", "25")
				callRun("set-priority", "task-c", "15")
			})
			It("succeeds", expectSuccess)
			Context("when show is called", func() {
				BeforeEach(func() {
					callRun("show")
				})
				It("shows the tasks in order of priority", func() {
					regexp := "WAITING tasks:\n  task-c.*\n  task-a.*\n  task-b.*"
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
			Context("when journal is called", func() {
				BeforeEach(func() {
					callRun("journal")
				})
				It("shows the 3 creation entries plus the 3 priority updates in reverse order", func() {
					regexp := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s",
						journalPrefixRegexp, ".*priority.*task-c.*", journalPrefixRegexp, ".*priority.*task-b.*",
						journalPrefixRegexp, ".*priority.*task-a.*", journalPrefixRegexp, ".*Created.*task-c.*",
						journalPrefixRegexp, ".*Created.*task-b.*", journalPrefixRegexp, ".*Created.*task-a.*")
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
		})
		Context("when the state is set on the tasks", func() {
			BeforeEach(func() {
				callRun("set-running", "task-a")
				callRun("set-blocked", "task-b")
				callRun("set-blocked", "task-a")
				callRun("set-running", "task-c")
				callRun("set-finished", "task-c")
				callRun("set-running", "task-a")
				// At the end of these calls, task-a is running, task-b is blocked, and task-c is finished.
			})
			It("succeeds", expectSuccess)
			Context("when show is called", func() {
				BeforeEach(func() {
					callRun("show")
				})
				It("shows the tasks in the correct state section", func() {
					regexp := fmt.Sprintf("%s\n%s\n%s\n%s", "RUNNING tasks:\n  task-a.*",
						"BLOCKED tasks:\n  task-b.*", "WAITING tasks:", "FINISHED tasks:\n  task-c.*")
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
			Context("when journal is called", func() {
				BeforeEach(func() {
					callRun("journal")
				})
				It("shows the 3 creation entries plus the 6 priority updates in reverse order", func() {
					Expect(strings.Count(output.String(), "\n")).To(Equal(9))
				})
			})
		})
		Context("when one task is deleted", func() {
			BeforeEach(func() {
				callRun("delete", "task-b")
			})
			It("succeeds", expectSuccess)
			Context("when show is called", func() {
				BeforeEach(func() {
					callRun("show")
				})
				It("only shows the 2 remaining tasks", func() {
					regexp := "WAITING tasks:\n  task-a.*\n  task-c.*"
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
			Context("when journal is called", func() {
				BeforeEach(func() {
					callRun("journal")
				})
				It("shows the 3 creation entries plus the 1 deletion entry in reverse order", func() {
					regexp := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s", journalPrefixRegexp,
						".*Deleted.*task-b", journalPrefixRegexp, ".*Created.*task-c.*", journalPrefixRegexp,
						".*Created.*task-b.*", journalPrefixRegexp, ".*Created.*task-a.*")
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
			Context("when journal is called on one of the tasks that have not been deleted", func() {
				BeforeEach(func() {
					callRun("journal", "task-c")
				})
				It("shows the 1 creation event related to that task", func() {
					regexp := fmt.Sprintf("%s: %s", journalPrefixRegexp, ".*Created.*task-c")
					Expect(output.String()).To(MatchRegexp(regexp))
					Expect(strings.Count(output.String(), "\n")).To(Equal(1))
				})
			})
		})
		Context("when a note is added", func() {
			BeforeEach(func() {
				callRun("note", "task-a", "Here is a note for task-a")
				callRun("note", "task-c", "Here is a note for task-c")
			})
			Context("when journal is called", func() {
				BeforeEach(func() {
					callRun("journal")
				})
				It("shows the 3 creation entries plus the 2 notes in reverse order", func() {
					regexp := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s", journalPrefixRegexp,
						".*Note.*task-c.*Here is a note for task-c.*", journalPrefixRegexp,
						".*Note.*task-a.*Here is a note for task-a.*", journalPrefixRegexp,
						".*Created.*task-c.*", journalPrefixRegexp, ".*Created.*task-b.*", journalPrefixRegexp,
						".*Created.*task-a.*")
					Expect(output.String()).To(MatchRegexp(regexp))
				})
			})
		})
	})

	Context("when debug is on", func() {
		BeforeEach(func() {
			callRun("-debug", "create", "task-a")
		})
		It("reporting information about saving/loading manager", func() {
			Expect(output.String()).To(ContainSubstring("Manager is"))
			Expect(output.String()).To(ContainSubstring("Persisting manager back to disk"))
		})
	})

	Context("when the journal is requested for a task with a similar name to another task", func() {
		BeforeEach(func() {
			callRun("create", "task-a")
			callRun("create", "ask-a")
			callRun("set-running", "task-a")
			callRun("journal", "ask-a")
		})
		It("succeeds", expectSuccess)
		It("returns only the creation event for ask-a", func() {
			Expect(strings.Count(output.String(), "\n")).To(Equal(1))
		})
	})
})
