package lag_test

import (
	"bytes"

	"github.com/ankeesler/anwork/lag"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("L", func() {
	It("prints formatted stuff to an io.Writer depending on the level", func() {
		buf := bytes.NewBuffer([]byte{})
		l := lag.New(buf)
		l.Level = lag.I

		l.P(lag.D, "here is a %s line", "debug")
		l.P(lag.I, "here is an %s line", "info")
		l.P(lag.E, "here is an %s %s", "error", "line")

		out := buf.String()
		Expect(out).To(MatchRegexp(`\[.*\] \(INFO\) here is an info line
\[.*\] \(ERROR\) here is an error line
`))
	})
})
