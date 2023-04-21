// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/incentive/incentive.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Incentive Info
type IncentiveInfo struct {
	// reward amount
	Amount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount" yaml:"amount"`
	// epoch identifier
	EpochIdentifier string `protobuf:"bytes,2,opt,name=epoch_identifier,json=epochIdentifier,proto3" json:"epoch_identifier,omitempty"`
	// start_time of the distribution
	StartTime time.Time `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time" yaml:"start_time"`
	// distribution duration
	NumEpochs    int64 `protobuf:"varint,4,opt,name=num_epochs,json=numEpochs,proto3" json:"num_epochs,omitempty"`
	CurrentEpoch int64 `protobuf:"varint,5,opt,name=current_epoch,json=currentEpoch,proto3" json:"current_epoch,omitempty"`
}

func (m *IncentiveInfo) Reset()         { *m = IncentiveInfo{} }
func (m *IncentiveInfo) String() string { return proto.CompactTextString(m) }
func (*IncentiveInfo) ProtoMessage()    {}
func (*IncentiveInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed0e67c7f36f3313, []int{0}
}
func (m *IncentiveInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IncentiveInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IncentiveInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IncentiveInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IncentiveInfo.Merge(m, src)
}
func (m *IncentiveInfo) XXX_Size() int {
	return m.Size()
}
func (m *IncentiveInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_IncentiveInfo.DiscardUnknown(m)
}

var xxx_messageInfo_IncentiveInfo proto.InternalMessageInfo

func (m *IncentiveInfo) GetEpochIdentifier() string {
	if m != nil {
		return m.EpochIdentifier
	}
	return ""
}

func (m *IncentiveInfo) GetStartTime() time.Time {
	if m != nil {
		return m.StartTime
	}
	return time.Time{}
}

func (m *IncentiveInfo) GetNumEpochs() int64 {
	if m != nil {
		return m.NumEpochs
	}
	return 0
}

func (m *IncentiveInfo) GetCurrentEpoch() int64 {
	if m != nil {
		return m.CurrentEpoch
	}
	return 0
}

func init() {
	proto.RegisterType((*IncentiveInfo)(nil), "elysnetwork.elys.incentive.IncentiveInfo")
}

func init() { proto.RegisterFile("elys/incentive/incentive.proto", fileDescriptor_ed0e67c7f36f3313) }

var fileDescriptor_ed0e67c7f36f3313 = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x51, 0xcd, 0x4a, 0xeb, 0x40,
	0x14, 0xce, 0xb4, 0xf7, 0x16, 0x3a, 0xf7, 0x16, 0x35, 0xb8, 0x08, 0x81, 0x26, 0x25, 0x82, 0xd4,
	0x45, 0x67, 0x40, 0x77, 0x6e, 0x84, 0x82, 0x60, 0xb6, 0x41, 0x50, 0xdc, 0x94, 0x34, 0x9d, 0xa6,
	0xa1, 0x9d, 0x99, 0x90, 0x99, 0xa8, 0x7d, 0x8b, 0xbe, 0x8a, 0x6f, 0xd1, 0x65, 0x97, 0xe2, 0x22,
	0x4a, 0xfb, 0x06, 0x7d, 0x02, 0x99, 0x49, 0xfa, 0xb3, 0x9a, 0x73, 0xbe, 0xf9, 0xce, 0xf9, 0xce,
	0x77, 0x0e, 0x74, 0xc8, 0x6c, 0x2e, 0x70, 0xc2, 0x22, 0xc2, 0x64, 0xf2, 0x4a, 0x0e, 0x11, 0x4a,
	0x33, 0x2e, 0xb9, 0x69, 0xab, 0x7f, 0x46, 0xe4, 0x1b, 0xcf, 0xa6, 0x48, 0xc5, 0x68, 0xcf, 0xb0,
	0xcf, 0x63, 0x1e, 0x73, 0x4d, 0xc3, 0x2a, 0x2a, 0x2b, 0x6c, 0x37, 0xe6, 0x3c, 0x9e, 0x11, 0xac,
	0xb3, 0x61, 0x3e, 0xc6, 0x32, 0xa1, 0x44, 0xc8, 0x90, 0xa6, 0x25, 0xc1, 0xfb, 0xa8, 0xc1, 0x96,
	0xbf, 0x6b, 0xe2, 0xb3, 0x31, 0x37, 0x9f, 0x60, 0x23, 0xa4, 0x3c, 0x67, 0xd2, 0x02, 0x1d, 0xd0,
	0x6d, 0xf6, 0xef, 0x96, 0x85, 0x6b, 0x7c, 0x15, 0xee, 0x65, 0x9c, 0xc8, 0x49, 0x3e, 0x44, 0x11,
	0xa7, 0x38, 0xe2, 0x82, 0x72, 0x51, 0x3d, 0x3d, 0x31, 0x9a, 0x62, 0x39, 0x4f, 0x89, 0x40, 0x3e,
	0x93, 0xdb, 0xc2, 0x6d, 0xcd, 0x43, 0x3a, 0xbb, 0xf5, 0xca, 0x2e, 0x5e, 0x50, 0xb5, 0x33, 0xaf,
	0xe0, 0x29, 0x49, 0x79, 0x34, 0x19, 0x24, 0x23, 0x25, 0x37, 0x4e, 0x48, 0x66, 0xd5, 0x94, 0x44,
	0x70, 0xa2, 0x71, 0x7f, 0x0f, 0x9b, 0xcf, 0x10, 0x0a, 0x19, 0x66, 0x72, 0xa0, 0xc6, 0xb5, 0xea,
	0x1d, 0xd0, 0xfd, 0x77, 0x6d, 0xa3, 0xd2, 0x0b, 0xda, 0x79, 0x41, 0x8f, 0x3b, 0x2f, 0xfd, 0xb6,
	0x9a, 0x71, 0x5b, 0xb8, 0x67, 0xa5, 0xf2, 0xa1, 0xd6, 0x5b, 0x7c, 0xbb, 0x20, 0x68, 0x6a, 0x40,
	0xd1, 0xcd, 0x36, 0x84, 0x2c, 0xa7, 0x03, 0x2d, 0x28, 0xac, 0x3f, 0x1d, 0xd0, 0xad, 0x07, 0x4d,
	0x96, 0xd3, 0x7b, 0x0d, 0x98, 0x17, 0xb0, 0x15, 0xe5, 0x59, 0x46, 0x98, 0x2c, 0x29, 0xd6, 0x5f,
	0xcd, 0xf8, 0x5f, 0x81, 0x9a, 0xd5, 0x7f, 0x58, 0xae, 0x1d, 0xb0, 0x5a, 0x3b, 0xe0, 0x67, 0xed,
	0x80, 0xc5, 0xc6, 0x31, 0x56, 0x1b, 0xc7, 0xf8, 0xdc, 0x38, 0xc6, 0x0b, 0x3a, 0xda, 0x91, 0xba,
	0x4f, 0xaf, 0x3a, 0x96, 0x4e, 0xf0, 0xfb, 0xd1, 0x69, 0xf5, 0xbe, 0x86, 0x0d, 0xed, 0xe5, 0xe6,
	0x37, 0x00, 0x00, 0xff, 0xff, 0x4c, 0x97, 0xec, 0x4c, 0xf9, 0x01, 0x00, 0x00,
}

