





package credentials 

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc/attributes"
	icredentials "google.golang.org/grpc/internal/credentials"
	"google.golang.org/protobuf/proto"
)



type PerRPCCredentials interface {
	
	
	
	
	
	
	
	
	
	
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	
	
	RequireTransportSecurity() bool
}




type SecurityLevel int

const (
	
	
	InvalidSecurityLevel SecurityLevel = iota
	
	NoSecurity
	
	IntegrityOnly
	
	PrivacyAndIntegrity
)


func (s SecurityLevel) String() string {
	switch s {
	case NoSecurity:
		return "NoSecurity"
	case IntegrityOnly:
		return "IntegrityOnly"
	case PrivacyAndIntegrity:
		return "PrivacyAndIntegrity"
	}
	return fmt.Sprintf("invalid SecurityLevel: %v", int(s))
}






type CommonAuthInfo struct {
	SecurityLevel SecurityLevel
}


func (c CommonAuthInfo) GetCommonAuthInfo() CommonAuthInfo {
	return c
}



type ProtocolInfo struct {
	
	ProtocolVersion string
	
	SecurityProtocol string
	
	
	
	
	
	SecurityVersion string
	
	ServerName string
}




type AuthInfo interface {
	AuthType() string
}



var ErrConnDispatched = errors.New("credentials: rawConn is dispatched out of gRPC")



type TransportCredentials interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
	
	
	
	
	
	
	ServerHandshake(net.Conn) (net.Conn, AuthInfo, error)
	
	Info() ProtocolInfo
	
	Clone() TransportCredentials
	
	
	
	
	
	
	
	OverrideServerName(string) error
}










type Bundle interface {
	
	
	
	
	
	TransportCredentials() TransportCredentials

	
	
	
	PerRPCCredentials() PerRPCCredentials

	
	
	
	
	NewWithMode(mode string) (Bundle, error)
}




type RequestInfo struct {
	
	Method string
	
	AuthInfo AuthInfo
}




func RequestInfoFromContext(ctx context.Context) (ri RequestInfo, ok bool) {
	ri, ok = icredentials.RequestInfoFromContext(ctx).(RequestInfo)
	return ri, ok
}







type ClientHandshakeInfo struct {
	
	
	Attributes *attributes.Attributes
}





func ClientHandshakeInfoFromContext(ctx context.Context) ClientHandshakeInfo {
	chi, _ := icredentials.ClientHandshakeInfoFromContext(ctx).(ClientHandshakeInfo)
	return chi
}






func CheckSecurityLevel(ai AuthInfo, level SecurityLevel) error {
	type internalInfo interface {
		GetCommonAuthInfo() CommonAuthInfo
	}
	if ai == nil {
		return errors.New("AuthInfo is nil")
	}
	if ci, ok := ai.(internalInfo); ok {
		
		if ci.GetCommonAuthInfo().SecurityLevel == InvalidSecurityLevel {
			return nil
		}
		if ci.GetCommonAuthInfo().SecurityLevel < level {
			return fmt.Errorf("requires SecurityLevel %v; connection has %v", level, ci.GetCommonAuthInfo().SecurityLevel)
		}
	}
	
	return nil
}





type ChannelzSecurityInfo interface {
	GetSecurityValue() ChannelzSecurityValue
}






type ChannelzSecurityValue interface {
	isChannelzSecurityValue()
}







type OtherChannelzSecurityValue struct {
	ChannelzSecurityValue
	Name  string
	Value proto.Message
}
