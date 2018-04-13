// +build !pro,!ent

package state

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/nomad/structs"
)

// Nodes returns an iterator over all the nodes
func (s *StateStore) Namespaces(ws memdb.WatchSet) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	// Walk the entire nodes table
	iter, err := txn.Get("namespace", "id")
	if err != nil {
		return nil, err
	}
	ws.Add(iter.WatchCh())
	return iter, nil
}

//NamespacesByNamePrefix is used to lookup policies by prefix
func (s *StateStore) NamespacesByNamePrefix(ws memdb.WatchSet, prefix string) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	iter, err := txn.Get("namespace", "id_prefix", prefix)
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %v", err)
	}
	ws.Add(iter.WatchCh())

	return iter, nil
}

func (s *StateStore) UpsertNamespace(index uint64, ns *structs.Namespace) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	// Ensure the policy hash is non-nil. This should be done outside the state store
	// for performance reasons, but we check here for defense in depth.
	if len(ns.Hash) == 0 {
		ns.SetHash()
	}

	// Check if the token already exists
	existing, err := txn.First("namespace", "id", ns.Name)
	if err != nil {
		return fmt.Errorf("namespace lookup failed: %v", err)
	}

	// Update all the indexes
	if existing != nil {
		existTK := existing.(*structs.Namespace)
		ns.CreateIndex = existTK.CreateIndex
		ns.ModifyIndex = index

		// Do not allow SecretID or create time to change
		ns.Name = existTK.Name
		ns.CreateTime = existTK.CreateTime

	} else {
		ns.CreateIndex = index
		ns.ModifyIndex = index
	}

	// Update the token
	if err := txn.Insert("namespace", ns); err != nil {
		return fmt.Errorf("upserting namespace failed: %v", err)
	}

	// Update the indexes table
	if err := txn.Insert("index", &IndexEntry{"namespace", index}); err != nil {
		return fmt.Errorf("index update failed: %v", err)
	}
	txn.Commit()
	return nil
}

//DeleteNamespace deletes the tokens with the given accessor ids
func (s *StateStore) DeleteNamespace(index uint64, id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	// Lookup the namespace
	existing, err := txn.First("namespace", "id", id)
	if err != nil {
		return fmt.Errorf("namespace lookup failed: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("namespace not found")
	}

	// Delete the node
	if err := txn.Delete("namespace", existing); err != nil {
		return fmt.Errorf("namespace delete failed: %v", err)
	}
	if err := txn.Insert("index", &IndexEntry{"namespace", index}); err != nil {
		return fmt.Errorf("index update failed: %v", err)
	}
	txn.Commit()
	return nil
}

func (s *StateStore) NamespaceByID(ws memdb.WatchSet, id string) (*structs.Namespace, error) {
	txn := s.db.Txn(false)

	watchCh, existing, err := txn.FirstWatch("namespace", "id", id)
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %v", err)
	}
	ws.Add(watchCh)

	if existing != nil {
		return existing.(*structs.Namespace), nil
	}
	return nil, nil
}

////quota

//QuotaByNamePrefix is used to lookup policies by prefix
func (s *StateStore) QuotaByNamePrefix(ws memdb.WatchSet, prefix string) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	iter, err := txn.Get("quota", "id_prefix", prefix)
	if err != nil {
		return nil, fmt.Errorf("quota lookup failed: %v", err)
	}
	ws.Add(iter.WatchCh())

	return iter, nil
}

//Quotas returns an iterator over all the nodes
func (s *StateStore) Quotas(ws memdb.WatchSet) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	// Walk the entire nodes table
	iter, err := txn.Get("quota", "id")
	if err != nil {
		return nil, err
	}
	ws.Add(iter.WatchCh())
	return iter, nil
}

//QuotaUsageByNamePrefix is used to lookup policies by prefix
func (s *StateStore) QuotaUsageByNamePrefix(ws memdb.WatchSet, prefix string) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	iter, err := txn.Get("quota_usage", "id_prefix", prefix)
	if err != nil {
		return nil, fmt.Errorf("quota_usage lookup failed: %v", err)
	}
	ws.Add(iter.WatchCh())

	return iter, nil
}

//QuotaUsages returns an iterator over all the nodes
func (s *StateStore) QuotaUsages(ws memdb.WatchSet) (memdb.ResultIterator, error) {
	txn := s.db.Txn(false)

	// Walk the entire nodes table
	iter, err := txn.Get("quota_usage", "id")
	if err != nil {
		return nil, err
	}
	ws.Add(iter.WatchCh())
	return iter, nil
}

