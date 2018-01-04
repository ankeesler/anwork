package main

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("anwork", func() {
	var output *bytes.Buffer
	var args []string
	JustBeforeEach(func() {
		output = new(bytes.Buffer)
		run(args, output)
	})

	Context("when no args are passed", func() {
		BeforeEach(func() {
			args = []string{"anwork"}
		})
		It("errors", func() {
			Expect(string(output.Bytes())).To(ContainSubstring("Error"))
			Expect(string(output.Bytes())).To(ContainSubstring("Expected command arguments"))
		})
	})
	Context("when help is requested", func() {
		BeforeEach(func() {
			args = []string{"anwork", "-h"}
		})
		It("shows the help", func() {
		})
	})
})
