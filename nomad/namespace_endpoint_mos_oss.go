package nomad

import (
	"fmt"
	"time"

	metrics "github.com/armon/go-metrics"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/nomad/state"
	"github.com/hashicorp/nomad/nomad/structs"
)

type Namespace struct {
	srv *Server
}

func (n *Namespace) List(args *structs.NamespaceListRequest, reply *structs.NamespaceListResponse) error {
	if done, err := n.srv.forward("Namespace.List", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "namespace", "list"}, time.Now())

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

			iter, err := state.Namespaces(ws)

			if err != nil {
				return err
			}

			var nodes []*structs.NamespaceListStub
			for {
				raw := iter.Next()
				if raw == nil {
					break
				}
				node := raw.(*structs.Namespace)
				nodes = append(nodes, node.Stub())
			}
			reply.Namespaces = nodes

			// Use the last index that affected the jobs table
			index, err := state.Index("namespace")
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

func (n *Namespace) Read(args *structs.NamespaceReadRequest, reply *structs.SingleNamespaceResponse) error {
	if done, err := n.srv.forward("Namespace.Read", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "namespace", "get"}, time.Now())

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
			out, err := state.NamespaceByID(ws, args.ID)
			if err != nil {
				return err
			}
			// Setup the output
			if out != nil {
				reply.Namespace = out.Stub()
				reply.Index = out.ModifyIndex
			} else {
				// Use the last index that affected the nodes table
				index, err := state.Index("namespace")
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
func (a *Namespace) Upsert(args *structs.NamespaceUpsertRequest, reply *structs.GenericResponse) error {
	args.Region = a.srv.config.AuthoritativeRegion

	if done, err := a.srv.forward("Namespace.Upsert", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "namespace", "upsert"}, time.Now())

	if acl, err := a.srv.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if acl == nil || !acl.IsManagement() {
		println(acl == nil)
		return structs.ErrPermissionDenied
	}

	if err := args.Namespace.Validate(); err != nil {
		return fmt.Errorf("Namespace invalid: %v", err)
	}
	args.Namespace.SetHash()

	// Update via Raft
	_, index, err := a.srv.raftApply(structs.NamespaceUpsertRequestType, args)
	if err != nil {
		return err
	}

	// Update the index
	reply.Index = index
	return nil

}

func (a *Namespace) Delete(args *structs.NamespaceDeregisterRequest, reply *structs.GenericResponse) error {
	args.Region = a.srv.config.AuthoritativeRegion

	if done, err := a.srv.forward("Namespace.Delete", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"nomad", "namespace", "delete"}, time.Now())

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
	_, index, err := a.srv.raftApply(structs.NamespaceDeleteRequestType, args)
	if err != nil {
		return err
	}

	// Update the index
	reply.Index = index
	return nil

}
