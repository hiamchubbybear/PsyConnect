package redis

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"

	"github.com/redis/go-redis/v9/internal"
	"github.com/redis/go-redis/v9/internal/pool"
	"github.com/redis/go-redis/v9/internal/proto"
)


var ErrClosed = pool.ErrClosed



var ErrPoolExhausted = pool.ErrPoolExhausted


var ErrPoolTimeout = pool.ErrPoolTimeout





var ErrCrossSlot = proto.RedisError("CROSSSLOT Keys in request don't hash to the same slot")


func HasErrorPrefix(err error, prefix string) bool {
	var rErr Error
	if !errors.As(err, &rErr) {
		return false
	}
	msg := rErr.Error()
	msg = strings.TrimPrefix(msg, "ERR ") 
	return strings.HasPrefix(msg, prefix)
}

type Error interface {
	error

	
	
	
	
	RedisError()
}

var _ Error = proto.RedisError("")

func isContextError(err error) bool {
	switch err {
	case context.Canceled, context.DeadlineExceeded:
		return true
	default:
		return false
	}
}

func shouldRetry(err error, retryTimeout bool) bool {
	switch err {
	case io.EOF, io.ErrUnexpectedEOF:
		return true
	case nil, context.Canceled, context.DeadlineExceeded:
		return false
	case pool.ErrPoolTimeout:
		
		return true
	}

	if v, ok := err.(timeoutError); ok {
		if v.Timeout() {
			return retryTimeout
		}
		return true
	}

	s := err.Error()
	if s == "ERR max number of clients reached" {
		return true
	}
	if strings.HasPrefix(s, "LOADING ") {
		return true
	}
	if strings.HasPrefix(s, "READONLY ") {
		return true
	}
	if strings.HasPrefix(s, "MASTERDOWN ") {
		return true
	}
	if strings.HasPrefix(s, "CLUSTERDOWN ") {
		return true
	}
	if strings.HasPrefix(s, "TRYAGAIN ") {
		return true
	}

	return false
}

func isRedisError(err error) bool {
	_, ok := err.(proto.RedisError)
	return ok
}

func isBadConn(err error, allowTimeout bool, addr string) bool {
	switch err {
	case nil:
		return false
	case context.Canceled, context.DeadlineExceeded:
		return true
	}

	if isRedisError(err) {
		switch {
		case isReadOnlyError(err):
			
			
			return true
		case isMovedSameConnAddr(err, addr):
			
			
			
			return true
		default:
			return false
		}
	}

	if allowTimeout {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false
		}
	}

	return true
}

func isMovedError(err error) (moved bool, ask bool, addr string) {
	if !isRedisError(err) {
		return
	}

	s := err.Error()
	switch {
	case strings.HasPrefix(s, "MOVED "):
		moved = true
	case strings.HasPrefix(s, "ASK "):
		ask = true
	default:
		return
	}

	ind := strings.LastIndex(s, " ")
	if ind == -1 {
		return false, false, ""
	}

	addr = s[ind+1:]
	addr = internal.GetAddr(addr)
	return
}

func isLoadingError(err error) bool {
	return strings.HasPrefix(err.Error(), "LOADING ")
}

func isReadOnlyError(err error) bool {
	return strings.HasPrefix(err.Error(), "READONLY ")
}

func isMovedSameConnAddr(err error, addr string) bool {
	redisError := err.Error()
	if !strings.HasPrefix(redisError, "MOVED ") {
		return false
	}
	return strings.HasSuffix(redisError, " "+addr)
}



type timeoutError interface {
	Timeout() bool
}
