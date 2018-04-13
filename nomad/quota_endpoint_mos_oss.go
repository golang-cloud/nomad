package nomad

import (
	"fmt"
	"time"

	metrics "github.com/armon/go-metrics"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/nomad/state"
	"github.com/hashicorp/nomad/nomad/structs"
)

type Quota struct {
	srv *Server
}

func (n *Quota) List(args *structs.QuotaListRequest, reply *structs.QuotaListResponse) error {
	if done, err := n.srv.forward("Quota.List", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quota", "list"}, time.Now())

	// Check node read permissions
	if aclObj, err := n.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowOperatorRead() {
		return structs.ErrPermissionDenied
	}

	// Setup the blocking query
	opts := blockingOptions{
		queryOpts: &args.QueryOptions,
		queryMeta: &reply.QueryMeta,
		run: func(ws memdb.WatchSet, state *state.StateStore) error {

			var err error
			var iter memdb.ResultIterator
			if prefix := args.QueryOptions.Prefix; prefix != "" {
				iter, err = state.QuotaByNamePrefix(ws, prefix)
			} else {
				iter, err = state.Quotas(ws)
			}

			if err != nil {
				return err
			}

			var quotas []*structs.QuotaSpec
			for {
				raw := iter.Next()
				if raw == nil {
					break
				}
				node := raw.(*structs.QuotaSpec)
				quotas = append(quotas, node)
			}
			reply.Quotas = quotas

			// Use the last index that affected the jobs table
			index, err := state.Index("quota")
			if err != nil {
				return err
			}
			reply.Index = index

			// Set the query response
			n.srv.setQueryMeta(&reply.QueryMeta)
			return nil
		}}
	return n.srv.blockingRPC(&opts)
}

func (n *Quota) ListUsage(args *structs.QuotaUsageListRequest, reply *structs.QuotaUsageListResponse) error {
	if done, err := n.srv.forward("Quota.ListUsage", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quotaUsage", "list"}, time.Now())

	// Check node read permissions
	if aclObj, err := n.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowOperatorRead() {
		return structs.ErrPermissionDenied
	}

	// Setup the blocking query
	opts := blockingOptions{
		queryOpts: &args.QueryOptions,
		queryMeta: &reply.QueryMeta,
		run: func(ws memdb.WatchSet, state *state.StateStore) error {

			var err error
			var iter memdb.ResultIterator
			if prefix := args.QueryOptions.Prefix; prefix != "" {
				iter, err = state.QuotaUsageByNamePrefix(ws, prefix)
			} else {
				iter, err = state.QuotaUsages(ws)
			}

			if err != nil {
				return err
			}

			var quotas []*structs.QuotaUsage
			for {
				raw := iter.Next()
				if raw == nil {
					break
				}
				node := raw.(*structs.QuotaUsage)
				quotas = append(quotas, node)
			}
			reply.Quotas = quotas

			// Use the last index that affected the jobs table
			index, err := state.Index("quotaUsage")
			if err != nil {
				return err
			}
			reply.Index = index

			// Set the query response
			n.srv.setQueryMeta(&reply.QueryMeta)
			return nil
		}}
	return n.srv.blockingRPC(&opts)
}

func (n *Quota) ReadUsage(args *structs.QuotaCommonReadRequest, reply *structs.SingleQuotaUsageResponse) error {
	if done, err := n.srv.forward("Quota.ReadUsage", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quota_usage", "get"}, time.Now())

	// Check for read-job permissions
	if aclObj, err := n.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowOperatorRead() {
		return structs.ErrPermissionDenied
	}

	// Setup the blocking query
	opts := blockingOptions{
		queryOpts: &args.QueryOptions,
		queryMeta: &reply.QueryMeta,
		run: func(ws memdb.WatchSet, state *state.StateStore) error {
			// Look for the job
			out, err := state.QuotaUsageByID(ws, args.ID)
			if err != nil {
				return err
			}
			// Setup the output
			if out != nil {
				reply.Usage = out
				reply.Index = out.ModifyIndex
			} else {
				// Use the last index that affected the nodes table
				index, err := state.Index("quota_usage")
				if err != nil {
					return err
				}
				reply.Index = index
			}

			// Set the query response
			n.srv.setQueryMeta(&reply.QueryMeta)
			return nil
		}}
	return n.srv.blockingRPC(&opts)
}

func (n *Quota) Read(args *structs.QuotaCommonReadRequest, reply *structs.SingleQuotaResponse) error {
	if done, err := n.srv.forward("Quota.Read", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quota", "get"}, time.Now())

	// Check for read-job permissions
	if aclObj, err := n.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowOperatorRead() {
		return structs.ErrPermissionDenied
	}

	// Setup the blocking query
	opts := blockingOptions{
		queryOpts: &args.QueryOptions,
		queryMeta: &reply.QueryMeta,
		run: func(ws memdb.WatchSet, state *state.StateStore) error {
			// Look for the job
			out, err := state.QuotaByID(ws, args.ID)
			if err != nil {
				return err
			}
			// Setup the output
			if out != nil {
				reply.Quota = out
				reply.Index = out.ModifyIndex
			} else {
				// Use the last index that affected the nodes table
				index, err := state.Index("quota")
				if err != nil {
					return err
				}
				reply.Index = index
			}

			// Set the query response
			n.srv.setQueryMeta(&reply.QueryMeta)
			return nil
		}}
	return n.srv.blockingRPC(&opts)
}

// Upsert is used to create or update a set of namespace
func (a *Quota) Upsert(args *structs.QuotaUpsertRequest, reply *structs.GenericResponse) error {
	args.Region = a.srv.config.AuthoritativeRegion

	if done, err := a.srv.forward("Quota.Upsert", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quota", "upsert"}, time.Now())

	if acl, err := a.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if acl == nil || !acl.IsManagement() {
		println(acl == nil)
		return structs.ErrPermissionDenied
	}

	if err := args.Quota.Validate(); err != nil {
		return fmt.Errorf("Quota invalid: %v", err)
	}
	args.Quota.SetHash()

	// Update via Raft
	_, index, err := a.srv.raftApply(structs.QuotaUpsertRequestType, args)
	if err != nil {
		return err
	}

	// Update the index
	reply.Index = index
	return nil

}

func (a *Quota) Delete(args *structs.QuotaDeleteRequest, reply *structs.GenericResponse) error {
	args.Region = a.srv.config.AuthoritativeRegion

	if done, err := a.srv.forward("Quota.Delete", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "quota", "delete"}, time.Now())

	// Check management level permissions
	if acl, err := a.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if acl == nil || !acl.IsManagement() {
		return structs.ErrPermissionDenied
	}

	// Validate non-zero set of policies
	if len(args.ID) == 0 {
		return fmt.Errorf("must specify as least one id")
	}

	// Update via Raft
	_, index, err := a.srv.raftApply(structs.QuotaDeleteRequestType, args)
	if err != nil {
		return err
	}

	// Update the index
	reply.Index = index
	return nil

}
