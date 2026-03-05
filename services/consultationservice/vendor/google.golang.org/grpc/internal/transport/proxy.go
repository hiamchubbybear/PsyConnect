

package transport

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/proxyattributes"
	"google.golang.org/grpc/resolver"
)

const proxyAuthHeaderKey = "Proxy-Authorization"





type bufConn struct {
	net.Conn
	r io.Reader
}

func (c *bufConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func doHTTPConnectHandshake(ctx context.Context, conn net.Conn, grpcUA string, opts proxyattributes.Options) (_ net.Conn, err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: opts.ConnectAddr},
		Header: map[string][]string{"User-Agent": {grpcUA}},
	}
	if user := opts.User; user != nil {
		u := user.Username()
		p, _ := user.Password()
		req.Header.Add(proxyAuthHeaderKey, "Basic "+basicAuth(u, p))
	}
	if err := sendHTTPRequest(ctx, req, conn); err != nil {
		return nil, fmt.Errorf("failed to write the HTTP request: %v", err)
	}

	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, req)
	if err != nil {
		return nil, fmt.Errorf("reading server HTTP response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to do connect handshake, status code: %s", resp.Status)
		}
		return nil, fmt.Errorf("failed to do connect handshake, response: %q", dump)
	}
	
	
	
	
	if r.Buffered() != 0 {
		return &bufConn{Conn: conn, r: r}, nil
	}
	return conn, nil
}


func proxyDial(ctx context.Context, addr resolver.Address, grpcUA string, opts proxyattributes.Options) (net.Conn, error) {
	conn, err := internal.NetDialerWithTCPKeepalive().DialContext(ctx, "tcp", addr.Addr)
	if err != nil {
		return nil, err
	}
	return doHTTPConnectHandshake(ctx, conn, grpcUA, opts)
}

func sendHTTPRequest(ctx context.Context, req *http.Request, conn net.Conn) error {
	req = req.WithContext(ctx)
	if err := req.Write(conn); err != nil {
		return fmt.Errorf("failed to write the HTTP request: %v", err)
	}
	return nil
}
