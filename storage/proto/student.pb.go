// Code generated by protoc-gen-go. DO NOT EDIT.
// source: student.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	student.proto

It has these top-level messages:
	StudentProtobuf
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type StudentProtobuf struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Id   int32  `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
}

func (m *StudentProtobuf) Reset()                    { *m = StudentProtobuf{} }
func (m *StudentProtobuf) String() string            { return proto1.CompactTextString(m) }
func (*StudentProtobuf) ProtoMessage()               {}
func (*StudentProtobuf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StudentProtobuf) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StudentProtobuf) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto1.RegisterType((*StudentProtobuf)(nil), "StudentProtobuf")
}

func init() { proto1.RegisterFile("student.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 95 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x2e, 0x29, 0x4d,
	0x49, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x32, 0xe5, 0xe2, 0x0f, 0x86, 0x08,
	0x04, 0x80, 0xf8, 0x49, 0xa5, 0x69, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x10, 0x1f, 0x17, 0x53, 0x66, 0x8a, 0x04, 0x93, 0x02,
	0xa3, 0x06, 0x6b, 0x10, 0x53, 0x66, 0x8a, 0x13, 0x7b, 0x14, 0x2b, 0x58, 0x7f, 0x12, 0x1b, 0x98,
	0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x0a, 0x39, 0xa6, 0x57, 0x00, 0x00, 0x00,
}