

package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)



var EnableTracing bool



func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") 
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] 
	}
	return m
}




type traceEventLog interface {
	Printf(format string, a ...any)
	Errorf(format string, a ...any)
	Finish()
}




type traceLog interface {
	LazyLog(x fmt.Stringer, sensitive bool)
	LazyPrintf(format string, a ...any)
	SetError()
	SetRecycler(f func(any))
	SetTraceInfo(traceID, spanID uint64)
	SetMaxEvents(m int)
	Finish()
}


type traceInfo struct {
	tr        traceLog
	firstLine firstLine
}




type firstLine struct {
	mu         sync.Mutex
	client     bool 
	remoteAddr net.Addr
	deadline   time.Duration 
}

func (f *firstLine) SetRemoteAddr(addr net.Addr) {
	f.mu.Lock()
	f.remoteAddr = addr
	f.mu.Unlock()
}

func (f *firstLine) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	var line bytes.Buffer
	io.WriteString(&line, "RPC: ")
	if f.client {
		io.WriteString(&line, "to")
	} else {
		io.WriteString(&line, "from")
	}
	fmt.Fprintf(&line, " %v deadline:", f.remoteAddr)
	if f.deadline != 0 {
		fmt.Fprint(&line, f.deadline)
	} else {
		io.WriteString(&line, "none")
	}
	return line.String()
}

const truncateSize = 100

func truncate(x string, l int) string {
	if l > len(x) {
		return x
	}
	return x[:l]
}


type payload struct {
	sent bool 
	msg  any  
	
}

func (p payload) String() string {
	if p.sent {
		return truncate(fmt.Sprintf("sent: %v", p.msg), truncateSize)
	}
	return truncate(fmt.Sprintf("recv: %v", p.msg), truncateSize)
}

type fmtStringer struct {
	format string
	a      []any
}

func (f *fmtStringer) String() string {
	return fmt.Sprintf(f.format, f.a...)
}

type stringer string

func (s stringer) String() string { return string(s) }
