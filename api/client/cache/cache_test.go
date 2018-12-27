package cache_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ankeesler/anwork/api/client/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		c *cache.Cache

		dir string
	)

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "anwork-cache-test")
		Expect(err).NotTo(HaveOccurred())

		c = cache.New(filepath.Join(dir, "cache"))
	})

	AfterEach(func() {
		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	Context("when the cache is empty to start", func() {
		BeforeEach(func() {
			_, ok := c.Get()
			Expect(ok).To(BeFalse())
		})

		It("returns stuff from Get() that you Set()", func() {
			c.Set("some string")

			s, ok := c.Get()
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal("some string"))
		})
	})

	Context("when the cache is not empty to start", func() {
		BeforeEach(func() {
			c.Set("some string 1")
		})

		It("properly updates", func() {
			c.Set("some string 2")

			s, ok := c.Get()
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal("some string 2"))
		})
	})
})
