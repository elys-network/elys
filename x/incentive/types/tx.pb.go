// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/incentive/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type MsgSetWithdrawAddress struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	WithdrawAddress  string `protobuf:"bytes,2,opt,name=withdraw_address,json=withdrawAddress,proto3" json:"withdraw_address,omitempty"`
}

func (m *MsgSetWithdrawAddress) Reset()         { *m = MsgSetWithdrawAddress{} }
func (m *MsgSetWithdrawAddress) String() string { return proto.CompactTextString(m) }
func (*MsgSetWithdrawAddress) ProtoMessage()    {}
func (*MsgSetWithdrawAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{0}
}
func (m *MsgSetWithdrawAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetWithdrawAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetWithdrawAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetWithdrawAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetWithdrawAddress.Merge(m, src)
}
func (m *MsgSetWithdrawAddress) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetWithdrawAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetWithdrawAddress.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetWithdrawAddress proto.InternalMessageInfo

func (m *MsgSetWithdrawAddress) GetDelegatorAddress() string {
	if m != nil {
		return m.DelegatorAddress
	}
	return ""
}

func (m *MsgSetWithdrawAddress) GetWithdrawAddress() string {
	if m != nil {
		return m.WithdrawAddress
	}
	return ""
}

type MsgSetWithdrawAddressResponse struct {
}

func (m *MsgSetWithdrawAddressResponse) Reset()         { *m = MsgSetWithdrawAddressResponse{} }
func (m *MsgSetWithdrawAddressResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetWithdrawAddressResponse) ProtoMessage()    {}
func (*MsgSetWithdrawAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{1}
}
func (m *MsgSetWithdrawAddressResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetWithdrawAddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetWithdrawAddressResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetWithdrawAddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetWithdrawAddressResponse.Merge(m, src)
}
func (m *MsgSetWithdrawAddressResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetWithdrawAddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetWithdrawAddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetWithdrawAddressResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSetWithdrawAddress)(nil), "elysnetwork.elys.incentive.MsgSetWithdrawAddress")
	proto.RegisterType((*MsgSetWithdrawAddressResponse)(nil), "elysnetwork.elys.incentive.MsgSetWithdrawAddressResponse")
}

func init() { proto.RegisterFile("elys/incentive/tx.proto", fileDescriptor_59dc3bedfb1cce84) }

var fileDescriptor_59dc3bedfb1cce84 = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0xcd, 0xa9, 0x2c,
	0xd6, 0xcf, 0xcc, 0x4b, 0x4e, 0xcd, 0x2b, 0xc9, 0x2c, 0x4b, 0xd5, 0x2f, 0xa9, 0xd0, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0x92, 0x02, 0x49, 0xe4, 0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0xeb, 0x81,
	0xd8, 0x7a, 0x70, 0x45, 0x4a, 0xf9, 0x5c, 0xa2, 0xbe, 0xc5, 0xe9, 0xc1, 0xa9, 0x25, 0xe1, 0x99,
	0x25, 0x19, 0x29, 0x45, 0x89, 0xe5, 0x8e, 0x29, 0x29, 0x45, 0xa9, 0xc5, 0xc5, 0x42, 0xda, 0x5c,
	0x82, 0x29, 0xa9, 0x39, 0xa9, 0xe9, 0x89, 0x25, 0xf9, 0x45, 0xf1, 0x89, 0x10, 0x41, 0x09, 0x46,
	0x05, 0x46, 0x0d, 0xce, 0x20, 0x01, 0xb8, 0x04, 0x4c, 0xb1, 0x26, 0x97, 0x40, 0x39, 0x54, 0x3f,
	0x5c, 0x2d, 0x13, 0x58, 0x2d, 0x7f, 0x39, 0xaa, 0xb9, 0x4a, 0xf2, 0x5c, 0xb2, 0x58, 0x2d, 0x0c,
	0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x35, 0xea, 0x62, 0xe4, 0x62, 0xf6, 0x2d, 0x4e, 0x17,
	0x6a, 0x62, 0xe4, 0x12, 0xc2, 0xe2, 0x2e, 0x43, 0x3d, 0xdc, 0xbe, 0xd1, 0xc3, 0x6a, 0xb2, 0x94,
	0x25, 0xc9, 0x5a, 0x60, 0x8e, 0x71, 0xf2, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6,
	0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39,
	0x86, 0x28, 0xbd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x90, 0x91,
	0xba, 0x50, 0xf3, 0xc1, 0x1c, 0xfd, 0x0a, 0xe4, 0x78, 0xa8, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03,
	0xc7, 0x85, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x13, 0xc6, 0x5a, 0xdd, 0xa6, 0x01, 0x00, 0x00,
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
	SetWithdrawAddress(ctx context.Context, in *MsgSetWithdrawAddress, opts ...grpc.CallOption) (*MsgSetWithdrawAddressResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetWithdrawAddress(ctx context.Context, in *MsgSetWithdrawAddress, opts ...grpc.CallOption) (*MsgSetWithdrawAddressResponse, error) {
	out := new(MsgSetWithdrawAddressResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.incentive.Msg/SetWithdrawAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetWithdrawAddress(context.Context, *MsgSetWithdrawAddress) (*MsgSetWithdrawAddressResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SetWithdrawAddress(ctx context.Context, req *MsgSetWithdrawAddress) (*MsgSetWithdrawAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetWithdrawAddress not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SetWithdrawAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetWithdrawAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetWithdrawAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.incentive.Msg/SetWithdrawAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetWithdrawAddress(ctx, req.(*MsgSetWithdrawAddress))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elysnetwork.elys.incentive.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetWithdrawAddress",
			Handler:    _Msg_SetWithdrawAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/incentive/tx.proto",
}

func (m *MsgSetWithdrawAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetWithdrawAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetWithdrawAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.WithdrawAddress) > 0 {
		i -= len(m.WithdrawAddress)
		copy(dAtA[i:], m.WithdrawAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.WithdrawAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetWithdrawAddressResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetWithdrawAddressResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetWithdrawAddressResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgSetWithdrawAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.WithdrawAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSetWithdrawAddressResponse) Size() (n int) {
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
func (m *MsgSetWithdrawAddress) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetWithdrawAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetWithdrawAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
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
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawAddress", wireType)
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
			m.WithdrawAddress = string(dAtA[iNdEx:postIndex])
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
func (m *MsgSetWithdrawAddressResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetWithdrawAddressResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetWithdrawAddressResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
