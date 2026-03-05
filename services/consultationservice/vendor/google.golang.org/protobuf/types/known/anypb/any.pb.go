

















































































































package anypb

import (
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	strings "strings"
	sync "sync"
	unsafe "unsafe"
)






















































































type Any struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	TypeUrl string `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	
	Value         []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}


func New(src proto.Message) (*Any, error) {
	dst := new(Any)
	if err := dst.MarshalFrom(src); err != nil {
		return nil, err
	}
	return dst, nil
}





func MarshalFrom(dst *Any, src proto.Message, opts proto.MarshalOptions) error {
	const urlPrefix = "type.googleapis.com/"
	if src == nil {
		return protoimpl.X.NewError("invalid nil source message")
	}
	b, err := opts.Marshal(src)
	if err != nil {
		return err
	}
	dst.TypeUrl = urlPrefix + string(src.ProtoReflect().Descriptor().FullName())
	dst.Value = b
	return nil
}






func UnmarshalTo(src *Any, dst proto.Message, opts proto.UnmarshalOptions) error {
	if src == nil {
		return protoimpl.X.NewError("invalid nil source message")
	}
	if !src.MessageIs(dst) {
		got := dst.ProtoReflect().Descriptor().FullName()
		want := src.MessageName()
		return protoimpl.X.NewError("mismatched message type: got %q, want %q", got, want)
	}
	return opts.Unmarshal(src.GetValue(), dst)
}








func UnmarshalNew(src *Any, opts proto.UnmarshalOptions) (dst proto.Message, err error) {
	if src.GetTypeUrl() == "" {
		return nil, protoimpl.X.NewError("invalid empty type URL")
	}
	if opts.Resolver == nil {
		opts.Resolver = protoregistry.GlobalTypes
	}
	r, ok := opts.Resolver.(protoregistry.MessageTypeResolver)
	if !ok {
		return nil, protoregistry.NotFound
	}
	mt, err := r.FindMessageByURL(src.GetTypeUrl())
	if err != nil {
		if err == protoregistry.NotFound {
			return nil, err
		}
		return nil, protoimpl.X.NewError("could not resolve %q: %v", src.GetTypeUrl(), err)
	}
	dst = mt.New().Interface()
	return dst, opts.Unmarshal(src.GetValue(), dst)
}


func (x *Any) MessageIs(m proto.Message) bool {
	if m == nil {
		return false
	}
	url := x.GetTypeUrl()
	name := string(m.ProtoReflect().Descriptor().FullName())
	if !strings.HasSuffix(url, name) {
		return false
	}
	return len(url) == len(name) || url[len(url)-len(name)-1] == '/'
}



func (x *Any) MessageName() protoreflect.FullName {
	url := x.GetTypeUrl()
	name := protoreflect.FullName(url)
	if i := strings.LastIndexByte(url, '/'); i >= 0 {
		name = name[i+len("/"):]
	}
	if !name.IsValid() {
		return ""
	}
	return name
}


func (x *Any) MarshalFrom(m proto.Message) error {
	return MarshalFrom(x, m, proto.MarshalOptions{})
}




func (x *Any) UnmarshalTo(m proto.Message) error {
	return UnmarshalTo(x, m, proto.UnmarshalOptions{})
}




func (x *Any) UnmarshalNew() (proto.Message, error) {
	return UnmarshalNew(x, proto.UnmarshalOptions{})
}

func (x *Any) Reset() {
	*x = Any{}
	mi := &file_google_protobuf_any_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Any) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Any) ProtoMessage() {}

func (x *Any) ProtoReflect() protoreflect.Message {
	mi := &file_google_protobuf_any_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*Any) Descriptor() ([]byte, []int) {
	return file_google_protobuf_any_proto_rawDescGZIP(), []int{0}
}

func (x *Any) GetTypeUrl() string {
	if x != nil {
		return x.TypeUrl
	}
	return ""
}

func (x *Any) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_google_protobuf_any_proto protoreflect.FileDescriptor

const file_google_protobuf_any_proto_rawDesc = "" +
	"\n" +
	"\x19google/protobuf/any.proto\x12\x0fgoogle.protobuf\"6\n" +
	"\x03Any\x12\x19\n" +
	"\btype_url\x18\x01 \x01(\tR\atypeUrl\x12\x14\n" +
	"\x05value\x18\x02 \x01(\fR\x05valueBv\n" +
	"\x13com.google.protobufB\bAnyProtoP\x01Z,google.golang.org/protobuf/types/known/anypb\xa2\x02\x03GPB\xaa\x02\x1eGoogle.Protobuf.WellKnownTypesb\x06proto3"

var (
	file_google_protobuf_any_proto_rawDescOnce sync.Once
	file_google_protobuf_any_proto_rawDescData []byte
)

func file_google_protobuf_any_proto_rawDescGZIP() []byte {
	file_google_protobuf_any_proto_rawDescOnce.Do(func() {
		file_google_protobuf_any_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_google_protobuf_any_proto_rawDesc), len(file_google_protobuf_any_proto_rawDesc)))
	})
	return file_google_protobuf_any_proto_rawDescData
}

var file_google_protobuf_any_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_google_protobuf_any_proto_goTypes = []any{
	(*Any)(nil), 
}
var file_google_protobuf_any_proto_depIdxs = []int32{
	0, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_google_protobuf_any_proto_init() }
func file_google_protobuf_any_proto_init() {
	if File_google_protobuf_any_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_google_protobuf_any_proto_rawDesc), len(file_google_protobuf_any_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_protobuf_any_proto_goTypes,
		DependencyIndexes: file_google_protobuf_any_proto_depIdxs,
		MessageInfos:      file_google_protobuf_any_proto_msgTypes,
	}.Build()
	File_google_protobuf_any_proto = out.File
	file_google_protobuf_any_proto_goTypes = nil
	file_google_protobuf_any_proto_depIdxs = nil
}
