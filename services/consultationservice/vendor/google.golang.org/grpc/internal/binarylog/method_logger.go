

package binarylog

import (
	"context"
	"net"
	"strings"
	"sync/atomic"
	"time"

	binlogpb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type callIDGenerator struct {
	id uint64
}

func (g *callIDGenerator) next() uint64 {
	id := atomic.AddUint64(&g.id, 1)
	return id
}


func (g *callIDGenerator) reset() {
	g.id = 0
}

var idGen callIDGenerator





type MethodLogger interface {
	Log(context.Context, LogEntryConfig)
}



type TruncatingMethodLogger struct {
	headerMaxLen, messageMaxLen uint64

	callID          uint64
	idWithinCallGen *callIDGenerator

	sink Sink 
}





func NewTruncatingMethodLogger(h, m uint64) *TruncatingMethodLogger {
	return &TruncatingMethodLogger{
		headerMaxLen:  h,
		messageMaxLen: m,

		callID:          idGen.next(),
		idWithinCallGen: &callIDGenerator{},

		sink: DefaultSink, 
	}
}




func (ml *TruncatingMethodLogger) Build(c LogEntryConfig) *binlogpb.GrpcLogEntry {
	m := c.toProto()
	timestamp := timestamppb.Now()
	m.Timestamp = timestamp
	m.CallId = ml.callID
	m.SequenceIdWithinCall = ml.idWithinCallGen.next()

	switch pay := m.Payload.(type) {
	case *binlogpb.GrpcLogEntry_ClientHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ClientHeader.GetMetadata())
	case *binlogpb.GrpcLogEntry_ServerHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ServerHeader.GetMetadata())
	case *binlogpb.GrpcLogEntry_Message:
		m.PayloadTruncated = ml.truncateMessage(pay.Message)
	}
	return m
}


func (ml *TruncatingMethodLogger) Log(_ context.Context, c LogEntryConfig) {
	ml.sink.Write(ml.Build(c))
}

func (ml *TruncatingMethodLogger) truncateMetadata(mdPb *binlogpb.Metadata) (truncated bool) {
	if ml.headerMaxLen == maxUInt {
		return false
	}
	var (
		bytesLimit = ml.headerMaxLen
		index      int
	)
	
	
	
	
	for ; index < len(mdPb.Entry); index++ {
		entry := mdPb.Entry[index]
		if entry.Key == "grpc-trace-bin" {
			
			
			continue
		}
		currentEntryLen := uint64(len(entry.GetKey())) + uint64(len(entry.GetValue()))
		if currentEntryLen > bytesLimit {
			break
		}
		bytesLimit -= currentEntryLen
	}
	truncated = index < len(mdPb.Entry)
	mdPb.Entry = mdPb.Entry[:index]
	return truncated
}

func (ml *TruncatingMethodLogger) truncateMessage(msgPb *binlogpb.Message) (truncated bool) {
	if ml.messageMaxLen == maxUInt {
		return false
	}
	if ml.messageMaxLen >= uint64(len(msgPb.Data)) {
		return false
	}
	msgPb.Data = msgPb.Data[:ml.messageMaxLen]
	return true
}





type LogEntryConfig interface {
	toProto() *binlogpb.GrpcLogEntry
}


type ClientHeader struct {
	OnClientSide bool
	Header       metadata.MD
	MethodName   string
	Authority    string
	Timeout      time.Duration
	
	PeerAddr net.Addr
}

