



package impl

import (
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"math"
	"reflect"

	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)





func (Export) LegacyEnumName(ed protoreflect.EnumDescriptor) string {
	return legacyEnumName(ed)
}



func (Export) LegacyMessageTypeOf(m protoiface.MessageV1, name protoreflect.FullName) protoreflect.MessageType {
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect().Type()
	}
	return legacyLoadMessageType(reflect.TypeOf(m), name)
}




func (Export) UnmarshalJSONEnum(ed protoreflect.EnumDescriptor, b []byte) (protoreflect.EnumNumber, error) {
	if b[0] == '"' {
		var name protoreflect.Name
		if err := json.Unmarshal(b, &name); err != nil {
			return 0, errors.New("invalid input for enum %v: %s", ed.FullName(), b)
		}
		ev := ed.Values().ByName(name)
		if ev == nil {
			return 0, errors.New("invalid value for enum %v: %s", ed.FullName(), name)
		}
		return ev.Number(), nil
	} else {
		var num protoreflect.EnumNumber
		if err := json.Unmarshal(b, &num); err != nil {
			return 0, errors.New("invalid input for enum %v: %s", ed.FullName(), b)
		}
		return num, nil
	}
}



func (Export) CompressGZIP(in []byte) (out []byte) {
	
	var gzipHeader = [10]byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}

	
	var blockHeader [5]byte
	const maxBlockSize = math.MaxUint16
	numBlocks := 1 + len(in)/maxBlockSize

	
	var gzipFooter [8]byte
	binary.LittleEndian.PutUint32(gzipFooter[0:4], crc32.ChecksumIEEE(in))
	binary.LittleEndian.PutUint32(gzipFooter[4:8], uint32(len(in)))

	
	out = make([]byte, 0, len(gzipHeader)+len(blockHeader)*numBlocks+len(in)+len(gzipFooter))
	out = append(out, gzipHeader[:]...)
	for blockHeader[0] == 0 {
		blockSize := maxBlockSize
		if blockSize > len(in) {
			blockHeader[0] = 0x01 
			blockSize = len(in)
		}
		binary.LittleEndian.PutUint16(blockHeader[1:3], uint16(blockSize))
		binary.LittleEndian.PutUint16(blockHeader[3:5], ^uint16(blockSize))
		out = append(out, blockHeader[:]...)
		out = append(out, in[:blockSize]...)
		in = in[blockSize:]
	}
	out = append(out, gzipFooter[:]...)
	return out
}
