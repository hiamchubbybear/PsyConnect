





package gprc_grpc_generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProfileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     string                 `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileRequest) Reset() {
	*x = ProfileRequest{}
	mi := &file_check_existed_profile_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileRequest) ProtoMessage() {}

func (x *ProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_check_existed_profile_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*ProfileRequest) Descriptor() ([]byte, []int) {
	return file_check_existed_profile_proto_rawDescGZIP(), []int{0}
}

func (x *ProfileRequest) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

type ProfileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Exists        bool                   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileResponse) Reset() {
	*x = ProfileResponse{}
	mi := &file_check_existed_profile_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileResponse) ProtoMessage() {}

func (x *ProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_check_existed_profile_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*ProfileResponse) Descriptor() ([]byte, []int) {
	return file_check_existed_profile_proto_rawDescGZIP(), []int{1}
}

func (x *ProfileResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

type ProfileRequestV1 struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     string                 `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileRequestV1) Reset() {
	*x = ProfileRequestV1{}
	mi := &file_check_existed_profile_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileRequestV1) ProtoMessage() {}

func (x *ProfileRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_check_existed_profile_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*ProfileRequestV1) Descriptor() ([]byte, []int) {
	return file_check_existed_profile_proto_rawDescGZIP(), []int{2}
}

func (x *ProfileRequestV1) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

type ProfileResponseV1 struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Exists        bool                   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	AvatarUri     string                 `protobuf:"bytes,2,opt,name=avatarUri,proto3" json:"avatarUri,omitempty"`
	Address       string                 `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Gender        string                 `protobuf:"bytes,4,opt,name=gender,proto3" json:"gender,omitempty"`
	Name          string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileResponseV1) Reset() {
	*x = ProfileResponseV1{}
	mi := &file_check_existed_profile_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileResponseV1) ProtoMessage() {}

func (x *ProfileResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_check_existed_profile_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*ProfileResponseV1) Descriptor() ([]byte, []int) {
	return file_check_existed_profile_proto_rawDescGZIP(), []int{3}
}

func (x *ProfileResponseV1) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

func (x *ProfileResponseV1) GetAvatarUri() string {
	if x != nil {
		return x.AvatarUri
	}
	return ""
}

func (x *ProfileResponseV1) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ProfileResponseV1) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *ProfileResponseV1) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_check_existed_profile_proto protoreflect.FileDescriptor

const file_check_existed_profile_proto_rawDesc = "" +
	"\n" +
	"\x1bcheck_existed_profile.proto\x12\aprofile\"/\n" +
	"\x0eProfileRequest\x12\x1d\n" +
	"\n" +
	"profile_id\x18\x01 \x01(\tR\tprofileId\")\n" +
	"\x0fProfileResponse\x12\x16\n" +
	"\x06exists\x18\x01 \x01(\bR\x06exists\"1\n" +
	"\x10ProfileRequestV1\x12\x1d\n" +
	"\n" +
	"profile_id\x18\x01 \x01(\tR\tprofileId\"\x8f\x01\n" +
	"\x11ProfileResponseV1\x12\x16\n" +
	"\x06exists\x18\x01 \x01(\bR\x06exists\x12\x1c\n" +
	"\tavatarUri\x18\x02 \x01(\tR\tavatarUri\x12\x18\n" +
	"\aaddress\x18\x03 \x01(\tR\aaddress\x12\x16\n" +
	"\x06gender\x18\x04 \x01(\tR\x06gender\x12\x12\n" +
	"\x04name\x18\x05 \x01(\tR\x04name2\xad\x01\n" +
	"\x13CheckProfileService\x12G\n" +
	"\x12CheckProfileExists\x12\x17.profile.ProfileRequest\x1a\x18.profile.ProfileResponse\x12M\n" +
	"\x14CheckProfileExistsV1\x12\x19.profile.ProfileRequestV1\x1a\x1a.profile.ProfileResponseV1B\x15Z\x13gprc.grpc_generatedb\x06proto3"

var (
	file_check_existed_profile_proto_rawDescOnce sync.Once
	file_check_existed_profile_proto_rawDescData []byte
)

func file_check_existed_profile_proto_rawDescGZIP() []byte {
	file_check_existed_profile_proto_rawDescOnce.Do(func() {
		file_check_existed_profile_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_check_existed_profile_proto_rawDesc), len(file_check_existed_profile_proto_rawDesc)))
	})
	return file_check_existed_profile_proto_rawDescData
}

var file_check_existed_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_check_existed_profile_proto_goTypes = []any{
	(*ProfileRequest)(nil),    
	(*ProfileResponse)(nil),   
	(*ProfileRequestV1)(nil),  
	(*ProfileResponseV1)(nil), 
}
var file_check_existed_profile_proto_depIdxs = []int32{
	0, 
	2, 
	1, 
	3, 
	2, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_check_existed_profile_proto_init() }
func file_check_existed_profile_proto_init() {
	if File_check_existed_profile_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_check_existed_profile_proto_rawDesc), len(file_check_existed_profile_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_check_existed_profile_proto_goTypes,
		DependencyIndexes: file_check_existed_profile_proto_depIdxs,
		MessageInfos:      file_check_existed_profile_proto_msgTypes,
	}.Build()
	File_check_existed_profile_proto = out.File
	file_check_existed_profile_proto_goTypes = nil
	file_check_existed_profile_proto_depIdxs = nil
}
