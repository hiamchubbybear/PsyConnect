





package genid

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

const File_google_protobuf_go_features_proto = "google/protobuf/go_features.proto"


const (
	GoFeatures_message_name     protoreflect.Name     = "GoFeatures"
	GoFeatures_message_fullname protoreflect.FullName = "pb.GoFeatures"
)


const (
	GoFeatures_LegacyUnmarshalJsonEnum_field_name protoreflect.Name = "legacy_unmarshal_json_enum"
	GoFeatures_ApiLevel_field_name                protoreflect.Name = "api_level"
	GoFeatures_StripEnumPrefix_field_name         protoreflect.Name = "strip_enum_prefix"

	GoFeatures_LegacyUnmarshalJsonEnum_field_fullname protoreflect.FullName = "pb.GoFeatures.legacy_unmarshal_json_enum"
	GoFeatures_ApiLevel_field_fullname                protoreflect.FullName = "pb.GoFeatures.api_level"
	GoFeatures_StripEnumPrefix_field_fullname         protoreflect.FullName = "pb.GoFeatures.strip_enum_prefix"
)


const (
	GoFeatures_LegacyUnmarshalJsonEnum_field_number protoreflect.FieldNumber = 1
	GoFeatures_ApiLevel_field_number                protoreflect.FieldNumber = 2
	GoFeatures_StripEnumPrefix_field_number         protoreflect.FieldNumber = 3
)


const (
	GoFeatures_APILevel_enum_fullname = "pb.GoFeatures.APILevel"
	GoFeatures_APILevel_enum_name     = "APILevel"
)


const (
	GoFeatures_API_LEVEL_UNSPECIFIED_enum_value = 0
	GoFeatures_API_OPEN_enum_value              = 1
	GoFeatures_API_HYBRID_enum_value            = 2
	GoFeatures_API_OPAQUE_enum_value            = 3
)


const (
	GoFeatures_StripEnumPrefix_enum_fullname = "pb.GoFeatures.StripEnumPrefix"
	GoFeatures_StripEnumPrefix_enum_name     = "StripEnumPrefix"
)


const (
	GoFeatures_STRIP_ENUM_PREFIX_UNSPECIFIED_enum_value   = 0
	GoFeatures_STRIP_ENUM_PREFIX_KEEP_enum_value          = 1
	GoFeatures_STRIP_ENUM_PREFIX_GENERATE_BOTH_enum_value = 2
	GoFeatures_STRIP_ENUM_PREFIX_STRIP_enum_value         = 3
)


const (
	FeatureSet_Go_ext_number protoreflect.FieldNumber = 1002
)
