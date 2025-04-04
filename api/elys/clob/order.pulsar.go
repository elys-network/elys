// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package clob

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_PerpetualOrder              protoreflect.MessageDescriptor
	fd_PerpetualOrder_market_id    protoreflect.FieldDescriptor
	fd_PerpetualOrder_order_type   protoreflect.FieldDescriptor
	fd_PerpetualOrder_price        protoreflect.FieldDescriptor
	fd_PerpetualOrder_block_height protoreflect.FieldDescriptor
	fd_PerpetualOrder_owner        protoreflect.FieldDescriptor
	fd_PerpetualOrder_amount       protoreflect.FieldDescriptor
	fd_PerpetualOrder_filled       protoreflect.FieldDescriptor
)

func init() {
	file_elys_clob_order_proto_init()
	md_PerpetualOrder = File_elys_clob_order_proto.Messages().ByName("PerpetualOrder")
	fd_PerpetualOrder_market_id = md_PerpetualOrder.Fields().ByName("market_id")
	fd_PerpetualOrder_order_type = md_PerpetualOrder.Fields().ByName("order_type")
	fd_PerpetualOrder_price = md_PerpetualOrder.Fields().ByName("price")
	fd_PerpetualOrder_block_height = md_PerpetualOrder.Fields().ByName("block_height")
	fd_PerpetualOrder_owner = md_PerpetualOrder.Fields().ByName("owner")
	fd_PerpetualOrder_amount = md_PerpetualOrder.Fields().ByName("amount")
	fd_PerpetualOrder_filled = md_PerpetualOrder.Fields().ByName("filled")
}

var _ protoreflect.Message = (*fastReflection_PerpetualOrder)(nil)

type fastReflection_PerpetualOrder PerpetualOrder

func (x *PerpetualOrder) ProtoReflect() protoreflect.Message {
	return (*fastReflection_PerpetualOrder)(x)
}

