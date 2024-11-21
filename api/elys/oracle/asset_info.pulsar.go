// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package oracle

import (
	fmt "fmt"
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
	md_AssetInfo             protoreflect.MessageDescriptor
	fd_AssetInfo_denom       protoreflect.FieldDescriptor
	fd_AssetInfo_display     protoreflect.FieldDescriptor
	fd_AssetInfo_band_ticker protoreflect.FieldDescriptor
	fd_AssetInfo_elys_ticker protoreflect.FieldDescriptor
	fd_AssetInfo_decimal     protoreflect.FieldDescriptor
)

func init() {
	file_elys_oracle_asset_info_proto_init()
	md_AssetInfo = File_elys_oracle_asset_info_proto.Messages().ByName("AssetInfo")
	fd_AssetInfo_denom = md_AssetInfo.Fields().ByName("denom")
	fd_AssetInfo_display = md_AssetInfo.Fields().ByName("display")
	fd_AssetInfo_band_ticker = md_AssetInfo.Fields().ByName("band_ticker")
	fd_AssetInfo_elys_ticker = md_AssetInfo.Fields().ByName("elys_ticker")
	fd_AssetInfo_decimal = md_AssetInfo.Fields().ByName("decimal")
}

var _ protoreflect.Message = (*fastReflection_AssetInfo)(nil)

type fastReflection_AssetInfo AssetInfo

func (x *AssetInfo) ProtoReflect() protoreflect.Message {
	return (*fastReflection_AssetInfo)(x)
}