func (s *StateStore) QuotaUsageByID(ws memdb.WatchSet, id string) (*structs.QuotaUsage, error) {
	txn := s.db.Txn(false)

	watchCh, existing, err := txn.FirstWatch("quota_usage", "id", id)
	if err != nil {
		return nil, fmt.Errorf("quota_usage lookup failed: %v", err)
	}
	ws.Add(watchCh)

	if existing != nil {
		return existing.(*structs.QuotaUsage), nil
	}
	return nil, nil
}

func (s *StateStore) QuotaByID(ws memdb.WatchSet, id string) (*structs.QuotaSpec, error) {
	txn := s.db.Txn(false)

	watchCh, existing, err := txn.FirstWatch("quota", "id", id)
	if err != nil {
		return nil, fmt.Errorf("quota lookup failed: %v", err)
	}
	ws.Add(watchCh)

	if existing != nil {
		return existing.(*structs.QuotaSpec), nil
	}
	return nil, nil
}

//DeleteQuota deletes the tokens with the given accessor ids
func (s *StateStore) DeleteQuota(index uint64, id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	// Lookup the namespace
	existing, err := txn.First("quota", "id", id)
	if err != nil {
		return fmt.Errorf("quota lookup failed: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("quota not found")
	}

	// Delete the node
	if err := txn.Delete("quota", existing); err != nil {
		return fmt.Errorf("quota delete failed: %v", err)
	}

	// Lookup the quota_usage
	usage, _ := txn.First("quota_usage", "id", id)

	if usage != nil {
		// Delete the quota_usage
		if err := txn.Delete("quota_usage", existing); err != nil {
			return fmt.Errorf("quota_usage delete failed: %v", err)
		}
	}

	if err := txn.Insert("index", &IndexEntry{"quota", index}); err != nil {
		return fmt.Errorf("index update failed: %v", err)
	}
	txn.Commit()
	return nil
}

func (s *StateStore) QuotaUpsert(index uint64, ns *structs.QuotaSpec) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	ns.SetHash()

	// Check if the token already exists
	existing, err := txn.First("quota", "id", ns.Name)
	if err != nil {
		return fmt.Errorf("quota lookup failed: %v", err)
	}

	// Check if the token already exists
	usage, err := txn.First("quota_usage", "id", ns.Name)

	if err != nil {
		return fmt.Errorf("quota_usage lookup failed: %v", err)
	}

	// Update all the indexes
	if existing != nil {
		existTK := existing.(*structs.QuotaSpec)
		ns.CreateIndex = existTK.CreateIndex
		ns.ModifyIndex = index

		// Do not allow SecretID or create time to change
		ns.Name = existTK.Name

	} else {
		ns.CreateIndex = index
		ns.ModifyIndex = index
	}

	if usage == nil {
		usage = &structs.QuotaUsage{Name: ns.Name}
		// Update the token
		if err := txn.Insert("quota_usage", usage); err != nil {
			return fmt.Errorf("upserting quota_usage failed: %v", err)
		}
	}

	// Update the token
	if err := txn.Insert("quota", ns); err != nil {
		return fmt.Errorf("upserting quota failed: %v", err)
	}

	// Update the indexes table
	if err := txn.Insert("index", &IndexEntry{"quota", index}); err != nil {
		return fmt.Errorf("index update failed: %v", err)
	}
	txn.Commit()
	return nil
}

//QuotaUsageUpsert 叠加资源使用
func (s *StateStore) QuotaUsageUpsert(index uint64, add *structs.QuotaUsage) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	// Check if the token already exists
	usage, err := txn.First("quota_usage", "id", add.Name)

	if err != nil || usage == nil {
		return fmt.Errorf("quota_usage lookup failed: %v", err)
	}

	existTK := usage.(*structs.QuotaUsage)
	existTK.ModifyIndex = index

	//
	used := existTK.Used
	for k, v := range add.Used {
		if existUsed, ok := used[k]; ok {
			existUsed.RegionLimit.Add(v.RegionLimit)
		} else {
			used[k] = v
		}
	}

	existTK.SetHash()
	// Update the token
	if err := txn.Insert("quota_usage", existTK); err != nil {
		return fmt.Errorf("upserting quota_usage failed: %v", err)
	}

	txn.Commit()
	return nil
}
