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

const root = "test-data"
const context = "test-context"

var _ = Describe("anwork", func() {
	var output *bytes.Buffer
	var args []string
	var ret int
	JustBeforeEach(func() {
		output = new(bytes.Buffer)
		ret = run(append([]string{"anwork", "-context", context, "-root", root}, args...), output)
	})

	AfterEach(func() {
		p := storage.Persister{root}
		Expect(p.Delete(context)).To(Succeed())
	})

	expectSuccess := func() {
		Expect(ret).To(BeEquivalentTo(0))
	}
	expectFailure := func() {
		Expect(ret).ToNot(BeEquivalentTo(0))
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

	Context("when no args are passed", func() {
		BeforeEach(func() {
			args = []string{}
		})
		It("fails", expectFailure)
		It("prints usage", expectUsagePrinted)
		It("prints error", func() {
			Expect(output.String()).To(ContainSubstring("Error! Expected command arguments"))
		})
	})
	Context("when help is requested", func() {
		BeforeEach(func() {
			args = []string{"-h"}
		})
		It("succeeds", expectSuccess)
		It("prints usage", expectUsagePrinted)
		It("prints usage only once!", func() {
			firstIndex := strings.Index(output.String(), "Usage of anwork")
			Expect(firstIndex).ToNot(BeEquivalentTo(-1))

			secondIndex := strings.Index(output.String()[firstIndex+1:], "Usage of anwork")
			Expect(secondIndex).To(BeEquivalentTo(-1))
		})
	})
	Context("when a bad flag is passed", func() {
		BeforeEach(func() {
			args = []string{"-tuna"}
		})
		It("fails", expectFailure)
		It("prints usage", expectUsagePrinted)
	})
	Context("when a bad command is passed", func() {
		BeforeEach(func() {
			args = []string{"fish"}
		})
		It("fails", expectFailure)
		It("prints usage", expectUsagePrinted)
	})

	Context("when the version command is passed", func() {
		BeforeEach(func() {
			args = []string{"version"}
		})
		It("succeeds", expectSuccess)
		It("prints the version", func() {
			msg := fmt.Sprintf("ANWORK Version = %d", command.Version)
			Expect(output.String()).To(ContainSubstring(msg))
		})
	})

	Context("when the set-priority command is passed", func() {
	})

	// TODO: write tests for the rest of the commands!
})
