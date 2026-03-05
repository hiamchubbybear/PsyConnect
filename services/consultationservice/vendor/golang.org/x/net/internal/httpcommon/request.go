



package httpcommon

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptrace"
	"net/textproto"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/http2/hpack"
)

var (
	ErrRequestHeaderListSize = errors.New("request header list larger than peer's advertised limit")
)




type Request struct {
	URL                 *url.URL
	Method              string
	Host                string
	Header              map[string][]string
	Trailer             map[string][]string
	ActualContentLength int64 
}


type EncodeHeadersParam struct {
	Request Request

	
	
	AddGzipHeader bool

	
	PeerMaxHeaderListSize uint64

	
	
	DefaultUserAgent string
}


type EncodeHeadersResult struct {
	HasBody     bool
	HasTrailers bool
}





func EncodeHeaders(ctx context.Context, param EncodeHeadersParam, headerf func(name, value string)) (res EncodeHeadersResult, _ error) {
	req := param.Request

	
	if err := checkConnHeaders(req.Header); err != nil {
		return res, err
	}

	if req.URL == nil {
		return res, errors.New("Request.URL is nil")
	}

	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	host, err := httpguts.PunycodeHostPort(host)
	if err != nil {
		return res, err
	}
	if !httpguts.ValidHostHeader(host) {
		return res, errors.New("invalid Host header")
	}

	
	isNormalConnect := false
	var protocol string
	if vv := req.Header[":protocol"]; len(vv) > 0 {
		protocol = vv[0]
	}
	if req.Method == "CONNECT" && protocol == "" {
		isNormalConnect = true
	} else if protocol != "" && req.Method != "CONNECT" {
		return res, errors.New("invalid :protocol header in non-CONNECT request")
	}

	
	var path string
	if !isNormalConnect {
		path = req.URL.RequestURI()
		if !validPseudoPath(path) {
			orig := path
			path = strings.TrimPrefix(path, req.URL.Scheme+"://"+host)
			if !validPseudoPath(path) {
				if req.URL.Opaque != "" {
					return res, fmt.Errorf("invalid request :path %q from URL.Opaque = %q", orig, req.URL.Opaque)
				} else {
					return res, fmt.Errorf("invalid request :path %q", orig)
				}
			}
		}
	}

	
	
	
	if err := validateHeaders(req.Header); err != "" {
		return res, fmt.Errorf("invalid HTTP header %s", err)
	}
	if err := validateHeaders(req.Trailer); err != "" {
		return res, fmt.Errorf("invalid HTTP trailer %s", err)
	}

	trailers, err := commaSeparatedTrailers(req.Trailer)
	if err != nil {
		return res, err
	}

	enumerateHeaders := func(f func(name, value string)) {
		
		
		
		
		
		f(":authority", host)
		m := req.Method
		if m == "" {
			m = "GET"
		}
		f(":method", m)
		if !isNormalConnect {
			f(":path", path)
			f(":scheme", req.URL.Scheme)
		}
		if protocol != "" {
			f(":protocol", protocol)
		}
		if trailers != "" {
			f("trailer", trailers)
		}

		var didUA bool
		for k, vv := range req.Header {
			if asciiEqualFold(k, "host") || asciiEqualFold(k, "content-length") {
				
				
				continue
			} else if asciiEqualFold(k, "connection") ||
				asciiEqualFold(k, "proxy-connection") ||
				asciiEqualFold(k, "transfer-encoding") ||
				asciiEqualFold(k, "upgrade") ||
				asciiEqualFold(k, "keep-alive") {
				
				
				
				
				continue
			} else if asciiEqualFold(k, "user-agent") {
				
				
				
				
				didUA = true
				if len(vv) < 1 {
					continue
				}
				vv = vv[:1]
				if vv[0] == "" {
					continue
				}
			} else if asciiEqualFold(k, "cookie") {
				
				
				
				for _, v := range vv {
					for {
						p := strings.IndexByte(v, ';')
						if p < 0 {
							break
						}
						f("cookie", v[:p])
						p++
						
						for p+1 <= len(v) && v[p] == ' ' {
							p++
						}
						v = v[p:]
					}
					if len(v) > 0 {
						f("cookie", v)
					}
				}
				continue
			} else if k == ":protocol" {
				
				continue
			}

			for _, v := range vv {
				f(k, v)
			}
		}
		if shouldSendReqContentLength(req.Method, req.ActualContentLength) {
			f("content-length", strconv.FormatInt(req.ActualContentLength, 10))
		}
		if param.AddGzipHeader {
			f("accept-encoding", "gzip")
		}
		if !didUA {
			f("user-agent", param.DefaultUserAgent)
		}
	}

	
	
	
	
	if param.PeerMaxHeaderListSize > 0 {
		hlSize := uint64(0)
		enumerateHeaders(func(name, value string) {
			hf := hpack.HeaderField{Name: name, Value: value}
			hlSize += uint64(hf.Size())
		})

		if hlSize > param.PeerMaxHeaderListSize {
			return res, ErrRequestHeaderListSize
		}
	}

	trace := httptrace.ContextClientTrace(ctx)

	
	enumerateHeaders(func(name, value string) {
		name, ascii := LowerHeader(name)
		if !ascii {
			
			
			return
		}

		headerf(name, value)

		if trace != nil && trace.WroteHeaderField != nil {
			trace.WroteHeaderField(name, []string{value})
		}
	})

	res.HasBody = req.ActualContentLength != 0
	res.HasTrailers = trailers != ""
	return res, nil
}



