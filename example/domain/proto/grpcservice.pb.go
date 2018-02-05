// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpcservice.proto

/*
Package grpcservice is a generated protocol buffer package.

It is generated from these files:
	grpcservice.proto

It has these top-level messages:
	WrinkledItem
	SmoothItem
	Item
*/
package grpcservice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// the city where the branch of our company is located
type Location int32

const (
	Location_Basel  Location = 0
	Location_Zurich Location = 1
	Location_Luzern Location = 2
	Location_Other  Location = 3
)

var Location_name = map[int32]string{
	0: "Basel",
	1: "Zurich",
	2: "Luzern",
	3: "Other",
}
var Location_value = map[string]int32{
	"Basel":  0,
	"Zurich": 1,
	"Luzern": 2,
	"Other":  3,
}

func (x Location) String() string {
	return proto.EnumName(Location_name, int32(x))
}
func (Location) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// WrinkledItem is used for items with multiple wrinkes
type WrinkledItem struct {
	// customer will contain the name of the customer
	Customer string   `protobuf:"bytes,1,opt,name=Customer" json:"Customer,omitempty"`
	Location Location `protobuf:"varint,2,opt,name=Location,enum=grpcservice.Location" json:"Location,omitempty"`
	// nested messages can be embedded directly using the tag compose:"embed"
	Item *Item `protobuf:"bytes,4,opt,name=Item" json:"Item,omitempty"`
}

func (m *WrinkledItem) Reset()                    { *m = WrinkledItem{} }
func (m *WrinkledItem) String() string            { return proto.CompactTextString(m) }
func (*WrinkledItem) ProtoMessage()               {}
func (*WrinkledItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *WrinkledItem) GetCustomer() string {
	if m != nil {
		return m.Customer
	}
	return ""
}

func (m *WrinkledItem) GetLocation() Location {
	if m != nil {
		return m.Location
	}
	return Location_Basel
}

func (m *WrinkledItem) GetItem() *Item {
	if m != nil {
		return m.Item
	}
	return nil
}

// SmoothItem is used as response with associated costs
type SmoothItem struct {
	Item *Item `protobuf:"bytes,1,opt,name=Item" json:"Item,omitempty"`
	Cost int32 `protobuf:"varint,2,opt,name=Cost" json:"Cost,omitempty"`
}

func (m *SmoothItem) Reset()                    { *m = SmoothItem{} }
func (m *SmoothItem) String() string            { return proto.CompactTextString(m) }
func (*SmoothItem) ProtoMessage()               {}
func (*SmoothItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SmoothItem) GetItem() *Item {
	if m != nil {
		return m.Item
	}
	return nil
}

func (m *SmoothItem) GetCost() int32 {
	if m != nil {
		return m.Cost
	}
	return 0
}

// Item contains the information about a specific item
type Item struct {
	Name     string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Wrinkels int32  `protobuf:"varint,2,opt,name=Wrinkels" json:"Wrinkels,omitempty"`
}

func (m *Item) Reset()                    { *m = Item{} }
func (m *Item) String() string            { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()               {}
func (*Item) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Item) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Item) GetWrinkels() int32 {
	if m != nil {
		return m.Wrinkels
	}
	return 0
}