func (x *AssetInfo) slowProtoReflect() protoreflect.Message {
	mi := &file_elys_oracle_asset_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_AssetInfo_messageType fastReflection_AssetInfo_messageType
var _ protoreflect.MessageType = fastReflection_AssetInfo_messageType{}

type fastReflection_AssetInfo_messageType struct{}

func (x fastReflection_AssetInfo_messageType) Zero() protoreflect.Message {
	return (*fastReflection_AssetInfo)(nil)
}
func (x fastReflection_AssetInfo_messageType) New() protoreflect.Message {
	return new(fastReflection_AssetInfo)
}
func (x fastReflection_AssetInfo_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_AssetInfo
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_AssetInfo) Descriptor() protoreflect.MessageDescriptor {
	return md_AssetInfo
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_AssetInfo) Type() protoreflect.MessageType {
	return _fastReflection_AssetInfo_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_AssetInfo) New() protoreflect.Message {
	return new(fastReflection_AssetInfo)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_AssetInfo) Interface() protoreflect.ProtoMessage {
	return (*AssetInfo)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_AssetInfo) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.Denom != "" {
		value := protoreflect.ValueOfString(x.Denom)
		if !f(fd_AssetInfo_denom, value) {
			return
		}
	}
	if x.Display != "" {
		value := protoreflect.ValueOfString(x.Display)
		if !f(fd_AssetInfo_display, value) {
			return
		}
	}
	if x.BandTicker != "" {
		value := protoreflect.ValueOfString(x.BandTicker)
		if !f(fd_AssetInfo_band_ticker, value) {
			return
		}
	}
	if x.ElysTicker != "" {
		value := protoreflect.ValueOfString(x.ElysTicker)
		if !f(fd_AssetInfo_elys_ticker, value) {
			return
		}
	}
	if x.Decimal != uint64(0) {
		value := protoreflect.ValueOfUint64(x.Decimal)
		if !f(fd_AssetInfo_decimal, value) {
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
func (x *fastReflection_AssetInfo) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "elys.oracle.AssetInfo.denom":
		return x.Denom != ""
	case "elys.oracle.AssetInfo.display":
		return x.Display != ""
	case "elys.oracle.AssetInfo.band_ticker":
		return x.BandTicker != ""
	case "elys.oracle.AssetInfo.elys_ticker":
		return x.ElysTicker != ""
	case "elys.oracle.AssetInfo.decimal":
		return x.Decimal != uint64(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_AssetInfo) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "elys.oracle.AssetInfo.denom":
		x.Denom = ""
	case "elys.oracle.AssetInfo.display":
		x.Display = ""
	case "elys.oracle.AssetInfo.band_ticker":
		x.BandTicker = ""
	case "elys.oracle.AssetInfo.elys_ticker":
		x.ElysTicker = ""
	case "elys.oracle.AssetInfo.decimal":
		x.Decimal = uint64(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_AssetInfo) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "elys.oracle.AssetInfo.denom":
		value := x.Denom
		return protoreflect.ValueOfString(value)
	case "elys.oracle.AssetInfo.display":
		value := x.Display
		return protoreflect.ValueOfString(value)
	case "elys.oracle.AssetInfo.band_ticker":
		value := x.BandTicker
		return protoreflect.ValueOfString(value)
	case "elys.oracle.AssetInfo.elys_ticker":
		value := x.ElysTicker
		return protoreflect.ValueOfString(value)
	case "elys.oracle.AssetInfo.decimal":
		value := x.Decimal
		return protoreflect.ValueOfUint64(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", descriptor.FullName()))
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
func (x *fastReflection_AssetInfo) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "elys.oracle.AssetInfo.denom":
		x.Denom = value.Interface().(string)
	case "elys.oracle.AssetInfo.display":
		x.Display = value.Interface().(string)
	case "elys.oracle.AssetInfo.band_ticker":
		x.BandTicker = value.Interface().(string)
	case "elys.oracle.AssetInfo.elys_ticker":
		x.ElysTicker = value.Interface().(string)
	case "elys.oracle.AssetInfo.decimal":
		x.Decimal = value.Uint()
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", fd.FullName()))
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
func (x *fastReflection_AssetInfo) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.oracle.AssetInfo.denom":
		panic(fmt.Errorf("field denom of message elys.oracle.AssetInfo is not mutable"))
	case "elys.oracle.AssetInfo.display":
		panic(fmt.Errorf("field display of message elys.oracle.AssetInfo is not mutable"))
	case "elys.oracle.AssetInfo.band_ticker":
		panic(fmt.Errorf("field band_ticker of message elys.oracle.AssetInfo is not mutable"))
	case "elys.oracle.AssetInfo.elys_ticker":
		panic(fmt.Errorf("field elys_ticker of message elys.oracle.AssetInfo is not mutable"))
	case "elys.oracle.AssetInfo.decimal":
		panic(fmt.Errorf("field decimal of message elys.oracle.AssetInfo is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_AssetInfo) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.oracle.AssetInfo.denom":
		return protoreflect.ValueOfString("")
	case "elys.oracle.AssetInfo.display":
		return protoreflect.ValueOfString("")
	case "elys.oracle.AssetInfo.band_ticker":
		return protoreflect.ValueOfString("")
	case "elys.oracle.AssetInfo.elys_ticker":
		return protoreflect.ValueOfString("")
	case "elys.oracle.AssetInfo.decimal":
		return protoreflect.ValueOfUint64(uint64(0))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.oracle.AssetInfo"))
		}
		panic(fmt.Errorf("message elys.oracle.AssetInfo does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_AssetInfo) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in elys.oracle.AssetInfo", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_AssetInfo) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_AssetInfo) SetUnknown(fields protoreflect.RawFields) {
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
func (x *fastReflection_AssetInfo) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_AssetInfo) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*AssetInfo)
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
		l = len(x.Denom)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Display)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.BandTicker)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.ElysTicker)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.Decimal != 0 {
			n += 1 + runtime.Sov(uint64(x.Decimal))
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
		x := input.Message.Interface().(*AssetInfo)
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
		if x.Decimal != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.Decimal))
			i--
			dAtA[i] = 0x28
		}
		if len(x.ElysTicker) > 0 {
			i -= len(x.ElysTicker)
			copy(dAtA[i:], x.ElysTicker)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.ElysTicker)))
			i--
			dAtA[i] = 0x22
		}
		if len(x.BandTicker) > 0 {
			i -= len(x.BandTicker)
			copy(dAtA[i:], x.BandTicker)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.BandTicker)))
			i--
			dAtA[i] = 0x1a
		}
		if len(x.Display) > 0 {
			i -= len(x.Display)
			copy(dAtA[i:], x.Display)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Display)))
			i--
			dAtA[i] = 0x12
		}
		if len(x.Denom) > 0 {
			i -= len(x.Denom)
			copy(dAtA[i:], x.Denom)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Denom)))
			i--
			dAtA[i] = 0xa
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
		x := input.Message.Interface().(*AssetInfo)
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: AssetInfo: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: AssetInfo: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
				x.Denom = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Display", wireType)
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
				x.Display = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field BandTicker", wireType)
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
				x.BandTicker = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field ElysTicker", wireType)
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
				x.ElysTicker = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 5:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Decimal", wireType)
				}
				x.Decimal = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.Decimal |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
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
// source: elys/oracle/asset_info.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AssetInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Denom      string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Display    string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	BandTicker string `protobuf:"bytes,3,opt,name=band_ticker,json=bandTicker,proto3" json:"band_ticker,omitempty"`
	ElysTicker string `protobuf:"bytes,4,opt,name=elys_ticker,json=elysTicker,proto3" json:"elys_ticker,omitempty"`
	Decimal    uint64 `protobuf:"varint,5,opt,name=decimal,proto3" json:"decimal,omitempty"`
}

