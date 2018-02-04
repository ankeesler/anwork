package storage

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

const (
	goodRoot = "test-data/good-root"
	badRoot  = "test-data/bad-root"

	emptyContext     = "empty-context"
	singletonContext = "singleton-context"
	badContext       = "bad-context"

	tmpRoot = "test-data/root.tmp"

	tmpEmptyContext     = "empty-context.tmp"
	tmpSingletonContext = "singleton-context.tmp"
)

var _ = Describe("FilePersister", func() {
	Describe("existence", func() {
		var (
			root    string
			context string
			p       FilePersister
		)

		JustBeforeEach(func() {
			p = FilePersister{root}
		})

		Context("when the root is valid", func() {
			BeforeEach(func() {
				root = goodRoot
			})
			Context("when the context is valid", func() {
				BeforeEach(func() {
					context = emptyContext
				})
				It("returns true", func() {
					Expect(p.Exists(context)).To(BeTrue())
				})
			})
			Context("when the context is invalid", func() {
				BeforeEach(func() {
					context = badContext
				})
				It("returns false", func() {
					Expect(p.Exists(context)).To(BeFalse())
				})
			})
		})

		Context("when the root is invalid", func() {
			BeforeEach(func() {
				root = badRoot
			})
			Context("when the context is valid", func() {
				BeforeEach(func() {
					context = emptyContext
				})
				It("returns false", func() {
					Expect(p.Exists(context)).To(BeFalse())
				})
			})
			Context("when the context is invalid", func() {
				BeforeEach(func() {
					context = badContext
				})
				It("returns false", func() {
					Expect(p.Exists(context)).To(BeFalse())
				})
			})
		})
	})
	Describe("persist", func() {
		var (
			root    string
			context string
			s       Serializable
			p       FilePersister
		)

		BeforeEach(func() {
			Expect(tmpRoot).ToNot(BeAnExistingFile(),
				"Error: cannot run test when tmpRoot (%s) exists", tmpRoot)
		})
		AfterEach(func() {
			Expect(os.RemoveAll(tmpRoot)).To(Succeed())
		})

		JustBeforeEach(func() {
			p = FilePersister{root}
		})

		Context("when the root is bad", func() {
			BeforeEach(func() {
				root = "/this/file/doesnt/exist"
				context = tmpEmptyContext
				s = &GoodSerializable{}
			})
			It("fails", func() {
				Expect(p.Persist(context, s)).ToNot(Succeed())
			})
		})

		Context("when the root is good", func() {
			BeforeEach(func() {
				root = tmpRoot
			})
			Context("when the serializer fails", func() {
				BeforeEach(func() {
					s = &BadSerializable{}
				})
				It("fails", func() {
					Expect(p.Persist(context, s)).ToNot(Succeed())
				})
			})
			Context("when the serializer is good", func() {
				Context("when no bytes are written", func() {
					BeforeEach(func() {
						s = &GoodSerializable{}
						context = tmpEmptyContext
					})
					It("passes", func() {
						Expect(p.Persist(context, s)).To(Succeed())
						Expect(p.Unpersist(context, s)).To(Succeed())
						Expect(s).To(BeAssignableToTypeOf(&GoodSerializable{}))
						gs := s.(*GoodSerializable)
						Expect(gs.actualBytes).To(Equal(gs.ExpectedBytes))
					})
				})
				Context("when some bytes are written once", func() {
					BeforeEach(func() {
						s = &GoodSerializable{ExpectedBytes: []byte{'a', 'b', 'c'}}
						context = singletonContext
					})
					It("passes", func() {
						Expect(p.Persist(context, s)).To(Succeed())
						Expect(p.Unpersist(context, s)).To(Succeed())
						Expect(s).To(BeAssignableToTypeOf(&GoodSerializable{}))
						gs := s.(*GoodSerializable)
						Expect(gs.actualBytes).To(Equal(gs.ExpectedBytes))
					})
				})
				Context("when some bytes are written twice", func() {
					BeforeEach(func() {
						s = &GoodSerializable{ExpectedBytes: []byte{'c', 'b', 'a'}}
						context = singletonContext
					})
					It("passes", func() {
						Expect(p.Persist(context, s)).To(Succeed())
						Expect(p.Unpersist(context, s)).To(Succeed())
						Expect(s).To(BeAssignableToTypeOf(&GoodSerializable{}))
						gs := s.(*GoodSerializable)
						Expect(gs.actualBytes).To(Equal(gs.ExpectedBytes))
					})
				})
			})
		})
	})

	Describe("unpersist", func() {
		var (
			root    string
			context string
			s       Serializable
			p       FilePersister
		)

		JustBeforeEach(func() {
			p = FilePersister{root}
		})

		Context("when the root doesn't exist", func() {
			BeforeEach(func() {
				root = badRoot
				context = emptyContext
				s = &GoodSerializable{}
			})
			It("indicates that the context does not exist", func() {
				Expect(p.Exists(context)).To(BeFalse())
			})
			It("fails", func() {
				Expect(p.Unpersist(context, s)).ToNot(Succeed())
			})
		})

		Context("when the root does exist", func() {
			BeforeEach(func() {
				root = goodRoot
			})
			Context("when the context doesn't exist", func() {
				BeforeEach(func() {
					context = badRoot
					s = &GoodSerializable{}
				})
				It("indicates that the context does not exist", func() {
					Expect(p.Exists(context)).To(BeFalse())
				})
				It("fails", func() {
					Expect(p.Unpersist(context, s)).ToNot(Succeed())
				})
			})
			Context("when the context does exist", func() {
				Context("when the context has 0 bytes", func() {
					BeforeEach(func() {
						context = emptyContext
						s = &GoodSerializable{}
					})
					It("indicates that the context does exist", func() {
						Expect(p.Exists(context)).To(BeTrue())
					})
					It("passes", func() {
						Expect(p.Unpersist(context, s)).To(Succeed())
						Expect(s).To(BeAssignableToTypeOf(&GoodSerializable{}))
						gs := s.(*GoodSerializable)
						Expect(gs.actualBytes).To(Equal(gs.ExpectedBytes))
					})
				})
				Context("when the context has more than 1 bytes", func() {
					BeforeEach(func() {
						context = singletonContext
						s = &GoodSerializable{ExpectedBytes: []byte{'a', 'b', 'c'}}
					})
					It("indicates that the context does exist", func() {
						Expect(p.Exists(context)).To(BeTrue())
					})
					It("passes", func() {
						Expect(p.Unpersist(context, s)).To(Succeed())
						Expect(s).To(BeAssignableToTypeOf(&GoodSerializable{}))
						gs := s.(*GoodSerializable)
						Expect(gs.actualBytes).To(Equal(gs.ExpectedBytes))
					})
				})
			})
		})
	})
	Describe("delete", func() {
		var (
			context string
			p       FilePersister
		)

		BeforeEach(func() {
			Expect(tmpRoot).ToNot(BeAnExistingFile(),
				"Error: cannot run test when tmpRoot (%s) exists", tmpRoot)
		})
		AfterEach(func() {
			Expect(os.RemoveAll(tmpRoot)).To(Succeed())
		})

		JustBeforeEach(func() {
			p = FilePersister{tmpRoot}
		})

		Context("when the root does not exist", func() {
			BeforeEach(func() {
				context = tmpEmptyContext
			})
			It("returns no error", func() {
				Expect(p.Delete(context)).To(Succeed())
			})
		})

		Context("when the root does exist", func() {
			BeforeEach(func() {
				Expect(p.Persist(tmpSingletonContext, &GoodSerializable{})).To(Succeed())
			})
			Context("when the context does not exist", func() {
				BeforeEach(func() {
					context = badContext
				})
				It("reports that the context does not exist", func() {
					Expect(p.Exists(context)).To(BeFalse())
				})
				It("returns no error", func() {
					Expect(p.Delete(context)).To(Succeed())
				})
			})
			Context("when the context does exist", func() {
				BeforeEach(func() {
					context = tmpSingletonContext
				})
				It("reports that the context does exist", func() {
					Expect(p.Exists(context)).To(BeTrue())
				})
				It("returns no error and removes the context", func() {
					Expect(p.Delete(context)).To(Succeed())
					Expect(p.Exists(context)).To(BeFalse())
				})
			})
		})
	})
})

type GoodSerializable struct {
	ExpectedBytes []byte
	actualBytes   []byte
}

func (s *GoodSerializable) Serialize() ([]byte, error) {
	return s.ExpectedBytes, nil
}

func (s *GoodSerializable) Unserialize(bytes []byte) error {
	s.actualBytes = bytes
	return nil
}

type BadSerializable struct {
}

func (s *BadSerializable) Serialize() ([]byte, error) {
	return nil, errors.New("Failure!")
}

func (s *BadSerializable) Unserialize(bytes []byte) error {
	return errors.New("Failure!")
}