func (x *PerpetualOrder) slowProtoReflect() protoreflect.Message {
	mi := &file_elys_clob_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_PerpetualOrder_messageType fastReflection_PerpetualOrder_messageType
var _ protoreflect.MessageType = fastReflection_PerpetualOrder_messageType{}

type fastReflection_PerpetualOrder_messageType struct{}

func (x fastReflection_PerpetualOrder_messageType) Zero() protoreflect.Message {
	return (*fastReflection_PerpetualOrder)(nil)
}
func (x fastReflection_PerpetualOrder_messageType) New() protoreflect.Message {
	return new(fastReflection_PerpetualOrder)
}
func (x fastReflection_PerpetualOrder_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_PerpetualOrder
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_PerpetualOrder) Descriptor() protoreflect.MessageDescriptor {
	return md_PerpetualOrder
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_PerpetualOrder) Type() protoreflect.MessageType {
	return _fastReflection_PerpetualOrder_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_PerpetualOrder) New() protoreflect.Message {
	return new(fastReflection_PerpetualOrder)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_PerpetualOrder) Interface() protoreflect.ProtoMessage {
	return (*PerpetualOrder)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_PerpetualOrder) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.MarketId != uint64(0) {
		value := protoreflect.ValueOfUint64(x.MarketId)
		if !f(fd_PerpetualOrder_market_id, value) {
			return
		}
	}
	if x.OrderType != 0 {
		value := protoreflect.ValueOfEnum((protoreflect.EnumNumber)(x.OrderType))
		if !f(fd_PerpetualOrder_order_type, value) {
			return
		}
	}
	if x.Price != "" {
		value := protoreflect.ValueOfString(x.Price)
		if !f(fd_PerpetualOrder_price, value) {
			return
		}
	}
	if x.BlockHeight != uint64(0) {
		value := protoreflect.ValueOfUint64(x.BlockHeight)
		if !f(fd_PerpetualOrder_block_height, value) {
			return
		}
	}
	if x.Owner != "" {
		value := protoreflect.ValueOfString(x.Owner)
		if !f(fd_PerpetualOrder_owner, value) {
			return
		}
	}
	if x.Amount != "" {
		value := protoreflect.ValueOfString(x.Amount)
		if !f(fd_PerpetualOrder_amount, value) {
			return
		}
	}
	if x.Filled != "" {
		value := protoreflect.ValueOfString(x.Filled)
		if !f(fd_PerpetualOrder_filled, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_PerpetualOrder) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		return x.MarketId != uint64(0)
	case "elys.clob.PerpetualOrder.order_type":
		return x.OrderType != 0
	case "elys.clob.PerpetualOrder.price":
		return x.Price != ""
	case "elys.clob.PerpetualOrder.block_height":
		return x.BlockHeight != uint64(0)
	case "elys.clob.PerpetualOrder.owner":
		return x.Owner != ""
	case "elys.clob.PerpetualOrder.amount":
		return x.Amount != ""
	case "elys.clob.PerpetualOrder.filled":
		return x.Filled != ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_PerpetualOrder) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		x.MarketId = uint64(0)
	case "elys.clob.PerpetualOrder.order_type":
		x.OrderType = 0
	case "elys.clob.PerpetualOrder.price":
		x.Price = ""
	case "elys.clob.PerpetualOrder.block_height":
		x.BlockHeight = uint64(0)
	case "elys.clob.PerpetualOrder.owner":
		x.Owner = ""
	case "elys.clob.PerpetualOrder.amount":
		x.Amount = ""
	case "elys.clob.PerpetualOrder.filled":
		x.Filled = ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_PerpetualOrder) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		value := x.MarketId
		return protoreflect.ValueOfUint64(value)
	case "elys.clob.PerpetualOrder.order_type":
		value := x.OrderType
		return protoreflect.ValueOfEnum((protoreflect.EnumNumber)(value))
	case "elys.clob.PerpetualOrder.price":
		value := x.Price
		return protoreflect.ValueOfString(value)
	case "elys.clob.PerpetualOrder.block_height":
		value := x.BlockHeight
		return protoreflect.ValueOfUint64(value)
	case "elys.clob.PerpetualOrder.owner":
		value := x.Owner
		return protoreflect.ValueOfString(value)
	case "elys.clob.PerpetualOrder.amount":
		value := x.Amount
		return protoreflect.ValueOfString(value)
	case "elys.clob.PerpetualOrder.filled":
		value := x.Filled
		return protoreflect.ValueOfString(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_PerpetualOrder) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		x.MarketId = value.Uint()
	case "elys.clob.PerpetualOrder.order_type":
		x.OrderType = (OrderType)(value.Enum())
	case "elys.clob.PerpetualOrder.price":
		x.Price = value.Interface().(string)
	case "elys.clob.PerpetualOrder.block_height":
		x.BlockHeight = value.Uint()
	case "elys.clob.PerpetualOrder.owner":
		x.Owner = value.Interface().(string)
	case "elys.clob.PerpetualOrder.amount":
		x.Amount = value.Interface().(string)
	case "elys.clob.PerpetualOrder.filled":
		x.Filled = value.Interface().(string)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_PerpetualOrder) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		panic(fmt.Errorf("field market_id of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.order_type":
		panic(fmt.Errorf("field order_type of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.price":
		panic(fmt.Errorf("field price of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.block_height":
		panic(fmt.Errorf("field block_height of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.owner":
		panic(fmt.Errorf("field owner of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.amount":
		panic(fmt.Errorf("field amount of message elys.clob.PerpetualOrder is not mutable"))
	case "elys.clob.PerpetualOrder.filled":
		panic(fmt.Errorf("field filled of message elys.clob.PerpetualOrder is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_PerpetualOrder) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.clob.PerpetualOrder.market_id":
		return protoreflect.ValueOfUint64(uint64(0))
	case "elys.clob.PerpetualOrder.order_type":
		return protoreflect.ValueOfEnum(0)
	case "elys.clob.PerpetualOrder.price":
		return protoreflect.ValueOfString("")
	case "elys.clob.PerpetualOrder.block_height":
		return protoreflect.ValueOfUint64(uint64(0))
	case "elys.clob.PerpetualOrder.owner":
		return protoreflect.ValueOfString("")
	case "elys.clob.PerpetualOrder.amount":
		return protoreflect.ValueOfString("")
	case "elys.clob.PerpetualOrder.filled":
		return protoreflect.ValueOfString("")
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.clob.PerpetualOrder"))
		}
		panic(fmt.Errorf("message elys.clob.PerpetualOrder does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_PerpetualOrder) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in elys.clob.PerpetualOrder", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_PerpetualOrder) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_PerpetualOrder) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_PerpetualOrder) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_PerpetualOrder) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*PerpetualOrder)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		if x.MarketId != 0 {
			n += 1 + runtime.Sov(uint64(x.MarketId))
		}
		if x.OrderType != 0 {
			n += 1 + runtime.Sov(uint64(x.OrderType))
		}
		l = len(x.Price)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.BlockHeight != 0 {
			n += 1 + runtime.Sov(uint64(x.BlockHeight))
		}
		l = len(x.Owner)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Amount)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Filled)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*PerpetualOrder)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if len(x.Filled) > 0 {
			i -= len(x.Filled)
			copy(dAtA[i:], x.Filled)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Filled)))
			i--
			dAtA[i] = 0x3a
		}
		if len(x.Amount) > 0 {
			i -= len(x.Amount)
			copy(dAtA[i:], x.Amount)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Amount)))
			i--
			dAtA[i] = 0x32
		}
		if len(x.Owner) > 0 {
			i -= len(x.Owner)
			copy(dAtA[i:], x.Owner)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Owner)))
			i--
			dAtA[i] = 0x2a
		}
		if x.BlockHeight != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.BlockHeight))
			i--
			dAtA[i] = 0x20
		}
		if len(x.Price) > 0 {
			i -= len(x.Price)
			copy(dAtA[i:], x.Price)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Price)))
			i--
			dAtA[i] = 0x1a
		}
		if x.OrderType != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.OrderType))
			i--
			dAtA[i] = 0x10
		}
		if x.MarketId != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.MarketId))
			i--
			dAtA[i] = 0x8
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*PerpetualOrder)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: PerpetualOrder: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: PerpetualOrder: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field MarketId", wireType)
				}
				x.MarketId = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.MarketId |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 2:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field OrderType", wireType)
				}
				x.OrderType = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.OrderType |= OrderType(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Price = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 4:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
				}
				x.BlockHeight = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.BlockHeight |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 5:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Owner = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 6:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Amount = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 7:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Filled", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Filled = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: elys/clob/order.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderType int32

