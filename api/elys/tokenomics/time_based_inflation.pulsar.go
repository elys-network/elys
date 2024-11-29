// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package tokenomics

import (
	fmt "fmt"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_TimeBasedInflation                    protoreflect.MessageDescriptor
	fd_TimeBasedInflation_start_block_height protoreflect.FieldDescriptor
	fd_TimeBasedInflation_end_block_height   protoreflect.FieldDescriptor
	fd_TimeBasedInflation_description        protoreflect.FieldDescriptor
	fd_TimeBasedInflation_inflation          protoreflect.FieldDescriptor
	fd_TimeBasedInflation_authority          protoreflect.FieldDescriptor
)

func init() {
	file_elys_tokenomics_time_based_inflation_proto_init()
	md_TimeBasedInflation = File_elys_tokenomics_time_based_inflation_proto.Messages().ByName("TimeBasedInflation")
	fd_TimeBasedInflation_start_block_height = md_TimeBasedInflation.Fields().ByName("start_block_height")
	fd_TimeBasedInflation_end_block_height = md_TimeBasedInflation.Fields().ByName("end_block_height")
	fd_TimeBasedInflation_description = md_TimeBasedInflation.Fields().ByName("description")
	fd_TimeBasedInflation_inflation = md_TimeBasedInflation.Fields().ByName("inflation")
	fd_TimeBasedInflation_authority = md_TimeBasedInflation.Fields().ByName("authority")
}

var _ protoreflect.Message = (*fastReflection_TimeBasedInflation)(nil)

type fastReflection_TimeBasedInflation TimeBasedInflation

func (x *TimeBasedInflation) ProtoReflect() protoreflect.Message {
	return (*fastReflection_TimeBasedInflation)(x)
}

