




package protoadapt

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)


type MessageV1 = protoiface.MessageV1



type MessageV2 = proto.Message



func MessageV1Of(m MessageV2) MessageV1 {
	return protoimpl.X.ProtoMessageV1Of(m)
}



func MessageV2Of(m MessageV1) MessageV2 {
	return protoimpl.X.ProtoMessageV2Of(m)
}
