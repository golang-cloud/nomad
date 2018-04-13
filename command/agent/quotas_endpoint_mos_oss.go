// +build !pro,!ent

package agent

import (
	"net/http"
	"strings"

	"github.com/hashicorp/nomad/nomad/structs"
)

func (s *HTTPServer) quotasRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	args := structs.QuotaListRequest{}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.QuotaListResponse
	if err := s.agent.RPC("Quota.List", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Quotas == nil {
		out.Quotas = make([]*structs.QuotaSpec, 0)
	}
	return out.Quotas, nil
}

func (s *HTTPServer) quotaUsages(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}

	args := structs.QuotaUsageListRequest{}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.QuotaUsageListResponse
	if err := s.agent.RPC("Quota.UsageList", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Quotas == nil {
		out.Quotas = make([]*structs.QuotaUsage, 0)
	}
	return out.Quotas, nil
}

func (s *HTTPServer) quotaUsageRead(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		return nil, CodedError(405, ErrInvalidMethod)
	}
	ID := strings.TrimPrefix(req.URL.Path, "/v1/quota/usage/")

	args := structs.QuotaCommonReadRequest{
		ID: ID,
	}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.SingleQuotaUsageResponse
	if err := s.agent.RPC("Quota.ReadUsage", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Usage == nil {
		return nil, CodedError(404, "Quota not found")
	}
	return out.Usage, nil
}

func (s *HTTPServer) quotaRead(resp http.ResponseWriter, req *http.Request, id string) (interface{}, error) {
	args := structs.QuotaCommonReadRequest{
		ID: id,
	}
	if s.parse(resp, req, &args.Region, &args.QueryOptions) {
		return nil, nil
	}

	var out structs.SingleQuotaResponse
	if err := s.agent.RPC("Quota.Read", &args, &out); err != nil {
		return nil, err
	}

	setMeta(resp, &out.QueryMeta)
	if out.Quota == nil {
		return nil, CodedError(404, "Quota not found")
	}
	return out.Quota, nil
}

func (s *HTTPServer) quotaCreate(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var quota structs.QuotaSpec

	if err := decodeBody(req, &quota); err != nil {
		return nil, CodedError(500, err.Error())
	}
	// Format the request
	args := structs.QuotaUpsertRequest{
		Quota: &quota,
	}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("Quota.Upsert", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}

func (s *HTTPServer) quotaUpdate(resp http.ResponseWriter, req *http.Request, id string) (interface{}, error) {

	return nil, nil
}

func (s *HTTPServer) quotaDelete(resp http.ResponseWriter, req *http.Request, id string) (interface{}, error) {
	args := structs.QuotaDeleteRequest{ID: id}
	s.parseWriteRequest(req, &args.WriteRequest)

	var out structs.GenericResponse
	if err := s.agent.RPC("Quota.Delete", &args, &out); err != nil {
		return nil, err
	}
	setIndex(resp, out.Index)
	return nil, nil
}

func (s *HTTPServer) quotaSpecificRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	ID := strings.TrimPrefix(req.URL.Path, "/v1/quota/")
	switch req.Method {
	case "GET":
		return s.quotaRead(resp, req, ID)
	case "POST":
		return s.quotaUpdate(resp, req, ID)
	case "DELETE":
		return s.quotaDelete(resp, req, ID)
	default:
		return nil, CodedError(405, ErrInvalidMethod)
	}
}
