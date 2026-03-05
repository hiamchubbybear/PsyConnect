



package insecure

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)





func NewCredentials() credentials.TransportCredentials {
	return insecureTC{}
}




type insecureTC struct{}

func (insecureTC) ClientHandshake(_ context.Context, _ string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, info{credentials.CommonAuthInfo{SecurityLevel: credentials.NoSecurity}}, nil
}

func (insecureTC) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, info{credentials.CommonAuthInfo{SecurityLevel: credentials.NoSecurity}}, nil
}

func (insecureTC) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{SecurityProtocol: "insecure"}
}

func (insecureTC) Clone() credentials.TransportCredentials {
	return insecureTC{}
}

func (insecureTC) OverrideServerName(string) error {
	return nil
}



type info struct {
	credentials.CommonAuthInfo
}


func (info) AuthType() string {
	return "insecure"
}




type insecureBundle struct{}


func NewBundle() credentials.Bundle {
	return insecureBundle{}
}


func (insecureBundle) NewWithMode(string) (credentials.Bundle, error) {
	return insecureBundle{}, nil
}



func (insecureBundle) PerRPCCredentials() credentials.PerRPCCredentials {
	return nil
}


func (insecureBundle) TransportCredentials() credentials.TransportCredentials {
	return NewCredentials()
}
