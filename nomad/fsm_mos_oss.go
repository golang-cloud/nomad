// +build !pro,!ent

package nomad

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/nomad/nomad/structs"
)

type ns struct {
	*nomadFSM
}

func (n *ns) applyNamespaceUpsert(buf []byte, index uint64) interface{} {

	defer metrics.MeasureSince([]string{"nomad", "fsm", "apply_namespace_upsert"}, time.Now())

	var req structs.NamespaceUpsertRequest

	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := n.state.UpsertNamespace(index, req.Namespace); err != nil {
		n.logger.Printf("[ERR] nomad.fsm: Namespace failed: %v", err)
		return err
	}

	return nil
}

func (n *ns) applyNamespaceDelete(buf []byte, index uint64) interface{} {
	defer metrics.MeasureSince([]string{"nomad", "fsm", "apply_namespace_delete"}, time.Now())
	var req structs.NamespaceDeregisterRequest
	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := n.state.DeleteNamespace(index, req.ID); err != nil {
		n.logger.Printf("[ERR] nomad.fsm: DeleteNamespace failed: %v", err)
		return err
	}
	return nil
}

func (n *ns) applyQuotaUpsert(buf []byte, index uint64) interface{} {

	defer metrics.MeasureSince([]string{"nomad", "fsm", "apply_quota_upsert"}, time.Now())

	var req structs.QuotaUpsertRequest

	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := n.state.QuotaUpsert(index, req.Quota); err != nil {
		n.logger.Printf("[ERR] nomad.fsm: Quota failed: %v", err)
		return err
	}

	return nil
}

func (n *ns) applyQuotaUsageUpsert(buf []byte, index uint64) interface{} {

	defer metrics.MeasureSince([]string{"nomad", "fsm", "apply_quota_usage_upsert"}, time.Now())

	var req structs.QuotaUsageUpsertRequest

	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := n.state.QuotaUsageUpsert(index, req.Quota); err != nil {
		n.logger.Printf("[ERR] nomad.fsm: Quota Usage failed: %v", err)
		return err
	}

	return nil
}

func (n *ns) applyQuotaDelete(buf []byte, index uint64) interface{} {
	defer metrics.MeasureSince([]string{"nomad", "fsm", "apply_quota_delete"}, time.Now())
	var req structs.QuotaDeleteRequest
	if err := structs.Decode(buf, &req); err != nil {
		panic(fmt.Errorf("failed to decode request: %v", err))
	}

	if err := n.state.DeleteQuota(index, req.ID); err != nil {
		n.logger.Printf("[ERR] nomad.fsm: DeleteQuota failed: %v", err)
		return err
	}
	return nil
}
