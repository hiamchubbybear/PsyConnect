












package wiremessage

import (
	"bytes"
	"encoding/binary"
	"strings"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type WireMessage []byte

var globalRequestID int32


func NextRequestID() int32 { return atomic.AddInt32(&globalRequestID, 1) }


type OpCode int32




const (
	OpReply  OpCode = 1
	_        OpCode = 1001
	OpUpdate OpCode = 2001
	OpInsert OpCode = 2002
	_        OpCode = 2003
	
	OpQuery        OpCode = 2004
	OpGetMore      OpCode = 2005
	OpDelete       OpCode = 2006
	OpKillCursors  OpCode = 2007
	OpCommand      OpCode = 2010
	OpCommandReply OpCode = 2011
	OpCompressed   OpCode = 2012
	OpMsg          OpCode = 2013
)


func (oc OpCode) String() string {
	switch oc {
	case OpReply:
		return "OP_REPLY"
	case OpUpdate:
		return "OP_UPDATE"
	case OpInsert:
		return "OP_INSERT"
	case OpQuery:
		return "OP_QUERY"
	case OpGetMore:
		return "OP_GET_MORE"
	case OpDelete:
		return "OP_DELETE"
	case OpKillCursors:
		return "OP_KILL_CURSORS"
	case OpCommand:
		return "OP_COMMAND"
	case OpCommandReply:
		return "OP_COMMANDREPLY"
	case OpCompressed:
		return "OP_COMPRESSED"
	case OpMsg:
		return "OP_MSG"
	default:
		return "<invalid opcode>"
	}
}


type QueryFlag int32


const (
	_ QueryFlag = 1 << iota
	TailableCursor
	SecondaryOK
	OplogReplay
	NoCursorTimeout
	AwaitData
	Exhaust
	Partial
)


func (qf QueryFlag) String() string {
	strs := make([]string, 0)
	if qf&TailableCursor == TailableCursor {
		strs = append(strs, "TailableCursor")
	}
	if qf&SecondaryOK == SecondaryOK {
		strs = append(strs, "SecondaryOK")
	}
	if qf&OplogReplay == OplogReplay {
		strs = append(strs, "OplogReplay")
	}
	if qf&NoCursorTimeout == NoCursorTimeout {
		strs = append(strs, "NoCursorTimeout")
	}
	if qf&AwaitData == AwaitData {
		strs = append(strs, "AwaitData")
	}
	if qf&Exhaust == Exhaust {
		strs = append(strs, "Exhaust")
	}
	if qf&Partial == Partial {
		strs = append(strs, "Partial")
	}
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
}


type MsgFlag uint32


const (
	ChecksumPresent MsgFlag = 1 << iota
	MoreToCome

	ExhaustAllowed MsgFlag = 1 << 16
)


type ReplyFlag int32


const (
	CursorNotFound ReplyFlag = 1 << iota
	QueryFailure
	ShardConfigStale
	AwaitCapable
)


func (rf ReplyFlag) String() string {
	strs := make([]string, 0)
	if rf&CursorNotFound == CursorNotFound {
		strs = append(strs, "CursorNotFound")
	}
	if rf&QueryFailure == QueryFailure {
		strs = append(strs, "QueryFailure")
	}
	if rf&ShardConfigStale == ShardConfigStale {
		strs = append(strs, "ShardConfigStale")
	}
	if rf&AwaitCapable == AwaitCapable {
		strs = append(strs, "AwaitCapable")
	}
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
}


type SectionType uint8


const (
	SingleDocument SectionType = iota
	DocumentSequence
)


type CompressorID uint8


const (
	CompressorNoOp CompressorID = iota
	CompressorSnappy
	CompressorZLib
	CompressorZstd
)


func (id CompressorID) String() string {
	switch id {
	case CompressorNoOp:
		return "CompressorNoOp"
	case CompressorSnappy:
		return "CompressorSnappy"
	case CompressorZLib:
		return "CompressorZLib"
	case CompressorZstd:
		return "CompressorZstd"
	default:
		return "CompressorInvalid"
	}
}

const (
	
	DefaultZlibLevel = 6
	
	
	DefaultZstdLevel = 6
)



func AppendHeaderStart(dst []byte, reqid, respto int32, opcode OpCode) (index int32, b []byte) {
	index, dst = bsoncore.ReserveLength(dst)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return index, dst
}


func AppendHeader(dst []byte, length, reqid, respto int32, opcode OpCode) []byte {
	dst = appendi32(dst, length)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return dst
}


func ReadHeader(src []byte) (length, requestID, responseTo int32, opcode OpCode, rem []byte, ok bool) {
	if len(src) < 16 {
		return 0, 0, 0, 0, src, false
	}

	length = readi32unsafe(src)
	requestID = readi32unsafe(src[4:])
	responseTo = readi32unsafe(src[8:])
	opcode = OpCode(readi32unsafe(src[12:]))
	return length, requestID, responseTo, opcode, src[16:], true
}


func AppendQueryFlags(dst []byte, flags QueryFlag) []byte {
	return appendi32(dst, int32(flags))
}


func AppendMsgFlags(dst []byte, flags MsgFlag) []byte {
	return appendi32(dst, int32(flags))
}


func AppendReplyFlags(dst []byte, flags ReplyFlag) []byte {
	return appendi32(dst, int32(flags))
}


func AppendMsgSectionType(dst []byte, stype SectionType) []byte {
	return append(dst, byte(stype))
}


func AppendQueryFullCollectionName(dst []byte, ns string) []byte {
	return appendCString(dst, ns)
}


func AppendQueryNumberToSkip(dst []byte, skip int32) []byte {
	return appendi32(dst, skip)
}


func AppendQueryNumberToReturn(dst []byte, nor int32) []byte {
	return appendi32(dst, nor)
}


func AppendReplyCursorID(dst []byte, id int64) []byte {
	return appendi64(dst, id)
}


func AppendReplyStartingFrom(dst []byte, sf int32) []byte {
	return appendi32(dst, sf)
}


func AppendReplyNumberReturned(dst []byte, nr int32) []byte {
	return appendi32(dst, nr)
}


func AppendCompressedOriginalOpCode(dst []byte, opcode OpCode) []byte {
	return appendi32(dst, int32(opcode))
}



func AppendCompressedUncompressedSize(dst []byte, size int32) []byte { return appendi32(dst, size) }


func AppendCompressedCompressorID(dst []byte, id CompressorID) []byte {
	return append(dst, byte(id))
}


func AppendCompressedCompressedMessage(dst []byte, msg []byte) []byte { return append(dst, msg...) }


func AppendGetMoreZero(dst []byte) []byte {
	return appendi32(dst, 0)
}


func AppendGetMoreFullCollectionName(dst []byte, ns string) []byte {
	return appendCString(dst, ns)
}


func AppendGetMoreNumberToReturn(dst []byte, numToReturn int32) []byte {
	return appendi32(dst, numToReturn)
}


func AppendGetMoreCursorID(dst []byte, cursorID int64) []byte {
	return appendi64(dst, cursorID)
}


func AppendKillCursorsZero(dst []byte) []byte {
	return appendi32(dst, 0)
}


func AppendKillCursorsNumberIDs(dst []byte, numIDs int32) []byte {
	return appendi32(dst, numIDs)
}


func AppendKillCursorsCursorIDs(dst []byte, cursors []int64) []byte {
	for _, cursor := range cursors {
		dst = appendi64(dst, cursor)
	}
	return dst
}


func ReadMsgFlags(src []byte) (flags MsgFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return MsgFlag(i32), rem, ok
}


func IsMsgMoreToCome(wm []byte) bool {
	return len(wm) >= 20 &&
		OpCode(readi32unsafe(wm[12:16])) == OpMsg &&
		MsgFlag(readi32unsafe(wm[16:20]))&MoreToCome == MoreToCome
}


func ReadMsgSectionType(src []byte) (stype SectionType, rem []byte, ok bool) {
	if len(src) < 1 {
		return 0, src, false
	}
	return SectionType(src[0]), src[1:], true
}


func ReadMsgSectionSingleDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}



