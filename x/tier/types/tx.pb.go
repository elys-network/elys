// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tier/tx.proto

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
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
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

type MsgSetPortfolio struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	User    string `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
}

func (m *MsgSetPortfolio) Reset()         { *m = MsgSetPortfolio{} }
func (m *MsgSetPortfolio) String() string { return proto.CompactTextString(m) }
func (*MsgSetPortfolio) ProtoMessage()    {}
func (*MsgSetPortfolio) Descriptor() ([]byte, []int) {
	return fileDescriptor_30a730293c45319c, []int{0}
}
func (m *MsgSetPortfolio) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPortfolio) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPortfolio.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPortfolio) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPortfolio.Merge(m, src)
}
func (m *MsgSetPortfolio) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPortfolio) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPortfolio.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPortfolio proto.InternalMessageInfo

func (m *MsgSetPortfolio) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgSetPortfolio) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

type MsgSetPortfolioResponse struct {
}

func (m *MsgSetPortfolioResponse) Reset()         { *m = MsgSetPortfolioResponse{} }
func (m *MsgSetPortfolioResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetPortfolioResponse) ProtoMessage()    {}
func (*MsgSetPortfolioResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_30a730293c45319c, []int{1}
}
func (m *MsgSetPortfolioResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPortfolioResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPortfolioResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPortfolioResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPortfolioResponse.Merge(m, src)
}
func (m *MsgSetPortfolioResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPortfolioResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPortfolioResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPortfolioResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSetPortfolio)(nil), "elys.tier.MsgSetPortfolio")
	proto.RegisterType((*MsgSetPortfolioResponse)(nil), "elys.tier.MsgSetPortfolioResponse")
}

func init() { proto.RegisterFile("elys/tier/tx.proto", fileDescriptor_30a730293c45319c) }

var fileDescriptor_30a730293c45319c = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0xcd, 0xa9, 0x2c,
	0xd6, 0x2f, 0xc9, 0x4c, 0x2d, 0xd2, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x04, 0x89, 0xe9, 0x81, 0xc4, 0xa4, 0x24, 0x11, 0xd2, 0x05, 0xf9, 0x45, 0x25, 0x69, 0xf9, 0x39,
	0x99, 0xf9, 0x10, 0x55, 0x52, 0xe2, 0xc9, 0xf9, 0xc5, 0xb9, 0xf9, 0xc5, 0xfa, 0xb9, 0xc5, 0xe9,
	0xfa, 0x65, 0x86, 0x20, 0x0a, 0x2a, 0x21, 0x98, 0x98, 0x9b, 0x99, 0x97, 0xaf, 0x0f, 0x26, 0xa1,
	0x42, 0x92, 0x10, 0xb5, 0xf1, 0x60, 0x9e, 0x3e, 0x84, 0x03, 0x91, 0x52, 0x6a, 0x60, 0xe4, 0xe2,
	0xf7, 0x2d, 0x4e, 0x0f, 0x4e, 0x2d, 0x09, 0x80, 0x59, 0x20, 0x64, 0xc4, 0xc5, 0x9e, 0x5c, 0x94,
	0x9a, 0x58, 0x92, 0x5f, 0x24, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24, 0x71, 0x69, 0x8b, 0xae,
	0x08, 0x54, 0x9b, 0x63, 0x4a, 0x4a, 0x51, 0x6a, 0x71, 0x71, 0x70, 0x49, 0x51, 0x66, 0x5e, 0x7a,
	0x10, 0x4c, 0xa1, 0x90, 0x10, 0x17, 0x4b, 0x69, 0x71, 0x6a, 0x91, 0x04, 0x13, 0x48, 0x43, 0x10,
	0x98, 0x6d, 0xa5, 0xda, 0xf4, 0x7c, 0x83, 0x16, 0x4c, 0x45, 0xd7, 0xf3, 0x0d, 0x5a, 0x22, 0x60,
	0x9f, 0xa0, 0x59, 0xa7, 0x24, 0xc9, 0x25, 0x8e, 0x26, 0x14, 0x94, 0x5a, 0x5c, 0x90, 0x9f, 0x57,
	0x9c, 0x6a, 0x14, 0xc3, 0xc5, 0xec, 0x5b, 0x9c, 0x2e, 0xe4, 0xc7, 0xc5, 0x83, 0xe2, 0x40, 0x29,
	0x3d, 0x78, 0x10, 0xe9, 0xa1, 0x69, 0x95, 0x52, 0xc2, 0x2d, 0x07, 0x33, 0x56, 0x8a, 0xb5, 0xe1,
	0xf9, 0x06, 0x2d, 0x46, 0x27, 0xb7, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0,
	0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88,
	0xd2, 0x49, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0x19, 0xa7, 0x9b,
	0x97, 0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x0d, 0xe6, 0xe8, 0x97, 0x99, 0xeb, 0x57, 0x40, 0xe3, 0xac,
	0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0x1c, 0x94, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6e,
	0xfb, 0x4b, 0x86, 0xcd, 0x01, 0x00, 0x00,
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
	SetPortfolio(ctx context.Context, in *MsgSetPortfolio, opts ...grpc.CallOption) (*MsgSetPortfolioResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetPortfolio(ctx context.Context, in *MsgSetPortfolio, opts ...grpc.CallOption) (*MsgSetPortfolioResponse, error) {
	out := new(MsgSetPortfolioResponse)
	err := c.cc.Invoke(ctx, "/elys.tier.Msg/SetPortfolio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetPortfolio(context.Context, *MsgSetPortfolio) (*MsgSetPortfolioResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SetPortfolio(ctx context.Context, req *MsgSetPortfolio) (*MsgSetPortfolioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPortfolio not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SetPortfolio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPortfolio)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPortfolio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.tier.Msg/SetPortfolio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPortfolio(ctx, req.(*MsgSetPortfolio))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.tier.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetPortfolio",
			Handler:    _Msg_SetPortfolio_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/tier/tx.proto",
}

func (m *MsgSetPortfolio) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPortfolio) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPortfolio) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintTx(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetPortfolioResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPortfolioResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPortfolioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSetPortfolio) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSetPortfolioResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSetPortfolio) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSetPortfolio: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPortfolio: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgSetPortfolioResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSetPortfolioResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPortfolioResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
