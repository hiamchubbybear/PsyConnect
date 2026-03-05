







































































package timestamppb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	time "time"
	unsafe "unsafe"
)


























































































type Timestamp struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	
	
	
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	
	
	
	
	Nanos         int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}


func Now() *Timestamp {
	return New(time.Now())
}


func New(t time.Time) *Timestamp {
	return &Timestamp{Seconds: int64(t.Unix()), Nanos: int32(t.Nanosecond())}
}


func (x *Timestamp) AsTime() time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).UTC()
}



func (x *Timestamp) IsValid() bool {
	return x.check() == 0
}





func (x *Timestamp) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return protoimpl.X.NewError("invalid nil Timestamp")
	case invalidUnderflow:
		return protoimpl.X.NewError("timestamp (%v) before 0001-01-01", x)
	case invalidOverflow:
		return protoimpl.X.NewError("timestamp (%v) after 9999-12-31", x)
	case invalidNanos:
		return protoimpl.X.NewError("timestamp (%v) has out-of-range nanos", x)
	default:
		return nil
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
)

func (x *Timestamp) check() uint {
	const minTimestamp = -62135596800  
	const maxTimestamp = +253402300799 
	secs := x.GetSeconds()
	nanos := x.GetNanos()
	switch {
	case x == nil:
		return invalidNil
	case secs < minTimestamp:
		return invalidUnderflow
	case secs > maxTimestamp:
		return invalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return invalidNanos
	default:
		return 0
	}
}

func (x *Timestamp) Reset() {
	*x = Timestamp{}
	mi := &file_google_protobuf_timestamp_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Timestamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timestamp) ProtoMessage() {}

func (x *Timestamp) ProtoReflect() protoreflect.Message {
	mi := &file_google_protobuf_timestamp_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*Timestamp) Descriptor() ([]byte, []int) {
	return file_google_protobuf_timestamp_proto_rawDescGZIP(), []int{0}
}

func (x *Timestamp) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *Timestamp) GetNanos() int32 {
	if x != nil {
		return x.Nanos
	}
	return 0
}

var File_google_protobuf_timestamp_proto protoreflect.FileDescriptor

const file_google_protobuf_timestamp_proto_rawDesc = "" +
	"\n" +
	"\x1fgoogle/protobuf/timestamp.proto\x12\x0fgoogle.protobuf\";\n" +
	"\tTimestamp\x12\x18\n" +
	"\aseconds\x18\x01 \x01(\x03R\aseconds\x12\x14\n" +
	"\x05nanos\x18\x02 \x01(\x05R\x05nanosB\x85\x01\n" +
	"\x13com.google.protobufB\x0eTimestampProtoP\x01Z2google.golang.org/protobuf/types/known/timestamppb\xf8\x01\x01\xa2\x02\x03GPB\xaa\x02\x1eGoogle.Protobuf.WellKnownTypesb\x06proto3"

var (
	file_google_protobuf_timestamp_proto_rawDescOnce sync.Once
	file_google_protobuf_timestamp_proto_rawDescData []byte
)

func file_google_protobuf_timestamp_proto_rawDescGZIP() []byte {
	file_google_protobuf_timestamp_proto_rawDescOnce.Do(func() {
		file_google_protobuf_timestamp_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_google_protobuf_timestamp_proto_rawDesc), len(file_google_protobuf_timestamp_proto_rawDesc)))
	})
	return file_google_protobuf_timestamp_proto_rawDescData
}

var file_google_protobuf_timestamp_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_google_protobuf_timestamp_proto_goTypes = []any{
	(*Timestamp)(nil), 
}
var file_google_protobuf_timestamp_proto_depIdxs = []int32{
	0, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_google_protobuf_timestamp_proto_init() }
func file_google_protobuf_timestamp_proto_init() {
	if File_google_protobuf_timestamp_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_google_protobuf_timestamp_proto_rawDesc), len(file_google_protobuf_timestamp_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_protobuf_timestamp_proto_goTypes,
		DependencyIndexes: file_google_protobuf_timestamp_proto_depIdxs,
		MessageInfos:      file_google_protobuf_timestamp_proto_msgTypes,
	}.Build()
	File_google_protobuf_timestamp_proto = out.File
	file_google_protobuf_timestamp_proto_goTypes = nil
	file_google_protobuf_timestamp_proto_depIdxs = nil
}
