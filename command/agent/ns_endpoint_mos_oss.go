// +build !pro,!ent

package agent

import (
	"net/http"
	"strings"

	"github.com/hashicorp/nomad/nomad/structs"
)

func (s *HTTPServer) namespacesRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	args := structs.NamespaceListRequest{}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.NamespaceListResponse
	if err := s.agent.RPC("Namespace.List", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Namespaces == nil {
		out.Namespaces = make([]*structs.NamespaceListStub, 0)
	}
	return out.Namespaces, nil
}

func (s *HTTPServer) namespaceSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	ID := strings.TrimPrefix(req.URL.Path, "/v1/namespace/")
	switch req.Method {
	case "GET":
		return s.namespaceRead(resp, req, ID)
	case "POST":
		return s.namespaceUpdate(resp, req, ID)
	case "DELETE":
		return s.namespaceDelete(resp, req, ID)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}

func (s *HTTPServer) namespaceRead(resp http.ResponseWriter, req *http.Request, ID string) (interface{}, error) {

	args := structs.NamespaceReadRequest{
		ID: ID,
	}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.SingleNamespaceResponse
	if err := s.agent.RPC("Namespace.Read", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Namespace == nil {
		return nil, CodedError(404, "Namespace not found")
	}
	return out.Namespace, nil
}

func (s *HTTPServer) namespaceCreate(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var namespace structs.Namespace

	if err := decodeBody(req, &namespace); err != nil {
		return nil, CodedError(500, err.Error())
	}
	// Format the request
	args := structs.NamespaceUpsertRequest{
		Namespace: &namespace,
	}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("Namespace.Upsert", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}

func (s *HTTPServer) namespaceUpdate(resp http.ResponseWriter, req *http.Request, ID string) (interface{}, error) {

	var namespace structs.Namespace
	if err := decodeBody(req, &namespace); err != nil {
		return nil, CodedError(500, err.Error())
	}

	// Ensure the policy name matches
	if namespace.Name != ID {
		return nil, CodedError(400, "Namespace name does not match request path")
	}

	// Format the request
	args := structs.NamespaceUpsertRequest{
		Namespace: &namespace,
	}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("Namespace.Upsert", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}

func (s *HTTPServer) namespaceDelete(resp http.ResponseWriter, req *http.Request, ID string) (interface{}, error) {

	args := structs.NamespaceDeregisterRequest{
		ID: ID,
	}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("Namespace.Delete", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}