const (
	OrderType_ORDER_TYPE_UNSPECIFIED OrderType = 0
	OrderType_ORDER_TYPE_LIMIT_BUY   OrderType = 1
	OrderType_ORDER_TYPE_LIMIT_SELL  OrderType = 2
	OrderType_ORDER_TYPE_MARKET_BUY  OrderType = 3
	OrderType_ORDER_TYPE_MARKET_SELL OrderType = 4
)

// Enum value maps for OrderType.
var (
	OrderType_name = map[int32]string{
		0: "ORDER_TYPE_UNSPECIFIED",
		1: "ORDER_TYPE_LIMIT_BUY",
		2: "ORDER_TYPE_LIMIT_SELL",
		3: "ORDER_TYPE_MARKET_BUY",
		4: "ORDER_TYPE_MARKET_SELL",
	}
	OrderType_value = map[string]int32{
		"ORDER_TYPE_UNSPECIFIED": 0,
		"ORDER_TYPE_LIMIT_BUY":   1,
		"ORDER_TYPE_LIMIT_SELL":  2,
		"ORDER_TYPE_MARKET_BUY":  3,
		"ORDER_TYPE_MARKET_SELL": 4,
	}
)

func (x OrderType) Enum() *OrderType {
	p := new(OrderType)
	*p = x
	return p
}

func (x OrderType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderType) Descriptor() protoreflect.EnumDescriptor {
	return file_elys_clob_order_proto_enumTypes[0].Descriptor()
}

func (OrderType) Type() protoreflect.EnumType {
	return &file_elys_clob_order_proto_enumTypes[0]
}

func (x OrderType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderType.Descriptor instead.
func (OrderType) EnumDescriptor() ([]byte, []int) {
	return file_elys_clob_order_proto_rawDescGZIP(), []int{0}
}

// key = market_id + is_long + price + block_height
type PerpetualOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MarketId    uint64    `protobuf:"varint,1,opt,name=market_id,json=marketId,proto3" json:"market_id,omitempty"`
	OrderType   OrderType `protobuf:"varint,2,opt,name=order_type,json=orderType,proto3,enum=elys.clob.OrderType" json:"order_type,omitempty"`
	Price       string    `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	BlockHeight uint64    `protobuf:"varint,4,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	Owner       string    `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Amount      string    `protobuf:"bytes,6,opt,name=amount,proto3" json:"amount,omitempty"`
	Filled      string    `protobuf:"bytes,7,opt,name=filled,proto3" json:"filled,omitempty"`
}

func (x *PerpetualOrder) Reset() {
	*x = PerpetualOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_elys_clob_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PerpetualOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PerpetualOrder) ProtoMessage() {}

// Deprecated: Use PerpetualOrder.ProtoReflect.Descriptor instead.
func (*PerpetualOrder) Descriptor() ([]byte, []int) {
	return file_elys_clob_order_proto_rawDescGZIP(), []int{0}
}

func (x *PerpetualOrder) GetMarketId() uint64 {
	if x != nil {
		return x.MarketId
	}
	return 0
}

func (x *PerpetualOrder) GetOrderType() OrderType {
	if x != nil {
		return x.OrderType
	}
	return OrderType_ORDER_TYPE_UNSPECIFIED
}

func (x *PerpetualOrder) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *PerpetualOrder) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *PerpetualOrder) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *PerpetualOrder) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

func (x *PerpetualOrder) GetFilled() string {
	if x != nil {
		return x.Filled
	}
	return ""
}

var File_elys_clob_order_proto protoreflect.FileDescriptor