func (c *ClientHeader) toProto() *binlogpb.GrpcLogEntry {
	
	
	clientHeader := &binlogpb.ClientHeader{
		Metadata:   mdToMetadataProto(c.Header),
		MethodName: c.MethodName,
		Authority:  c.Authority,
	}
	if c.Timeout > 0 {
		clientHeader.Timeout = durationpb.New(c.Timeout)
	}
	ret := &binlogpb.GrpcLogEntry{
		Type: binlogpb.GrpcLogEntry_EVENT_TYPE_CLIENT_HEADER,
		Payload: &binlogpb.GrpcLogEntry_ClientHeader{
			ClientHeader: clientHeader,
		},
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type ServerHeader struct {
	OnClientSide bool
	Header       metadata.MD
	
	PeerAddr net.Addr
}

func (c *ServerHeader) toProto() *binlogpb.GrpcLogEntry {
	ret := &binlogpb.GrpcLogEntry{
		Type: binlogpb.GrpcLogEntry_EVENT_TYPE_SERVER_HEADER,
		Payload: &binlogpb.GrpcLogEntry_ServerHeader{
			ServerHeader: &binlogpb.ServerHeader{
				Metadata: mdToMetadataProto(c.Header),
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type ClientMessage struct {
	OnClientSide bool
	
	
	Message any
}

func (c *ClientMessage) toProto() *binlogpb.GrpcLogEntry {
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok {
		data, err = proto.Marshal(m)
		if err != nil {
			grpclogLogger.Infof("binarylogging: failed to marshal proto message: %v", err)
		}
	} else if b, ok := c.Message.([]byte); ok {
		data = b
	} else {
		grpclogLogger.Infof("binarylogging: message to log is neither proto.message nor []byte")
	}
	ret := &binlogpb.GrpcLogEntry{
		Type: binlogpb.GrpcLogEntry_EVENT_TYPE_CLIENT_MESSAGE,
		Payload: &binlogpb.GrpcLogEntry_Message{
			Message: &binlogpb.Message{
				Length: uint32(len(data)),
				Data:   data,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ServerMessage struct {
	OnClientSide bool
	
	
	Message any
}

func (c *ServerMessage) toProto() *binlogpb.GrpcLogEntry {
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok {
		data, err = proto.Marshal(m)
		if err != nil {
			grpclogLogger.Infof("binarylogging: failed to marshal proto message: %v", err)
		}
	} else if b, ok := c.Message.([]byte); ok {
		data = b
	} else {
		grpclogLogger.Infof("binarylogging: message to log is neither proto.message nor []byte")
	}
	ret := &binlogpb.GrpcLogEntry{
		Type: binlogpb.GrpcLogEntry_EVENT_TYPE_SERVER_MESSAGE,
		Payload: &binlogpb.GrpcLogEntry_Message{
			Message: &binlogpb.Message{
				Length: uint32(len(data)),
				Data:   data,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ClientHalfClose struct {
	OnClientSide bool
}

func (c *ClientHalfClose) toProto() *binlogpb.GrpcLogEntry {
	ret := &binlogpb.GrpcLogEntry{
		Type:    binlogpb.GrpcLogEntry_EVENT_TYPE_CLIENT_HALF_CLOSE,
		Payload: nil, 
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ServerTrailer struct {
	OnClientSide bool
	Trailer      metadata.MD
	
	Err error
	
	
	PeerAddr net.Addr
}

func (c *ServerTrailer) toProto() *binlogpb.GrpcLogEntry {
	st, ok := status.FromError(c.Err)
	if !ok {
		grpclogLogger.Info("binarylogging: error in trailer is not a status error")
	}
	var (
		detailsBytes []byte
		err          error
	)
	stProto := st.Proto()
	if stProto != nil && len(stProto.Details) != 0 {
		detailsBytes, err = proto.Marshal(stProto)
		if err != nil {
			grpclogLogger.Infof("binarylogging: failed to marshal status proto: %v", err)
		}
	}
	ret := &binlogpb.GrpcLogEntry{
		Type: binlogpb.GrpcLogEntry_EVENT_TYPE_SERVER_TRAILER,
		Payload: &binlogpb.GrpcLogEntry_Trailer{
			Trailer: &binlogpb.Trailer{
				Metadata:      mdToMetadataProto(c.Trailer),
				StatusCode:    uint32(st.Code()),
				StatusMessage: st.Message(),
				StatusDetails: detailsBytes,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type Cancel struct {
	OnClientSide bool
}

func (c *Cancel) toProto() *binlogpb.GrpcLogEntry {
	ret := &binlogpb.GrpcLogEntry{
		Type:    binlogpb.GrpcLogEntry_EVENT_TYPE_CANCEL,
		Payload: nil,
	}
	if c.OnClientSide {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = binlogpb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}



func metadataKeyOmit(key string) bool {
	switch key {
	case "lb-token", ":path", ":authority", "content-encoding", "content-type", "user-agent", "te":
		return true
	case "grpc-trace-bin": 
		return false
	}
	return strings.HasPrefix(key, "grpc-")
}

func mdToMetadataProto(md metadata.MD) *binlogpb.Metadata {
	ret := &binlogpb.Metadata{}
	for k, vv := range md {
		if metadataKeyOmit(k) {
			continue
		}
		for _, v := range vv {
			ret.Entry = append(ret.Entry,
				&binlogpb.MetadataEntry{
					Key:   k,
					Value: []byte(v),
				},
			)
		}
	}
	return ret
}

func addrToProto(addr net.Addr) *binlogpb.Address {
	ret := &binlogpb.Address{}
	switch a := addr.(type) {
	case *net.TCPAddr:
		if a.IP.To4() != nil {
			ret.Type = binlogpb.Address_TYPE_IPV4
		} else if a.IP.To16() != nil {
			ret.Type = binlogpb.Address_TYPE_IPV6
		} else {
			ret.Type = binlogpb.Address_TYPE_UNKNOWN
			
			break
		}
		ret.Address = a.IP.String()
		ret.IpPort = uint32(a.Port)
	case *net.UnixAddr:
		ret.Type = binlogpb.Address_TYPE_UNIX
		ret.Address = a.String()
	default:
		ret.Type = binlogpb.Address_TYPE_UNKNOWN
	}
	return ret
}
