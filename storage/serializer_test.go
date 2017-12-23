package storage

import (
	"strings"

	pb "github.com/ankeesler/anwork/storage/proto"
	"github.com/golang/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serializable", func() {
	It("can serialize", func() {
		student := Student{"andrew", 18}
		bytes, err := student.Serialize()
		Expect(err).To(Succeed())

		unserializedStudent := Student{}
		err = unserializedStudent.Unserialize(bytes)
		Expect(err).To(Succeed())

		Expect(student).To(Equal(unserializedStudent))
	})
})

type Student struct {
	Name string
	Id   int32
}

func (s Student) Equal(otherS Student) bool {
	return strings.EqualFold(s.Name, otherS.Name) && s.Id == otherS.Id
}

func (s *Student) Serialize() ([]byte, error) {
	sProtobuf := pb.StudentProtobuf{Name: s.Name, Id: s.Id}
	return proto.Marshal(&sProtobuf)
}

func (s *Student) Unserialize(bytes []byte) error {
	sProtobuf := pb.StudentProtobuf{}
	err := proto.Unmarshal(bytes, &sProtobuf)
	if err != nil {
		return err
	}

	s.Name = sProtobuf.Name
	s.Id = sProtobuf.Id
	return nil
}