var file_elys_clob_order_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x62, 0x2f, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x6c, 0x79, 0x73, 0x2e, 0x63, 0x6c,
	0x6f, 0x62, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f,
	0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x88, 0x03, 0x0a, 0x0e, 0x50, 0x65, 0x72, 0x70, 0x65, 0x74, 0x75, 0x61,
	0x6c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6d, 0x61, 0x72, 0x6b, 0x65,
	0x74, 0x49, 0x64, 0x12, 0x33, 0x0a, 0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x65, 0x6c, 0x79, 0x73, 0x2e, 0x63,
	0x6c, 0x6f, 0x62, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x47, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x31, 0xc8, 0xde, 0x1f, 0x00, 0xda, 0xde, 0x1f,
	0x1b, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x6d, 0x61,
	0x74, 0x68, 0x2e, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x44, 0x65, 0x63, 0xd2, 0xb4, 0x2d, 0x0a,
	0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x63, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x18, 0xd2, 0xb4, 0x2d, 0x14, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x12, 0x43, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x2b, 0xc8, 0xde, 0x1f, 0x00, 0xda, 0xde, 0x1f, 0x15, 0x63, 0x6f,
	0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x6d, 0x61, 0x74, 0x68, 0x2e,
	0x49, 0x6e, 0x74, 0xd2, 0xb4, 0x2d, 0x0a, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x49, 0x6e,
	0x74, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x43, 0x0a, 0x06, 0x66, 0x69, 0x6c,
	0x6c, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2b, 0xc8, 0xde, 0x1f, 0x00, 0xda,
	0xde, 0x1f, 0x15, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f,
	0x6d, 0x61, 0x74, 0x68, 0x2e, 0x49, 0x6e, 0x74, 0xd2, 0xb4, 0x2d, 0x0a, 0x63, 0x6f, 0x73, 0x6d,
	0x6f, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x2a, 0x93,
	0x01, 0x0a, 0x09, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x4f, 0x52, 0x44, 0x45,
	0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x5f, 0x42, 0x55, 0x59,
	0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x10, 0x02, 0x12, 0x19, 0x0a,
	0x15, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41, 0x52, 0x4b,
	0x45, 0x54, 0x5f, 0x42, 0x55, 0x59, 0x10, 0x03, 0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x52, 0x44, 0x45,
	0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x5f, 0x53, 0x45,
	0x4c, 0x4c, 0x10, 0x04, 0x42, 0x8c, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x79,
	0x73, 0x2e, 0x63, 0x6c, 0x6f, 0x62, 0x42, 0x0a, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x65, 0x6c,
	0x79, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x62,
	0xa2, 0x02, 0x03, 0x45, 0x43, 0x58, 0xaa, 0x02, 0x09, 0x45, 0x6c, 0x79, 0x73, 0x2e, 0x43, 0x6c,
	0x6f, 0x62, 0xca, 0x02, 0x09, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x43, 0x6c, 0x6f, 0x62, 0xe2, 0x02,
	0x15, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x43, 0x6c, 0x6f, 0x62, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x45, 0x6c, 0x79, 0x73, 0x3a, 0x3a, 0x43,
	0x6c, 0x6f, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_elys_clob_order_proto_rawDescOnce sync.Once
	file_elys_clob_order_proto_rawDescData = file_elys_clob_order_proto_rawDesc
)

func file_elys_clob_order_proto_rawDescGZIP() []byte {
	file_elys_clob_order_proto_rawDescOnce.Do(func() {
		file_elys_clob_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_elys_clob_order_proto_rawDescData)
	})
	return file_elys_clob_order_proto_rawDescData
}

var file_elys_clob_order_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_elys_clob_order_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_elys_clob_order_proto_goTypes = []interface{}{
	(OrderType)(0),         // 0: elys.clob.OrderType
	(*PerpetualOrder)(nil), // 1: elys.clob.PerpetualOrder
}
var file_elys_clob_order_proto_depIdxs = []int32{
	0, // 0: elys.clob.PerpetualOrder.order_type:type_name -> elys.clob.OrderType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_elys_clob_order_proto_init() }
func file_elys_clob_order_proto_init() {
	if File_elys_clob_order_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_elys_clob_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PerpetualOrder); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_elys_clob_order_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_elys_clob_order_proto_goTypes,
		DependencyIndexes: file_elys_clob_order_proto_depIdxs,
		EnumInfos:         file_elys_clob_order_proto_enumTypes,
		MessageInfos:      file_elys_clob_order_proto_msgTypes,
	}.Build()
	File_elys_clob_order_proto = out.File
	file_elys_clob_order_proto_rawDesc = nil
	file_elys_clob_order_proto_goTypes = nil
	file_elys_clob_order_proto_depIdxs = nil
}
