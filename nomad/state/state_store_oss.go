// +build !pro,!ent

package state

import (
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/nomad/structs"
)

// namespaceExists returns whether a namespace exists
func (s *StateStore) namespaceExists(txn *memdb.Txn, namespace string) (bool, error) {
	existing, err := txn.First("namespace", "id", namespace)

	if err != nil {
		return false, err
	}

	return existing != nil, nil
}

// updateEntWithAlloc is used to update Nomad Enterprise objects when an allocation is
// added/modified/deleted
func (s *StateStore) updateEntWithAlloc(index uint64, new, existing *structs.Allocation, txn *memdb.Txn) error {
	return nil
}
