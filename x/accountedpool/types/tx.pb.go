// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/accountedpool/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("elys/accountedpool/tx.proto", fileDescriptor_8792482078796b08) }

var fileDescriptor_8792482078796b08 = []byte{
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0xcd, 0xa9, 0x2c,
	0xd6, 0x4f, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0x49, 0x4d, 0x29, 0xc8, 0xcf, 0xcf, 0xd1, 0x2f,
	0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x02, 0x49, 0xea, 0xa1, 0x48, 0x4a, 0x89,
	0x27, 0xe7, 0x17, 0xe7, 0xe6, 0x17, 0xeb, 0xe7, 0x16, 0xa7, 0xeb, 0x97, 0x19, 0x82, 0x28, 0x88,
	0x62, 0x29, 0x49, 0x88, 0x44, 0x3c, 0x98, 0xa7, 0x0f, 0xe1, 0x40, 0xa5, 0x04, 0x13, 0x73, 0x33,
	0xf3, 0xf2, 0xf5, 0xc1, 0x24, 0x44, 0xc8, 0x88, 0x87, 0x8b, 0xd9, 0xb7, 0x38, 0x5d, 0x8a, 0xb5,
	0xe1, 0xf9, 0x06, 0x2d, 0x46, 0x27, 0x9f, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c,
	0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63,
	0x88, 0x32, 0x4a, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0xb9, 0x46,
	0x37, 0x2f, 0xb5, 0xa4, 0x3c, 0xbf, 0x28, 0x1b, 0xcc, 0xd1, 0xaf, 0x40, 0x77, 0x79, 0x65, 0x41,
	0x6a, 0x71, 0x12, 0x1b, 0xd8, 0x0a, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x06, 0xe0, 0xab,
	0x03, 0xdc, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.accountedpool.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "elys/accountedpool/tx.proto",
}