func (m *IncentiveInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IncentiveInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IncentiveInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CurrentEpoch != 0 {
		i = encodeVarintIncentive(dAtA, i, uint64(m.CurrentEpoch))
		i--
		dAtA[i] = 0x28
	}
	if m.NumEpochs != 0 {
		i = encodeVarintIncentive(dAtA, i, uint64(m.NumEpochs))
		i--
		dAtA[i] = 0x20
	}
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.StartTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.StartTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintIncentive(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	if len(m.EpochIdentifier) > 0 {
		i -= len(m.EpochIdentifier)
		copy(dAtA[i:], m.EpochIdentifier)
		i = encodeVarintIncentive(dAtA, i, uint64(len(m.EpochIdentifier)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintIncentive(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintIncentive(dAtA []byte, offset int, v uint64) int {
	offset -= sovIncentive(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *IncentiveInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Amount.Size()
	n += 1 + l + sovIncentive(uint64(l))
	l = len(m.EpochIdentifier)
	if l > 0 {
		n += 1 + l + sovIncentive(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.StartTime)
	n += 1 + l + sovIncentive(uint64(l))
	if m.NumEpochs != 0 {
		n += 1 + sovIncentive(uint64(m.NumEpochs))
	}
	if m.CurrentEpoch != 0 {
		n += 1 + sovIncentive(uint64(m.CurrentEpoch))
	}
	return n
}

func sovIncentive(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIncentive(x uint64) (n int) {
	return sovIncentive(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IncentiveInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIncentive
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
			return fmt.Errorf("proto: IncentiveInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IncentiveInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIncentive
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
				return ErrInvalidLengthIncentive
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIncentive
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochIdentifier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIncentive
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
				return ErrInvalidLengthIncentive
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIncentive
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EpochIdentifier = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIncentive
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIncentive
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIncentive
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.StartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumEpochs", wireType)
			}
			m.NumEpochs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIncentive
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumEpochs |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentEpoch", wireType)
			}
			m.CurrentEpoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIncentive
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentEpoch |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIncentive(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIncentive
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
func skipIncentive(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIncentive
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
					return 0, ErrIntOverflowIncentive
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
					return 0, ErrIntOverflowIncentive
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
				return 0, ErrInvalidLengthIncentive
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIncentive
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIncentive
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIncentive        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIncentive          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIncentive = fmt.Errorf("proto: unexpected end of group")
)
