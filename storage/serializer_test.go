package storage

import (
	"strings"
	"testing"

	pb "github.com/ankeesler/anwork/storage/proto"
	"github.com/golang/protobuf/proto"
)

// This is a test object used for serialization. It is serialized via the StudentProtobuf definition.
//go:generate protoc --proto_path=proto --go_out=proto proto/student.proto
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

func TestSerialize(t *testing.T) {
	student := Student{"andrew", 18}
	bytes, err := student.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize student: %s", err)
	}

	unserializedStudent := Student{}
	err = unserializedStudent.Unserialize(bytes)
	if err != nil {
		t.Fatalf("Failed to unserialize student: %s", err)
	}

	if !student.Equal(unserializedStudent) {
		t.Fatalf("Unserialized student (%v) is not equal to original student (%v)",
			unserializedStudent, student)
	}
}

func TestSerializable(t *testing.T) {
	student := &Student{}
	func(s Serializable) {}(student)
}