func (x *AssetInfo) Reset() {
	*x = AssetInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_elys_oracle_asset_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInfo) ProtoMessage() {}

// Deprecated: Use AssetInfo.ProtoReflect.Descriptor instead.
func (*AssetInfo) Descriptor() ([]byte, []int) {
	return file_elys_oracle_asset_info_proto_rawDescGZIP(), []int{0}
}

func (x *AssetInfo) GetDenom() string {
	if x != nil {
		return x.Denom
	}
	return ""
}

func (x *AssetInfo) GetDisplay() string {
	if x != nil {
		return x.Display
	}
	return ""
}

func (x *AssetInfo) GetBandTicker() string {
	if x != nil {
		return x.BandTicker
	}
	return ""
}

func (x *AssetInfo) GetElysTicker() string {
	if x != nil {
		return x.ElysTicker
	}
	return ""
}

func (x *AssetInfo) GetDecimal() uint64 {
	if x != nil {
		return x.Decimal
	}
	return 0
}

var File_elys_oracle_asset_info_proto protoreflect.FileDescriptor

var file_elys_oracle_asset_info_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x2f, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x65, 0x6c, 0x79, 0x73, 0x2e, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x1a, 0x14, 0x67, 0x6f, 0x67,
	0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x97, 0x01, 0x0a, 0x09, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x14, 0x0a, 0x05, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x12,
	0x1f, 0x0a, 0x0b, 0x62, 0x61, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x61, 0x6e, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6c, 0x79, 0x73, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6c, 0x79, 0x73, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x72, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x42, 0x9c, 0x01, 0x0a, 0x0f,
	0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x79, 0x73, 0x2e, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x42,
	0x0e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c,
	0x79, 0x73, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0xa2,
	0x02, 0x03, 0x45, 0x4f, 0x58, 0xaa, 0x02, 0x0b, 0x45, 0x6c, 0x79, 0x73, 0x2e, 0x4f, 0x72, 0x61,
	0x63, 0x6c, 0x65, 0xca, 0x02, 0x0b, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x4f, 0x72, 0x61, 0x63, 0x6c,
	0x65, 0xe2, 0x02, 0x17, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x45, 0x6c,
	0x79, 0x73, 0x3a, 0x3a, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_elys_oracle_asset_info_proto_rawDescOnce sync.Once
	file_elys_oracle_asset_info_proto_rawDescData = file_elys_oracle_asset_info_proto_rawDesc
)

func file_elys_oracle_asset_info_proto_rawDescGZIP() []byte {
	file_elys_oracle_asset_info_proto_rawDescOnce.Do(func() {
		file_elys_oracle_asset_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_elys_oracle_asset_info_proto_rawDescData)
	})
	return file_elys_oracle_asset_info_proto_rawDescData
}

var file_elys_oracle_asset_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_elys_oracle_asset_info_proto_goTypes = []interface{}{
	(*AssetInfo)(nil), // 0: elys.oracle.AssetInfo
}
var file_elys_oracle_asset_info_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_elys_oracle_asset_info_proto_init() }
func file_elys_oracle_asset_info_proto_init() {
	if File_elys_oracle_asset_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_elys_oracle_asset_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInfo); i {
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
			RawDescriptor: file_elys_oracle_asset_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_elys_oracle_asset_info_proto_goTypes,
		DependencyIndexes: file_elys_oracle_asset_info_proto_depIdxs,
		MessageInfos:      file_elys_oracle_asset_info_proto_msgTypes,
	}.Build()
	File_elys_oracle_asset_info_proto = out.File
	file_elys_oracle_asset_info_proto_rawDesc = nil
	file_elys_oracle_asset_info_proto_goTypes = nil
	file_elys_oracle_asset_info_proto_depIdxs = nil
}
