package integration

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("anwork", func() {
	var (
		outBuf, errBuf *gbytes.Buffer
	)

	BeforeEach(func() {
		outBuf = gbytes.NewBuffer()
		errBuf = gbytes.NewBuffer()
	})

	Context("when no args are passed", func() {
		It("prints the usage", func() {
			run(outBuf, errBuf)
			Expect(outBuf).To(gbytes.Say("Usage of anwork"))
		})
	})

	Context("when help is requested", func() {
		It("prints the usage (only once!)", func() {
			run(outBuf, errBuf, "-h")
			Expect(outBuf).To(gbytes.Say("Usage of anwork"))
			Expect(outBuf).NotTo(gbytes.Say("Usage of anwork"))
		})
	})

	Context("when a bad flag is passed", func() {
		It("fails and prints the usage", func() {
			runWithStatus(1, outBuf, errBuf, "-tuna")
			Expect(outBuf).To(gbytes.Say("Usage of anwork"))
		})
	})

	Context("when a bad command is passed", func() {
		It("fails with a somewhat helpful error message", func() {
			runWithStatus(1, outBuf, errBuf, "tuna")
			Expect(errBuf).To(gbytes.Say("Unknown command: 'tuna'"))
		})
	})

	Context("when a command expects an arg but doesn't get one", func() {
		It("fails and prints the usage for that command", func() {
			runWithStatus(1, outBuf, errBuf, "create")
			Expect(errBuf).To(gbytes.Say("Got: \\[]"))
			Expect(errBuf).To(gbytes.Say("Expected: \\[task-name\\]"))
		})

		It("prints something about the missing argument", func() {
			runWithStatus(1, outBuf, errBuf, "create")
			Expect(errBuf).To(gbytes.Say("Invalid argument passed to command 'create'"))
		})
	})

	Context("when a bad context is passed", func() {
		It("fails", func() {
			if runWithApi {
				Skip("when connecting to API, we ignore the context flag")
			}
			runWithStatus(1, outBuf, errBuf, "-c", "/i/really/hope/this/file/doesnt/exist", "show")
		})
	})

	Context("when the context is corrupt", func() {
		It("fails", func() {
			if runWithApi {
				Skip("when connecting to API, we ignore the context flag")
			}
			runWithStatus(1, outBuf, errBuf, "-c", "data", "-o", "bad-context", "show")
		})
	})

	Context("when no output dir is passed", func() {
		var homeDir string
		BeforeEach(func() {
			var found bool
			homeDir, found = os.LookupEnv("HOME")
			Expect(found).To(BeTrue(), "Cannot find $HOME env var - this test needs that var to be set")
		})
		AfterEach(func() {
			Expect(os.RemoveAll(filepath.Join(homeDir, ".anwork", "home-dir-context"))).To(Succeed())
		})
		It("puts the data in the $HOME/.anwork/ directory", func() {
			if runWithApi {
				Skip("the API has its own problems regarding the use of the context...")
			}
			run(nil, nil, "-c", "home-dir-context", "create", "task-a")
			cmd := exec.Command(anworkBin, "-c", "home-dir-context", "create", "task-a")
			out, err := cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("out: %s", string(out)))

			Expect(filepath.Join(homeDir, ".anwork", "home-dir-context")).To(BeAnExistingFile())
		})
	})

	Context("when a zero-length output dir is passed", func() {
		It("fails", func() {
			if runWithApi {
				Skip("when connecting to API, we ignore the output flag")
			}
			runWithStatus(1, outBuf, errBuf, "-o", "", "show")
		})
	})

	Context("when the debug flag is passed", func() {
		It("prints out some extra stuff", func() {
			run(outBuf, errBuf, "-d", "show")
			Eventually(outBuf).Should(gbytes.Say("Manager is "))
		})
	})

	Context("when the version command is passed", func() {
		AfterEach(func() {
			run(nil, nil, "reset")
		})
		It("prints out the version, default build hash, and default date", func() {
			run(outBuf, errBuf, "version")
			Eventually(outBuf).Should(gbytes.Say("ANWORK Version = %d", version))
			Eventually(outBuf).Should(gbytes.Say(fmt.Sprintf("ANWORK Build Hash = \\(dev\\)\n")))
			Eventually(outBuf).Should(gbytes.Say(fmt.Sprintf("ANWORK Build Date = \\?\\?\\?\n")))
		})

		Context("when the binary is built via the official build script", func() {
			var (
				officialAnworkBin string
				outDir            string

				buildHash string
				buildDate string
			)
			BeforeEach(func() {
				buildHash = getBuildHash()
				buildDate = getBuildDate()
				officialAnworkBin = runOfficialBuildScript(buildHash, buildDate)

				var err error
				outDir, err = ioutil.TempDir("", "anwork-version-test")
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				Expect(os.RemoveAll(officialAnworkBin)).To(Succeed())
				Expect(os.RemoveAll(outDir)).To(Succeed())
			})
			It("prints out the version, build hash, and date", func() {
				cmd := exec.Command(officialAnworkBin, "-o", outDir, "version")
				out, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("output: %s", string(out)))

				outBuf := gbytes.BufferWithBytes([]byte(out))
				Expect(outBuf).To(gbytes.Say(fmt.Sprintf("ANWORK Version = %d\n", version)))
				Expect(outBuf).To(gbytes.Say(fmt.Sprintf("ANWORK Build Hash = %s\n", buildHash)))
				Expect(outBuf).To(gbytes.Say(fmt.Sprintf("ANWORK Build Date = %s\n", buildDate)))
			})
		})
	})

	Context("when creating a task", func() {
		BeforeEach(func() {
			run(outBuf, errBuf, "create", "task-a")
		})
		AfterEach(func() {
			run(nil, nil, "reset")
		})
		It("shows the task as ready", func() {
			run(outBuf, errBuf, "show")
			Expect(outBuf).To(gbytes.Say("READY tasks:\n  task-a \\(\\d+\\)\nFINISHED tasks"))
		})
		It("shows the correct task details", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("Name: task-a\nID: \\d+\nCreated: .*\nPriority: 10\nState: READY"))
		})
		It("records the event in the task's journal", func() {
			run(outBuf, errBuf, "journal", "task-a")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-a'"))
		})
		It("records the event in the globas journal", func() {
			run(outBuf, errBuf, "journal")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-a'"))
		})
	})

	Context("when creating multiple tasks", func() {
		BeforeEach(func() {
			run(outBuf, errBuf, "create", "task-a")
			run(outBuf, errBuf, "create", "task-b")
			run(outBuf, errBuf, "create", "task-c")
		})
		AfterEach(func() {
			run(nil, nil, "reset")
		})
		It("shows the tasks as ready in the order in which they were created", func() {
			run(outBuf, errBuf, "show")
			Expect(outBuf).To(gbytes.Say("READY tasks:\n  task-a \\(\\d+\\)\n  task-b \\(\\d+\\)\n  task-c \\(\\d+\\)\nFINISHED tasks"))
		})
		It("shows the correct task details", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("Name: task-a\nID: \\d+\nCreated: .*\nPriority: 10\nState: READY"))
			run(outBuf, errBuf, "show", "task-b")
			Expect(outBuf).To(gbytes.Say("Name: task-b\nID: \\d+\nCreated: .*\nPriority: 10\nState: READY"))
			run(outBuf, errBuf, "show", "task-c")
			Expect(outBuf).To(gbytes.Say("Name: task-c\nID: \\d+\nCreated: .*\nPriority: 10\nState: READY"))
		})
		It("records the events in each task's journal", func() {
			run(outBuf, errBuf, "journal", "task-a")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-a'"))
			run(outBuf, errBuf, "journal", "task-b")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-b'"))
			run(outBuf, errBuf, "journal", "task-c")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-c'"))
		})
		It("records the events in the global journal in most recent to oldest order", func() {
			run(outBuf, errBuf, "journal")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task 'task-c'\n\\[.*\\]: Created task 'task-b'\n\\[.*\\]: Created task 'task-a'"))
		})
		Context("when changing the priority on tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "set-priority", "task-a", "15")
				run(nil, nil, "set-priority", "task-b", "10")
				run(nil, nil, "set-priority", "task-c", "20")
			})
			It("properly shows the tasks in order of priority", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("READY tasks:"))
				Expect(outBuf).To(gbytes.Say("  task-b"))
				Expect(outBuf).To(gbytes.Say("  task-a"))
				Expect(outBuf).To(gbytes.Say("  task-c"))

			})
			It("properly records the priorities", func() {
				run(outBuf, errBuf, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: \\d+\nCreated: .*\nPriority: 15\nState: READY"))
				run(outBuf, errBuf, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: \\d+\nCreated: .*\nPriority: 10\nState: READY"))
				run(outBuf, errBuf, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: \\d+\nCreated: .*\nPriority: 20\nState: READY"))
			})
			It("records the events in each of the task's journals", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-a' from 10 to 15"))
				run(outBuf, errBuf, "journal", "task-b")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-b' from 10 to 10"))
				run(outBuf, errBuf, "journal", "task-c")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-c' from 10 to 20"))
			})
			It("records the events in the global journal", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-c' from 10 to 20"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-b' from 10 to 10"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task 'task-a' from 10 to 15"))
			})
		})
		Context("when changing the state on tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "set-running", "task-a")
				run(nil, nil, "set-finished", "task-b")
				run(nil, nil, "set-blocked", "task-c")
			})
			It("properly displays the states", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:\n  task-a"))
				Expect(outBuf).To(gbytes.Say("BLOCKED tasks:\n  task-c"))
				Expect(outBuf).To(gbytes.Say("FINISHED tasks:\n  task-b"))
			})
			It("properly records the states", func() {
				run(outBuf, errBuf, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: \\d+\nCreated: .*\nPriority: 10\nState: RUNNING"))
				run(outBuf, errBuf, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: \\d+\nCreated: .*\nPriority: 10\nState: FINISHED"))
				run(outBuf, errBuf, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: \\d+\nCreated: .*\nPriority: 10\nState: BLOCKED"))
			})
			It("records the events in each of the task's journals", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-a' from Ready to Running"))
				run(outBuf, errBuf, "journal", "task-b")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-b' from Ready to Finished"))
				run(outBuf, errBuf, "journal", "task-c")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-c' from Ready to Blocked"))
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-c' from Ready to Blocked"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-b' from Ready to Finished"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-a' from Ready to Running"))
			})
		})
		Context("when adding a note to tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "note", "task-a", "Here is a note")
			})
			It("records the note in task-a's journal", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Note added to task 'task-a': Here is a note"))
			})
			It("records the note in the global journal", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Note added to task 'task-a': Here is a note"))
			})
		})
		Context("when deleting a task", func() {
			BeforeEach(func() {
				run(nil, nil, "delete", "task-b")
			})
			It("properly displays the remaining tasks", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("READY tasks:\n  task-a.*\n  task-c"))
			})
			It("fails when we try to show the deleted task", func() {
				runWithStatus(1, outBuf, errBuf, "show", "task-b")
			})
			It("fails to show the journal for the deleted task", func() {
				runWithStatus(1, outBuf, errBuf, "journal", "task-b")
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-b'"))
			})
		})
		Context("when deleting all tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "delete-all")
			})
			It("no longer shows any tasks", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:\nBLOCKED tasks:\nREADY tasks:\nFINISHED tasks"))
			})
			It("fails when we try to show the deleted tasks", func() {
				runWithStatus(1, outBuf, errBuf, "show", "task-a")
				runWithStatus(1, outBuf, errBuf, "show", "task-b")
				runWithStatus(1, outBuf, errBuf, "show", "task-c")
			})
			It("fails to show the journal for the deleted tasks", func() {
				runWithStatus(1, outBuf, errBuf, "journal", "task-a")
				runWithStatus(1, outBuf, errBuf, "journal", "task-b")
				runWithStatus(1, outBuf, errBuf, "journal", "task-c")
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-c'"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-b'"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-a'"))
			})
		})
		Context("when doing a reset", func() {
			var (
				s     *gexec.Session
				stdin io.Writer
			)
			BeforeEach(func() {
				Expect(os.Unsetenv("ANWORK_TEST_RESET_ANSWER")).To(Succeed())

				var err error
				cmd := exec.Command(anworkBin, "-o", outputDir, "reset")
				stdin, err = cmd.StdinPipe()
				Expect(err).ToNot(HaveOccurred())

				s, err = gexec.Start(cmd, outBuf, errBuf)
				Expect(err).ToNot(HaveOccurred())

				Eventually(outBuf).Should(gbytes.Say("Are you sure you want to delete all data \\[y/n\\]: "))

				Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "y")).To(Succeed())
			})
			AfterEach(func() {
				s.Interrupt()
			})
			It("does not delete all data if the user says no", func() {
				stdin.Write([]byte("n\n"))
				Eventually(outBuf).Should(gbytes.Say("NOT deleting all data"))

				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("task-a"))
			})
			It("does delete all data if the user says yes", func() {
				stdin.Write([]byte("y\n"))
				Eventually(outBuf).Should(gbytes.Say("OK, deleting all data"))
				Eventually(s).Should(gexec.Exit(0))

				run(outBuf, errBuf, "show")
				Expect(outBuf).ToNot(gbytes.Say("task-a"))
			})
		})
		Context("when archiving", func() {
			BeforeEach(func() {
				run(nil, nil, "create", "task-d")
				run(nil, nil, "set-finished", "task-b")
				run(nil, nil, "set-running", "task-c")
				run(nil, nil, "set-finished", "task-a")
				run(nil, nil, "archive")
			})
			It("no longer shows any of the finished tasks", func() {
				run(outBuf, errBuf, "show")
				output := string(outBuf.Contents())
				expectedOutput := `RUNNING tasks:
  task-c \(\d+\)
BLOCKED tasks:
READY tasks:
  task-d \(\d+\)
FINISHED tasks:`
				Expect(output).To(MatchRegexp(expectedOutput))
			})
			It("fails when we try to show the finished tasks", func() {
				runWithStatus(1, nil, nil, "show", "task-a")
				runWithStatus(1, nil, nil, "show", "task-b")
				run(nil, nil, "show", "task-c")
				run(nil, nil, "show", "task-d")
			})
			It("fails to show the journal for the finished tasks", func() {
				runWithStatus(1, nil, nil, "journal", "task-a")
				runWithStatus(1, nil, nil, "journal", "task-b")
				run(nil, nil, "journal", "task-c")
				run(nil, nil, "journal", "task-d")
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-b'"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task 'task-a'"))
			})
		})
		Context("when performing a summary", func() {
			BeforeEach(func() {
				run(nil, nil, "set-finished", "task-b")
			})
			It("reports the finished tasks", func() {
				run(outBuf, errBuf, "summary", "1")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-b' from Ready to Finished"))
				Expect(outBuf).To(gbytes.Say("  took "))
			})
		})
	})

	Context("when creating a task, deleting it, and creating it again", func() {
		var id int
		BeforeEach(func() {
			run(nil, nil, "create", "task-a")
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("ID: \\d+"))

			re, err := regexp.Compile(`ID: (\d+)`)
			Expect(err).NotTo(HaveOccurred())

			matches := re.FindSubmatch(outBuf.Contents())
			Expect(matches).To(HaveLen(2), fmt.Sprintf("Matches is: %s", matches))

			id, err = strconv.Atoi(string(matches[1]))
			Expect(err).NotTo(HaveOccurred())

			run(nil, nil, "delete", "task-a")
			runWithStatus(1, outBuf, errBuf, "show", "task-a")

			run(nil, nil, "create", "task-a")
		})
		It("should give unique IDs to both tasks", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).ToNot(gbytes.Say(fmt.Sprintf("ID: %d", id)))
		})
	})

	Measure("CRUD'ing 10 tasks", func(b Benchmarker) {
		defer run(nil, nil, "reset")

		runtime := b.Time("runtime", func() {
			for i := 0; i < 10; i++ {
				name := fmt.Sprintf("task-%d", i)
				run(nil, nil, "create", name)
				run(nil, nil, "show", name)
				run(nil, nil, "set-finished", name)
				run(nil, nil, "delete", name)
			}
		})
		Expect(runtime.Seconds()).To(BeNumerically("<", 2))
	}, 5)
})