func init() {
	proto.RegisterType((*WrinkledItem)(nil), "grpcservice.WrinkledItem")
	proto.RegisterType((*SmoothItem)(nil), "grpcservice.SmoothItem")
	proto.RegisterType((*Item)(nil), "grpcservice.Item")
	proto.RegisterEnum("grpcservice.Location", Location_name, Location_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Butler service

type ButlerClient interface {
	// Iron will take a wrinkled item and remove all wrinkels
	Iron(ctx context.Context, in *WrinkledItem, opts ...grpc.CallOption) (*SmoothItem, error)
}

type butlerClient struct {
	cc *grpc.ClientConn
}

func NewButlerClient(cc *grpc.ClientConn) ButlerClient {
	return &butlerClient{cc}
}

func (c *butlerClient) Iron(ctx context.Context, in *WrinkledItem, opts ...grpc.CallOption) (*SmoothItem, error) {
	out := new(SmoothItem)
	err := grpc.Invoke(ctx, "/grpcservice.Butler/Iron", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Butler service

type ButlerServer interface {
	// Iron will take a wrinkled item and remove all wrinkels
	Iron(context.Context, *WrinkledItem) (*SmoothItem, error)
}

func RegisterButlerServer(s *grpc.Server, srv ButlerServer) {
	s.RegisterService(&_Butler_serviceDesc, srv)
}

func _Butler_Iron_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WrinkledItem)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ButlerServer).Iron(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcservice.Butler/Iron",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ButlerServer).Iron(ctx, req.(*WrinkledItem))
	}
	return interceptor(ctx, in, info, handler)
}

var _Butler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpcservice.Butler",
	HandlerType: (*ButlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Iron",
			Handler:    _Butler_Iron_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpcservice.proto",
}

func init() { proto.RegisterFile("grpcservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0xed, 0xd6, 0x34, 0xb4, 0x53, 0x91, 0x74, 0x40, 0x8c, 0x3d, 0x85, 0x80, 0x10, 0x3c, 0x14,
	0x8c, 0x20, 0x1e, 0x3c, 0xb5, 0xa0, 0x14, 0x8a, 0x42, 0x3c, 0x08, 0xde, 0x62, 0x1c, 0x4c, 0x30,
	0xc9, 0x96, 0xd9, 0x8d, 0x07, 0xcf, 0x7e, 0xb8, 0xec, 0xd2, 0xa6, 0xc9, 0xa5, 0xb7, 0x37, 0xfb,
	0xde, 0x1b, 0xde, 0xdb, 0x81, 0xd9, 0x17, 0x6f, 0x33, 0x45, 0xfc, 0x53, 0x64, 0xb4, 0xd8, 0xb2,
	0xd4, 0x12, 0xa7, 0x9d, 0xa7, 0xf0, 0x4f, 0xc0, 0xe9, 0x1b, 0x17, 0xf5, 0x77, 0x49, 0x9f, 0x6b,
	0x4d, 0x15, 0xce, 0x61, 0xbc, 0x6a, 0x94, 0x96, 0x15, 0xb1, 0x2f, 0x02, 0x11, 0x4d, 0x92, 0x76,
	0xc6, 0x1b, 0x18, 0x6f, 0x64, 0x96, 0xea, 0x42, 0xd6, 0xfe, 0x30, 0x10, 0xd1, 0x59, 0x7c, 0xbe,
	0xe8, 0xee, 0xdf, 0x93, 0x49, 0x2b, 0xc3, 0x2b, 0x70, 0xcc, 0x5a, 0xdf, 0x09, 0x44, 0x34, 0x8d,
	0x67, 0x3d, 0xb9, 0x21, 0x12, 0x4b, 0x87, 0x4f, 0x00, 0xaf, 0x95, 0x94, 0x3a, 0xb7, 0x19, 0xf6,
	0x26, 0x71, 0xd4, 0x84, 0x08, 0xce, 0x4a, 0x2a, 0x6d, 0xa3, 0x8c, 0x12, 0x8b, 0xc3, 0x3b, 0x68,
	0xb9, 0xe7, 0xb4, 0xa2, 0x5d, 0x05, 0x8b, 0x4d, 0x35, 0x5b, 0x95, 0x4a, 0xb5, 0xf3, 0xb4, 0xf3,
	0xf5, 0xfd, 0xa1, 0x1a, 0x4e, 0x60, 0xb4, 0x4c, 0x15, 0x95, 0xde, 0x00, 0x01, 0xdc, 0xf7, 0x86,
	0x8b, 0x2c, 0xf7, 0x84, 0xc1, 0x9b, 0xe6, 0x97, 0xb8, 0xf6, 0x86, 0x46, 0xf2, 0xa2, 0x73, 0x62,
	0xef, 0x24, 0x7e, 0x04, 0x77, 0xd9, 0xe8, 0x92, 0x18, 0x1f, 0xc0, 0x59, 0xb3, 0xac, 0xf1, 0xb2,
	0x17, 0xb8, 0xfb, 0xbb, 0xf3, 0x8b, 0x1e, 0x75, 0xa8, 0x1c, 0x0e, 0x3e, 0x5c, 0x7b, 0x9d, 0xdb,
	0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x61, 0x97, 0xb5, 0xe6, 0xb2, 0x01, 0x00, 0x00,
}
