





package scram

import "sync"




type Server struct {
	sync.RWMutex
	credentialCB CredentialLookup
	nonceGen     NonceGeneratorFcn
	hashGen      HashGeneratorFcn
}

func newServer(cl CredentialLookup, fcn HashGeneratorFcn) (*Server, error) {
	return &Server{
		credentialCB: cl,
		nonceGen:     defaultNonceGenerator,
		hashGen:      fcn,
	}, nil
}




func (s *Server) WithNonceGenerator(ng NonceGeneratorFcn) *Server {
	s.Lock()
	defer s.Unlock()
	s.nonceGen = ng
	return s
}




func (s *Server) NewConversation() *ServerConversation {
	s.RLock()
	defer s.RUnlock()
	return &ServerConversation{
		nonceGen:     s.nonceGen,
		hashGen:      s.hashGen,
		credentialCB: s.credentialCB,
	}
}