func ReadMsgSectionDocumentSequence(src []byte) (identifier string, docs []bsoncore.Document, rem []byte, ok bool) {
	identifier, rem, ret, ok := ReadMsgSectionRawDocumentSequence(src)
	if !ok {
		return "", nil, src, false
	}

	docs = make([]bsoncore.Document, 0)
	var doc bsoncore.Document
	for {
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok {
			break
		}
		docs = append(docs, doc)
	}
	if len(rem) > 0 {
		return "", nil, src, false
	}

	return identifier, docs, ret, true
}



func ReadMsgSectionRawDocumentSequence(src []byte) (identifier string, data []byte, rem []byte, ok bool) {
	length, rem, ok := readi32(src)
	if !ok || int(length) > len(src) || length-4 < 0 {
		return "", nil, src, false
	}

	
	
	rem, rest := rem[:length-4], rem[length-4:]

	identifier, rem, ok = readcstring(rem)
	if !ok {
		return "", nil, src, false
	}

	return identifier, rem, rest, true
}


func ReadMsgChecksum(src []byte) (checksum uint32, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return uint32(i32), rem, ok
}





func ReadQueryFlags(src []byte) (flags QueryFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return QueryFlag(i32), rem, ok
}





func ReadQueryFullCollectionName(src []byte) (collname string, rem []byte, ok bool) {
	return readcstring(src)
}





