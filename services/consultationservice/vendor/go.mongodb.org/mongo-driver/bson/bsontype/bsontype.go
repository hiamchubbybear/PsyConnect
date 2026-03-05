







package bsontype 




const (
	Double           Type = 0x01
	String           Type = 0x02
	EmbeddedDocument Type = 0x03
	Array            Type = 0x04
	Binary           Type = 0x05
	Undefined        Type = 0x06
	ObjectID         Type = 0x07
	Boolean          Type = 0x08
	DateTime         Type = 0x09
	Null             Type = 0x0A
	Regex            Type = 0x0B
	DBPointer        Type = 0x0C
	JavaScript       Type = 0x0D
	Symbol           Type = 0x0E
	CodeWithScope    Type = 0x0F
	Int32            Type = 0x10
	Timestamp        Type = 0x11
	Int64            Type = 0x12
	Decimal128       Type = 0x13
	MinKey           Type = 0xFF
	MaxKey           Type = 0x7F
)




const (
	BinaryGeneric     byte = 0x00
	BinaryFunction    byte = 0x01
	BinaryBinaryOld   byte = 0x02
	BinaryUUIDOld     byte = 0x03
	BinaryUUID        byte = 0x04
	BinaryMD5         byte = 0x05
	BinaryEncrypted   byte = 0x06
	BinaryColumn      byte = 0x07
	BinarySensitive   byte = 0x08
	BinaryUserDefined byte = 0x80
)


type Type byte


func (bt Type) String() string {
	switch bt {
	case '\x01':
		return "double"
	case '\x02':
		return "string"
	case '\x03':
		return "embedded document"
	case '\x04':
		return "array"
	case '\x05':
		return "binary"
	case '\x06':
		return "undefined"
	case '\x07':
		return "objectID"
	case '\x08':
		return "boolean"
	case '\x09':
		return "UTC datetime"
	case '\x0A':
		return "null"
	case '\x0B':
		return "regex"
	case '\x0C':
		return "dbPointer"
	case '\x0D':
		return "javascript"
	case '\x0E':
		return "symbol"
	case '\x0F':
		return "code with scope"
	case '\x10':
		return "32-bit integer"
	case '\x11':
		return "timestamp"
	case '\x12':
		return "64-bit integer"
	case '\x13':
		return "128-bit decimal"
	case '\xFF':
		return "min key"
	case '\x7F':
		return "max key"
	default:
		return "invalid"
	}
}


func (bt Type) IsValid() bool {
	switch bt {
	case Double, String, EmbeddedDocument, Array, Binary, Undefined, ObjectID, Boolean, DateTime, Null, Regex,
		DBPointer, JavaScript, Symbol, CodeWithScope, Int32, Timestamp, Int64, Decimal128, MinKey, MaxKey:
		return true
	default:
		return false
	}
}
