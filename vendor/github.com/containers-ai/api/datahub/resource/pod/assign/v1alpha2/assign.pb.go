// Code generated by protoc-gen-go. DO NOT EDIT.
// source: datahub/resource/pod/assign/v1alpha2/assign.proto

package v1alpha2 // import "github.com/containers-ai/api/datahub/resource/pod/assign/v1alpha2"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// *
// Represents the priority of a node
type NodePriority struct {
	Nodes                []string `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodePriority) Reset()         { *m = NodePriority{} }
func (m *NodePriority) String() string { return proto.CompactTextString(m) }
func (*NodePriority) ProtoMessage()    {}
func (*NodePriority) Descriptor() ([]byte, []int) {
	return fileDescriptor_assign_80d2fe6ada307263, []int{0}
}
func (m *NodePriority) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodePriority.Unmarshal(m, b)
}
func (m *NodePriority) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodePriority.Marshal(b, m, deterministic)
}
func (dst *NodePriority) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodePriority.Merge(dst, src)
}
func (m *NodePriority) XXX_Size() int {
	return xxx_messageInfo_NodePriority.Size(m)
}
func (m *NodePriority) XXX_DiscardUnknown() {
	xxx_messageInfo_NodePriority.DiscardUnknown(m)
}

var xxx_messageInfo_NodePriority proto.InternalMessageInfo

func (m *NodePriority) GetNodes() []string {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type Selector struct {
	Selector             map[string]string `protobuf:"bytes,1,rep,name=selector,proto3" json:"selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Selector) Reset()         { *m = Selector{} }
func (m *Selector) String() string { return proto.CompactTextString(m) }
func (*Selector) ProtoMessage()    {}
func (*Selector) Descriptor() ([]byte, []int) {
	return fileDescriptor_assign_80d2fe6ada307263, []int{1}
}
func (m *Selector) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Selector.Unmarshal(m, b)
}
func (m *Selector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Selector.Marshal(b, m, deterministic)
}
func (dst *Selector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Selector.Merge(dst, src)
}
func (m *Selector) XXX_Size() int {
	return xxx_messageInfo_Selector.Size(m)
}
func (m *Selector) XXX_DiscardUnknown() {
	xxx_messageInfo_Selector.DiscardUnknown(m)
}

var xxx_messageInfo_Selector proto.InternalMessageInfo

func (m *Selector) GetSelector() map[string]string {
	if m != nil {
		return m.Selector
	}
	return nil
}

func init() {
	proto.RegisterType((*NodePriority)(nil), "containersai.datahub.resource.pod.assign.v1alpha2.NodePriority")
	proto.RegisterType((*Selector)(nil), "containersai.datahub.resource.pod.assign.v1alpha2.Selector")
	proto.RegisterMapType((map[string]string)(nil), "containersai.datahub.resource.pod.assign.v1alpha2.Selector.SelectorEntry")
}

func init() {
	proto.RegisterFile("datahub/resource/pod/assign/v1alpha2/assign.proto", fileDescriptor_assign_80d2fe6ada307263)
}

var fileDescriptor_assign_80d2fe6ada307263 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x50, 0xc1, 0x4a, 0x03, 0x31,
	0x14, 0x24, 0x2d, 0x4a, 0x1b, 0x15, 0x64, 0xf1, 0xb0, 0x78, 0x2a, 0xc5, 0x43, 0x2f, 0x26, 0x6c,
	0xbd, 0x88, 0x9e, 0x54, 0x3c, 0x78, 0x11, 0x59, 0x6f, 0xde, 0x5e, 0x37, 0x8f, 0x36, 0xb8, 0xe6,
	0x85, 0x97, 0x6c, 0x61, 0x7f, 0xca, 0x6f, 0x94, 0xdd, 0x66, 0x2b, 0xde, 0xec, 0x6d, 0x26, 0x99,
	0x99, 0x37, 0x8c, 0x2c, 0x0c, 0x44, 0xd8, 0x34, 0x2b, 0xcd, 0x18, 0xa8, 0xe1, 0x0a, 0xb5, 0x27,
	0xa3, 0x21, 0x04, 0xbb, 0x76, 0x7a, 0x5b, 0x40, 0xed, 0x37, 0xb0, 0x4c, 0x5c, 0x79, 0xa6, 0x48,
	0x59, 0x51, 0x91, 0x8b, 0x60, 0x1d, 0x72, 0x00, 0xab, 0x92, 0x5f, 0x0d, 0x7e, 0xe5, 0xc9, 0xa8,
	0xa4, 0x1f, 0xfc, 0xf3, 0x2b, 0x79, 0xfa, 0x4a, 0x06, 0xdf, 0xd8, 0x12, 0xdb, 0xd8, 0x66, 0x17,
	0xf2, 0xc8, 0x91, 0xc1, 0x90, 0x8b, 0xd9, 0x78, 0x31, 0x2d, 0x77, 0x64, 0xfe, 0x2d, 0xe4, 0xe4,
	0x1d, 0x6b, 0xac, 0x22, 0x71, 0x86, 0x72, 0x12, 0x12, 0xee, 0x55, 0x27, 0xcb, 0x17, 0x75, 0xf0,
	0x61, 0x35, 0xc4, 0xed, 0xc1, 0xb3, 0x8b, 0xdc, 0x96, 0xfb, 0xe8, 0xcb, 0x7b, 0x79, 0xf6, 0xe7,
	0x2b, 0x3b, 0x97, 0xe3, 0x4f, 0x6c, 0x73, 0x31, 0x13, 0x8b, 0x69, 0xd9, 0xc1, 0xae, 0xec, 0x16,
	0xea, 0x06, 0xf3, 0x51, 0xff, 0xb6, 0x23, 0x77, 0xa3, 0x5b, 0xf1, 0xf8, 0xf4, 0xf1, 0xb0, 0xb6,
	0xb1, 0x2b, 0x51, 0xd1, 0x97, 0xfe, 0x6d, 0x77, 0x0d, 0x56, 0x83, 0xb7, 0xfa, 0x3f, 0xdb, 0xae,
	0x8e, 0xfb, 0x55, 0x6f, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x33, 0xa8, 0x0d, 0x8a, 0x01,
	0x00, 0x00,
}