func IsRequestGzip(method string, header map[string][]string, disableCompression bool) bool {
	
	if !disableCompression &&
		len(header["Accept-Encoding"]) == 0 &&
		len(header["Range"]) == 0 &&
		method != "HEAD" {
		
		
		
		
		
		
		
		
		
		
		
		
		return true
	}
	return false
}








func checkConnHeaders(h map[string][]string) error {
	if vv := h["Upgrade"]; len(vv) > 0 && (vv[0] != "" && vv[0] != "chunked") {
		return fmt.Errorf("invalid Upgrade request header: %q", vv)
	}
	if vv := h["Transfer-Encoding"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && vv[0] != "chunked") {
		return fmt.Errorf("invalid Transfer-Encoding request header: %q", vv)
	}
	if vv := h["Connection"]; len(vv) > 0 && (len(vv) > 1 || vv[0] != "" && !asciiEqualFold(vv[0], "close") && !asciiEqualFold(vv[0], "keep-alive")) {
		return fmt.Errorf("invalid Connection request header: %q", vv)
	}
	return nil
}

func commaSeparatedTrailers(trailer map[string][]string) (string, error) {
	keys := make([]string, 0, len(trailer))
	for k := range trailer {
		k = CanonicalHeader(k)
		switch k {
		case "Transfer-Encoding", "Trailer", "Content-Length":
			return "", fmt.Errorf("invalid Trailer key %q", k)
		}
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		sort.Strings(keys)
		return strings.Join(keys, ","), nil
	}
	return "", nil
}














func validPseudoPath(v string) bool {
	return (len(v) > 0 && v[0] == '/') || v == "*"
}

func validateHeaders(hdrs map[string][]string) string {
	for k, vv := range hdrs {
		if !httpguts.ValidHeaderFieldName(k) && k != ":protocol" {
			return fmt.Sprintf("name %q", k)
		}
		for _, v := range vv {
			if !httpguts.ValidHeaderFieldValue(v) {
				
				
				return fmt.Sprintf("value for header %q", k)
			}
		}
	}
	return ""
}






func shouldSendReqContentLength(method string, contentLength int64) bool {
	if contentLength > 0 {
		return true
	}
	if contentLength < 0 {
		return false
	}
	
	
	switch method {
	case "POST", "PUT", "PATCH":
		return true
	default:
		return false
	}
}


type ServerRequestParam struct {
	Method                  string
	Scheme, Authority, Path string
	Protocol                string
	Header                  map[string][]string
}


type ServerRequestResult struct {
	
	URL        *url.URL
	RequestURI string
	Trailer    map[string][]string

	NeedsContinue bool 

	
	
	
	
	InvalidReason string
}

func NewServerRequest(rp ServerRequestParam) ServerRequestResult {
	needsContinue := httpguts.HeaderValuesContainsToken(rp.Header["Expect"], "100-continue")
	if needsContinue {
		delete(rp.Header, "Expect")
	}
	
	if cookies := rp.Header["Cookie"]; len(cookies) > 1 {
		rp.Header["Cookie"] = []string{strings.Join(cookies, "; ")}
	}

	
	var trailer map[string][]string
	for _, v := range rp.Header["Trailer"] {
		for _, key := range strings.Split(v, ",") {
			key = textproto.CanonicalMIMEHeaderKey(textproto.TrimString(key))
			switch key {
			case "Transfer-Encoding", "Trailer", "Content-Length":
				
				
			default:
				if trailer == nil {
					trailer = make(map[string][]string)
				}
				trailer[key] = nil
			}
		}
	}
	delete(rp.Header, "Trailer")

	
	
	
	if strings.IndexByte(rp.Authority, '@') != -1 && (rp.Scheme == "http" || rp.Scheme == "https") {
		return ServerRequestResult{
			InvalidReason: "userinfo_in_authority",
		}
	}

	var url_ *url.URL
	var requestURI string
	if rp.Method == "CONNECT" && rp.Protocol == "" {
		url_ = &url.URL{Host: rp.Authority}
		requestURI = rp.Authority 
	} else {
		var err error
		url_, err = url.ParseRequestURI(rp.Path)
		if err != nil {
			return ServerRequestResult{
				InvalidReason: "bad_path",
			}
		}
		requestURI = rp.Path
	}

	return ServerRequestResult{
		URL:           url_,
		NeedsContinue: needsContinue,
		RequestURI:    requestURI,
		Trailer:       trailer,
	}
}
