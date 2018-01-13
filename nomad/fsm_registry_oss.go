// +build !pro,!ent

package nomad

import (
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/raft"
	"github.com/ugorji/go/codec"
)

// registerLogAppliers is a no-op for open-source only FSMs.
func (n *nomadFSM) registerLogAppliers() {
	ns := &ns{n}
	n.enterpriseAppliers[structs.NamespaceUpsertRequestType] = ns.applyNamespaceUpsert
	n.enterpriseAppliers[structs.NamespaceDeleteRequestType] = ns.applyNamespaceDelete
}

// registerSnapshotRestorers is a no-op for open-source only FSMs.
func (n *nomadFSM) registerSnapshotRestorers() {}

// persistEnterpriseTables is a no-op for open-source only FSMs.
func (s *nomadSnapshot) persistEnterpriseTables(sink raft.SnapshotSink, encoder *codec.Encoder) error {
	return nil
}
