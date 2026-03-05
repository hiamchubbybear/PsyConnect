



package http2

import (
	"math"
	"net/http"
	"time"
)
















type http2Config struct {
	MaxConcurrentStreams         uint32
	MaxDecoderHeaderTableSize    uint32
	MaxEncoderHeaderTableSize    uint32
	MaxReadFrameSize             uint32
	MaxUploadBufferPerConnection int32
	MaxUploadBufferPerStream     int32
	SendPingTimeout              time.Duration
	PingTimeout                  time.Duration
	WriteByteTimeout             time.Duration
	PermitProhibitedCipherSuites bool
	CountError                   func(errType string)
}



func configFromServer(h1 *http.Server, h2 *Server) http2Config {
	conf := http2Config{
		MaxConcurrentStreams:         h2.MaxConcurrentStreams,
		MaxEncoderHeaderTableSize:    h2.MaxEncoderHeaderTableSize,
		MaxDecoderHeaderTableSize:    h2.MaxDecoderHeaderTableSize,
		MaxReadFrameSize:             h2.MaxReadFrameSize,
		MaxUploadBufferPerConnection: h2.MaxUploadBufferPerConnection,
		MaxUploadBufferPerStream:     h2.MaxUploadBufferPerStream,
		SendPingTimeout:              h2.ReadIdleTimeout,
		PingTimeout:                  h2.PingTimeout,
		WriteByteTimeout:             h2.WriteByteTimeout,
		PermitProhibitedCipherSuites: h2.PermitProhibitedCipherSuites,
		CountError:                   h2.CountError,
	}
	fillNetHTTPServerConfig(&conf, h1)
	setConfigDefaults(&conf, true)
	return conf
}



func configFromTransport(h2 *Transport) http2Config {
	conf := http2Config{
		MaxEncoderHeaderTableSize: h2.MaxEncoderHeaderTableSize,
		MaxDecoderHeaderTableSize: h2.MaxDecoderHeaderTableSize,
		MaxReadFrameSize:          h2.MaxReadFrameSize,
		SendPingTimeout:           h2.ReadIdleTimeout,
		PingTimeout:               h2.PingTimeout,
		WriteByteTimeout:          h2.WriteByteTimeout,
	}

	
	
	if conf.MaxReadFrameSize < minMaxFrameSize {
		conf.MaxReadFrameSize = minMaxFrameSize
	} else if conf.MaxReadFrameSize > maxFrameSize {
		conf.MaxReadFrameSize = maxFrameSize
	}

	if h2.t1 != nil {
		fillNetHTTPTransportConfig(&conf, h2.t1)
	}
	setConfigDefaults(&conf, false)
	return conf
}

func setDefault[T ~int | ~int32 | ~uint32 | ~int64](v *T, minval, maxval, defval T) {
	if *v < minval || *v > maxval {
		*v = defval
	}
}

func setConfigDefaults(conf *http2Config, server bool) {
	setDefault(&conf.MaxConcurrentStreams, 1, math.MaxUint32, defaultMaxStreams)
	setDefault(&conf.MaxEncoderHeaderTableSize, 1, math.MaxUint32, initialHeaderTableSize)
	setDefault(&conf.MaxDecoderHeaderTableSize, 1, math.MaxUint32, initialHeaderTableSize)
	if server {
		setDefault(&conf.MaxUploadBufferPerConnection, initialWindowSize, math.MaxInt32, 1<<20)
	} else {
		setDefault(&conf.MaxUploadBufferPerConnection, initialWindowSize, math.MaxInt32, transportDefaultConnFlow)
	}
	if server {
		setDefault(&conf.MaxUploadBufferPerStream, 1, math.MaxInt32, 1<<20)
	} else {
		setDefault(&conf.MaxUploadBufferPerStream, 1, math.MaxInt32, transportDefaultStreamFlow)
	}
	setDefault(&conf.MaxReadFrameSize, minMaxFrameSize, maxFrameSize, defaultMaxReadFrameSize)
	setDefault(&conf.PingTimeout, 1, math.MaxInt64, 15*time.Second)
}



func adjustHTTP1MaxHeaderSize(n int64) int64 {
	
	
	const perFieldOverhead = 32 
	const typicalHeaders = 10   
	return n + typicalHeaders*perFieldOverhead
}
