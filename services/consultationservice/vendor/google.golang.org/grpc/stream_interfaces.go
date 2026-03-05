

package grpc




type ServerStreamingClient[Res any] interface {
	
	
	
	
	
	Recv() (*Res, error)

	
	
	
	ClientStream
}







type ServerStreamingServer[Res any] interface {
	
	
	
	
	Send(*Res) error

	
	
	
	ServerStream
}





type ClientStreamingClient[Req any, Res any] interface {
	
	
	
	
	
	Send(*Req) error

	
	
	
	
	CloseAndRecv() (*Res, error)

	
	
	
	ClientStream
}









type ClientStreamingServer[Req any, Res any] interface {
	
	
	
	
	
	
	Recv() (*Req, error)

	
	
	
	
	SendAndClose(*Res) error

	
	
	
	ServerStream
}





type BidiStreamingClient[Req any, Res any] interface {
	
	
	
	
	
	Send(*Req) error

	
	
	
	
	
	Recv() (*Res, error)

	
	
	
	ClientStream
}








type BidiStreamingServer[Req any, Res any] interface {
	
	
	
	
	
	
	Recv() (*Req, error)

	
	
	
	
	Send(*Res) error

	
	
	
	ServerStream
}



type GenericClientStream[Req any, Res any] struct {
	ClientStream
}

var _ ServerStreamingClient[string] = (*GenericClientStream[int, string])(nil)
var _ ClientStreamingClient[int, string] = (*GenericClientStream[int, string])(nil)
var _ BidiStreamingClient[int, string] = (*GenericClientStream[int, string])(nil)




func (x *GenericClientStream[Req, Res]) Send(m *Req) error {
	return x.ClientStream.SendMsg(m)
}




func (x *GenericClientStream[Req, Res]) Recv() (*Res, error) {
	m := new(Res)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}




func (x *GenericClientStream[Req, Res]) CloseAndRecv() (*Res, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Res)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}



type GenericServerStream[Req any, Res any] struct {
	ServerStream
}

var _ ServerStreamingServer[string] = (*GenericServerStream[int, string])(nil)
var _ ClientStreamingServer[int, string] = (*GenericServerStream[int, string])(nil)
var _ BidiStreamingServer[int, string] = (*GenericServerStream[int, string])(nil)




func (x *GenericServerStream[Req, Res]) Send(m *Res) error {
	return x.ServerStream.SendMsg(m)
}




func (x *GenericServerStream[Req, Res]) SendAndClose(m *Res) error {
	return x.ServerStream.SendMsg(m)
}




func (x *GenericServerStream[Req, Res]) Recv() (*Req, error) {
	m := new(Req)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
