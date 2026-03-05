

package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/status"
)







type PreparedMsg struct {
	
	encodedData mem.BufferSlice
	hdr         []byte
	payload     mem.BufferSlice
	pf          payloadFormat
}


func (p *PreparedMsg) Encode(s Stream, msg any) error {
	ctx := s.Context()
	rpcInfo, ok := rpcInfoFromContext(ctx)
	if !ok {
		return status.Errorf(codes.Internal, "grpc: unable to get rpcInfo")
	}

	
	if rpcInfo.preloaderInfo == nil {
		return status.Errorf(codes.Internal, "grpc: rpcInfo.preloaderInfo is nil")
	}
	if rpcInfo.preloaderInfo.codec == nil {
		return status.Errorf(codes.Internal, "grpc: rpcInfo.preloaderInfo.codec is nil")
	}

	
	data, err := encode(rpcInfo.preloaderInfo.codec, msg)
	if err != nil {
		return err
	}

	materializedData := data.Materialize()
	data.Free()
	p.encodedData = mem.BufferSlice{mem.SliceBuffer(materializedData)}

	
	
	
	var compData mem.BufferSlice
	compData, p.pf, err = compress(p.encodedData, rpcInfo.preloaderInfo.cp, rpcInfo.preloaderInfo.comp, mem.DefaultBufferPool())
	if err != nil {
		return err
	}

	if p.pf.isCompressed() {
		materializedCompData := compData.Materialize()
		compData.Free()
		compData = mem.BufferSlice{mem.SliceBuffer(materializedCompData)}
	}

	p.hdr, p.payload = msgHeader(p.encodedData, compData, p.pf)

	return nil
}