func (x *TimeBasedInflation) slowProtoReflect() protoreflect.Message {
	mi := &file_elys_tokenomics_time_based_inflation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_TimeBasedInflation_messageType fastReflection_TimeBasedInflation_messageType
var _ protoreflect.MessageType = fastReflection_TimeBasedInflation_messageType{}

type fastReflection_TimeBasedInflation_messageType struct{}

func (x fastReflection_TimeBasedInflation_messageType) Zero() protoreflect.Message {
	return (*fastReflection_TimeBasedInflation)(nil)
}
func (x fastReflection_TimeBasedInflation_messageType) New() protoreflect.Message {
	return new(fastReflection_TimeBasedInflation)
}
func (x fastReflection_TimeBasedInflation_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_TimeBasedInflation
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_TimeBasedInflation) Descriptor() protoreflect.MessageDescriptor {
	return md_TimeBasedInflation
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_TimeBasedInflation) Type() protoreflect.MessageType {
	return _fastReflection_TimeBasedInflation_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_TimeBasedInflation) New() protoreflect.Message {
	return new(fastReflection_TimeBasedInflation)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_TimeBasedInflation) Interface() protoreflect.ProtoMessage {
	return (*TimeBasedInflation)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_TimeBasedInflation) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.StartBlockHeight != uint64(0) {
		value := protoreflect.ValueOfUint64(x.StartBlockHeight)
		if !f(fd_TimeBasedInflation_start_block_height, value) {
			return
		}
	}
	if x.EndBlockHeight != uint64(0) {
		value := protoreflect.ValueOfUint64(x.EndBlockHeight)
		if !f(fd_TimeBasedInflation_end_block_height, value) {
			return
		}
	}
	if x.Description != "" {
		value := protoreflect.ValueOfString(x.Description)
		if !f(fd_TimeBasedInflation_description, value) {
			return
		}
	}
	if x.Inflation != nil {
		value := protoreflect.ValueOfMessage(x.Inflation.ProtoReflect())
		if !f(fd_TimeBasedInflation_inflation, value) {
			return
		}
	}
	if x.Authority != "" {
		value := protoreflect.ValueOfString(x.Authority)
		if !f(fd_TimeBasedInflation_authority, value) {
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
func (x *fastReflection_TimeBasedInflation) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		return x.StartBlockHeight != uint64(0)
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		return x.EndBlockHeight != uint64(0)
	case "elys.tokenomics.TimeBasedInflation.description":
		return x.Description != ""
	case "elys.tokenomics.TimeBasedInflation.inflation":
		return x.Inflation != nil
	case "elys.tokenomics.TimeBasedInflation.authority":
		return x.Authority != ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_TimeBasedInflation) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		x.StartBlockHeight = uint64(0)
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		x.EndBlockHeight = uint64(0)
	case "elys.tokenomics.TimeBasedInflation.description":
		x.Description = ""
	case "elys.tokenomics.TimeBasedInflation.inflation":
		x.Inflation = nil
	case "elys.tokenomics.TimeBasedInflation.authority":
		x.Authority = ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_TimeBasedInflation) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		value := x.StartBlockHeight
		return protoreflect.ValueOfUint64(value)
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		value := x.EndBlockHeight
		return protoreflect.ValueOfUint64(value)
	case "elys.tokenomics.TimeBasedInflation.description":
		value := x.Description
		return protoreflect.ValueOfString(value)
	case "elys.tokenomics.TimeBasedInflation.inflation":
		value := x.Inflation
		return protoreflect.ValueOfMessage(value.ProtoReflect())
	case "elys.tokenomics.TimeBasedInflation.authority":
		value := x.Authority
		return protoreflect.ValueOfString(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", descriptor.FullName()))
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
func (x *fastReflection_TimeBasedInflation) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		x.StartBlockHeight = value.Uint()
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		x.EndBlockHeight = value.Uint()
	case "elys.tokenomics.TimeBasedInflation.description":
		x.Description = value.Interface().(string)
	case "elys.tokenomics.TimeBasedInflation.inflation":
		x.Inflation = value.Message().Interface().(*InflationEntry)
	case "elys.tokenomics.TimeBasedInflation.authority":
		x.Authority = value.Interface().(string)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", fd.FullName()))
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
func (x *fastReflection_TimeBasedInflation) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.tokenomics.TimeBasedInflation.inflation":
		if x.Inflation == nil {
			x.Inflation = new(InflationEntry)
		}
		return protoreflect.ValueOfMessage(x.Inflation.ProtoReflect())
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		panic(fmt.Errorf("field start_block_height of message elys.tokenomics.TimeBasedInflation is not mutable"))
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		panic(fmt.Errorf("field end_block_height of message elys.tokenomics.TimeBasedInflation is not mutable"))
	case "elys.tokenomics.TimeBasedInflation.description":
		panic(fmt.Errorf("field description of message elys.tokenomics.TimeBasedInflation is not mutable"))
	case "elys.tokenomics.TimeBasedInflation.authority":
		panic(fmt.Errorf("field authority of message elys.tokenomics.TimeBasedInflation is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_TimeBasedInflation) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "elys.tokenomics.TimeBasedInflation.start_block_height":
		return protoreflect.ValueOfUint64(uint64(0))
	case "elys.tokenomics.TimeBasedInflation.end_block_height":
		return protoreflect.ValueOfUint64(uint64(0))
	case "elys.tokenomics.TimeBasedInflation.description":
		return protoreflect.ValueOfString("")
	case "elys.tokenomics.TimeBasedInflation.inflation":
		m := new(InflationEntry)
		return protoreflect.ValueOfMessage(m.ProtoReflect())
	case "elys.tokenomics.TimeBasedInflation.authority":
		return protoreflect.ValueOfString("")
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: elys.tokenomics.TimeBasedInflation"))
		}
		panic(fmt.Errorf("message elys.tokenomics.TimeBasedInflation does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_TimeBasedInflation) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in elys.tokenomics.TimeBasedInflation", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_TimeBasedInflation) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_TimeBasedInflation) SetUnknown(fields protoreflect.RawFields) {
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
func (x *fastReflection_TimeBasedInflation) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_TimeBasedInflation) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*TimeBasedInflation)
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
		if x.StartBlockHeight != 0 {
			n += 1 + runtime.Sov(uint64(x.StartBlockHeight))
		}
		if x.EndBlockHeight != 0 {
			n += 1 + runtime.Sov(uint64(x.EndBlockHeight))
		}
		l = len(x.Description)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.Inflation != nil {
			l = options.Size(x.Inflation)
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Authority)
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
		x := input.Message.Interface().(*TimeBasedInflation)
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
		if len(x.Authority) > 0 {
			i -= len(x.Authority)
			copy(dAtA[i:], x.Authority)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Authority)))
			i--
			dAtA[i] = 0x2a
		}
		if x.Inflation != nil {
			encoded, err := options.Marshal(x.Inflation)
			if err != nil {
				return protoiface.MarshalOutput{
					NoUnkeyedLiterals: input.NoUnkeyedLiterals,
					Buf:               input.Buf,
				}, err
			}
			i -= len(encoded)
			copy(dAtA[i:], encoded)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
			i--
			dAtA[i] = 0x22
		}
		if len(x.Description) > 0 {
			i -= len(x.Description)
			copy(dAtA[i:], x.Description)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Description)))
			i--
			dAtA[i] = 0x1a
		}
		if x.EndBlockHeight != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.EndBlockHeight))
			i--
			dAtA[i] = 0x10
		}
		if x.StartBlockHeight != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.StartBlockHeight))
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
		x := input.Message.Interface().(*TimeBasedInflation)
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: TimeBasedInflation: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: TimeBasedInflation: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field StartBlockHeight", wireType)
				}
				x.StartBlockHeight = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.StartBlockHeight |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 2:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field EndBlockHeight", wireType)
				}
				x.EndBlockHeight = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.EndBlockHeight |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
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
				x.Description = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Inflation", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if x.Inflation == nil {
					x.Inflation = &InflationEntry{}
				}
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.Inflation); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 5:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
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
				x.Authority = string(dAtA[iNdEx:postIndex])
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
// source: elys/tokenomics/time_based_inflation.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TimeBasedInflation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartBlockHeight uint64          `protobuf:"varint,1,opt,name=start_block_height,json=startBlockHeight,proto3" json:"start_block_height,omitempty"`
	EndBlockHeight   uint64          `protobuf:"varint,2,opt,name=end_block_height,json=endBlockHeight,proto3" json:"end_block_height,omitempty"`
	Description      string          `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Inflation        *InflationEntry `protobuf:"bytes,4,opt,name=inflation,proto3" json:"inflation,omitempty"`
	Authority        string          `protobuf:"bytes,5,opt,name=authority,proto3" json:"authority,omitempty"`
}

func (x *TimeBasedInflation) Reset() {
	*x = TimeBasedInflation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_elys_tokenomics_time_based_inflation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeBasedInflation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeBasedInflation) ProtoMessage() {}

// Deprecated: Use TimeBasedInflation.ProtoReflect.Descriptor instead.
func (*TimeBasedInflation) Descriptor() ([]byte, []int) {
	return file_elys_tokenomics_time_based_inflation_proto_rawDescGZIP(), []int{0}
}

func (x *TimeBasedInflation) GetStartBlockHeight() uint64 {
	if x != nil {
		return x.StartBlockHeight
	}
	return 0
}

func (x *TimeBasedInflation) GetEndBlockHeight() uint64 {
	if x != nil {
		return x.EndBlockHeight
	}
	return 0
}

func (x *TimeBasedInflation) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *TimeBasedInflation) GetInflation() *InflationEntry {
	if x != nil {
		return x.Inflation
	}
	return nil
}

func (x *TimeBasedInflation) GetAuthority() string {
	if x != nil {
		return x.Authority
	}
	return ""
}

var File_elys_tokenomics_time_based_inflation_proto protoreflect.FileDescriptor

var file_elys_tokenomics_time_based_inflation_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63,
	0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x66,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x65, 0x6c,
	0x79, 0x73, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x1a, 0x25, 0x65,
	0x6c, 0x79, 0x73, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x2f, 0x69,
	0x6e, 0x66, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x01, 0x0a, 0x12, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x61, 0x73,
	0x65, 0x64, 0x49, 0x6e, 0x66, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x12, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x73, 0x74, 0x61, 0x72, 0x74, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x65, 0x6e, 0x64,
	0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0e, 0x65, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x09, 0x69, 0x6e, 0x66, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x65, 0x6c, 0x79, 0x73, 0x2e,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x2e, 0x49, 0x6e, 0x66, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x69, 0x6e, 0x66, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x42, 0xbd, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x79, 0x73, 0x2e,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x42, 0x17, 0x54, 0x69, 0x6d, 0x65,
	0x42, 0x61, 0x73, 0x65, 0x64, 0x49, 0x6e, 0x66, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x65,
	0x6c, 0x79, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6c, 0x79, 0x73, 0x2f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0xa2, 0x02, 0x03, 0x45, 0x54, 0x58, 0xaa, 0x02, 0x0f,
	0x45, 0x6c, 0x79, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0xca,
	0x02, 0x0f, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x63,
	0x73, 0xe2, 0x02, 0x1b, 0x45, 0x6c, 0x79, 0x73, 0x5c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d,
	0x69, 0x63, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x10, 0x45, 0x6c, 0x79, 0x73, 0x3a, 0x3a, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x6f, 0x6d, 0x69,
	0x63, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_elys_tokenomics_time_based_inflation_proto_rawDescOnce sync.Once
	file_elys_tokenomics_time_based_inflation_proto_rawDescData = file_elys_tokenomics_time_based_inflation_proto_rawDesc
)

func file_elys_tokenomics_time_based_inflation_proto_rawDescGZIP() []byte {
	file_elys_tokenomics_time_based_inflation_proto_rawDescOnce.Do(func() {
		file_elys_tokenomics_time_based_inflation_proto_rawDescData = protoimpl.X.CompressGZIP(file_elys_tokenomics_time_based_inflation_proto_rawDescData)
	})
	return file_elys_tokenomics_time_based_inflation_proto_rawDescData
}

var file_elys_tokenomics_time_based_inflation_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_elys_tokenomics_time_based_inflation_proto_goTypes = []interface{}{
	(*TimeBasedInflation)(nil), // 0: elys.tokenomics.TimeBasedInflation
	(*InflationEntry)(nil),     // 1: elys.tokenomics.InflationEntry
}
var file_elys_tokenomics_time_based_inflation_proto_depIdxs = []int32{
	1, // 0: elys.tokenomics.TimeBasedInflation.inflation:type_name -> elys.tokenomics.InflationEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_elys_tokenomics_time_based_inflation_proto_init() }
func file_elys_tokenomics_time_based_inflation_proto_init() {
	if File_elys_tokenomics_time_based_inflation_proto != nil {
		return
	}
	file_elys_tokenomics_inflation_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_elys_tokenomics_time_based_inflation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeBasedInflation); i {
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
			RawDescriptor: file_elys_tokenomics_time_based_inflation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_elys_tokenomics_time_based_inflation_proto_goTypes,
		DependencyIndexes: file_elys_tokenomics_time_based_inflation_proto_depIdxs,
		MessageInfos:      file_elys_tokenomics_time_based_inflation_proto_msgTypes,
	}.Build()
	File_elys_tokenomics_time_based_inflation_proto = out.File
	file_elys_tokenomics_time_based_inflation_proto_rawDesc = nil
	file_elys_tokenomics_time_based_inflation_proto_goTypes = nil
	file_elys_tokenomics_time_based_inflation_proto_depIdxs = nil
}