func ReadQueryNumberToSkip(src []byte) (nts int32, rem []byte, ok bool) {
	return readi32(src)
}





func ReadQueryNumberToReturn(src []byte) (ntr int32, rem []byte, ok bool) {
	return readi32(src)
}





func ReadQueryQuery(src []byte) (query bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}





func ReadQueryReturnFieldsSelector(src []byte) (rfs bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}


func ReadReplyFlags(src []byte) (flags ReplyFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return ReplyFlag(i32), rem, ok
}


func ReadReplyCursorID(src []byte) (cursorID int64, rem []byte, ok bool) {
	return readi64(src)
}


func ReadReplyStartingFrom(src []byte) (startingFrom int32, rem []byte, ok bool) {
	return readi32(src)
}


func ReadReplyNumberReturned(src []byte) (numberReturned int32, rem []byte, ok bool) {
	return readi32(src)
}


func ReadReplyDocuments(src []byte) (docs []bsoncore.Document, rem []byte, ok bool) {
	rem = src
	for {
		var doc bsoncore.Document
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok {
			break
		}

		docs = append(docs, doc)
	}

	return docs, rem, true
}


func ReadReplyDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}


func ReadCompressedOriginalOpCode(src []byte) (opcode OpCode, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return OpCode(i32), rem, ok
}



func ReadCompressedUncompressedSize(src []byte) (size int32, rem []byte, ok bool) {
	return readi32(src)
}


func ReadCompressedCompressorID(src []byte) (id CompressorID, rem []byte, ok bool) {
	if len(src) < 1 {
		return 0, src, false
	}
	return CompressorID(src[0]), src[1:], true
}





func ReadCompressedCompressedMessage(src []byte, length int32) (msg []byte, rem []byte, ok bool) {
	if len(src) < int(length) || length < 0 {
		return nil, src, false
	}
	return src[:length], src[length:], true
}


func ReadKillCursorsZero(src []byte) (zero int32, rem []byte, ok bool) {
	return readi32(src)
}


func ReadKillCursorsNumberIDs(src []byte) (numIDs int32, rem []byte, ok bool) {
	return readi32(src)
}


func ReadKillCursorsCursorIDs(src []byte, numIDs int32) (cursorIDs []int64, rem []byte, ok bool) {
	var i int32
	var id int64
	for i = 0; i < numIDs; i++ {
		id, src, ok = readi64(src)
		if !ok {
			return cursorIDs, src, false
		}

		cursorIDs = append(cursorIDs, id)
	}
	return cursorIDs, src, true
}

func appendi32(dst []byte, x int32) []byte {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b, uint32(x))
	return append(dst, b...)
}

func appendi64(dst []byte, x int64) []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(b, uint64(x))
	return append(dst, b...)
}

func appendCString(b []byte, str string) []byte {
	b = append(b, str...)
	return append(b, 0x00)
}

func readi32(src []byte) (int32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}
	return readi32unsafe(src), src[4:], true
}

func readi32unsafe(src []byte) int32 {
	return int32(binary.LittleEndian.Uint32(src))
}

func readi64(src []byte) (int64, []byte, bool) {
	if len(src) < 8 {
		return 0, src, false
	}
	return int64(binary.LittleEndian.Uint64(src)), src[8:], true
}

func readcstring(src []byte) (string, []byte, bool) {
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 {
		return "", src, false
	}
	return string(src[:idx]), src[idx+1:], true
}
