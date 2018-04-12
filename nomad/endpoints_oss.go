// +build !pro,!ent

package nomad

import "net/rpc"

// EnterpriseEndpoints holds the set of enterprise only endpoints to register
type EnterpriseEndpoints struct {
	s *Server
}

// NewEnterpriseEndpoints returns a stub of the enterprise endpoints since there
// are none in oss
func NewEnterpriseEndpoints(s *Server) *EnterpriseEndpoints {
	return &EnterpriseEndpoints{s}
}

// Register is a no-op in oss.
func (e *EnterpriseEndpoints) Register(s *rpc.Server) {
	s.Register(&Namespace{e.s})
}
