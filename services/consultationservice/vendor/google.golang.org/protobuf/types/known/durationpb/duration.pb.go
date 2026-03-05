








































































package durationpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	math "math"
	reflect "reflect"
	sync "sync"
	time "time"
	unsafe "unsafe"
)



























































type Duration struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	
	
	
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	
	
	
	
	
	
	Nanos         int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}


func New(d time.Duration) *Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &Duration{Seconds: int64(secs), Nanos: int32(nanos)}
}



func (x *Duration) AsDuration() time.Duration {
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	d := time.Duration(secs) * time.Second
	overflow := d/time.Second != time.Duration(secs)
	d += time.Duration(nanos) * time.Nanosecond
	overflow = overflow || (secs < 0 && nanos < 0 && d > 0)
	overflow = overflow || (secs > 0 && nanos > 0 && d < 0)
	if overflow {
		switch {
		case secs < 0:
			return time.Duration(math.MinInt64)
		case secs > 0:
			return time.Duration(math.MaxInt64)
		}
	}
	return d
}



func (x *Duration) IsValid() bool {
	return x.check() == 0
}





func (x *Duration) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return protoimpl.X.NewError("invalid nil Duration")
	case invalidUnderflow:
		return protoimpl.X.NewError("duration (%v) exceeds -10000 years", x)
	case invalidOverflow:
		return protoimpl.X.NewError("duration (%v) exceeds +10000 years", x)
	case invalidNanosRange:
		return protoimpl.X.NewError("duration (%v) has out-of-range nanos", x)
	case invalidNanosSign:
		return protoimpl.X.NewError("duration (%v) has seconds and nanos with different signs", x)
	default:
		return nil
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanosRange
	invalidNanosSign
)

func (x *Duration) check() uint {
	const absDuration = 315576000000 
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < -absDuration:
		return invalidUnderflow
	case secs > +absDuration:
		return invalidOverflow
	case nanos <= -1e9 || nanos >= +1e9:
		return invalidNanosRange
	case (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0):
		return invalidNanosSign
	default:
		return 0
	}
}

func (x *Duration) Reset() {
	*x = Duration{}
	mi := &file_google_protobuf_duration_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Duration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Duration) ProtoMessage() {}

func (x *Duration) ProtoReflect() protoreflect.Message {
	mi := &file_google_protobuf_duration_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*Duration) Descriptor() ([]byte, []int) {
	return file_google_protobuf_duration_proto_rawDescGZIP(), []int{0}
}

func (x *Duration) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *Duration) GetNanos() int32 {
	if x != nil {
		return x.Nanos
	}
	return 0
}

var File_google_protobuf_duration_proto protoreflect.FileDescriptor

const file_google_protobuf_duration_proto_rawDesc = "" +
	"\n" +
	"\x1egoogle/protobuf/duration.proto\x12\x0fgoogle.protobuf\":\n" +
	"\bDuration\x12\x18\n" +
	"\aseconds\x18\x01 \x01(\x03R\aseconds\x12\x14\n" +
	"\x05nanos\x18\x02 \x01(\x05R\x05nanosB\x83\x01\n" +
	"\x13com.google.protobufB\rDurationProtoP\x01Z1google.golang.org/protobuf/types/known/durationpb\xf8\x01\x01\xa2\x02\x03GPB\xaa\x02\x1eGoogle.Protobuf.WellKnownTypesb\x06proto3"

var (
	file_google_protobuf_duration_proto_rawDescOnce sync.Once
	file_google_protobuf_duration_proto_rawDescData []byte
)

func file_google_protobuf_duration_proto_rawDescGZIP() []byte {
	file_google_protobuf_duration_proto_rawDescOnce.Do(func() {
		file_google_protobuf_duration_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_google_protobuf_duration_proto_rawDesc), len(file_google_protobuf_duration_proto_rawDesc)))
	})
	return file_google_protobuf_duration_proto_rawDescData
}

var file_google_protobuf_duration_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_google_protobuf_duration_proto_goTypes = []any{
	(*Duration)(nil), 
}
var file_google_protobuf_duration_proto_depIdxs = []int32{
	0, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_google_protobuf_duration_proto_init() }
func file_google_protobuf_duration_proto_init() {
	if File_google_protobuf_duration_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_google_protobuf_duration_proto_rawDesc), len(file_google_protobuf_duration_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_protobuf_duration_proto_goTypes,
		DependencyIndexes: file_google_protobuf_duration_proto_depIdxs,
		MessageInfos:      file_google_protobuf_duration_proto_msgTypes,
	}.Build()
	File_google_protobuf_duration_proto = out.File
	file_google_protobuf_duration_proto_goTypes = nil
	file_google_protobuf_duration_proto_depIdxs = nil
}
