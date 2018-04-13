package structs

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/crypto/blake2b"
)

var (
	validNamespaceName = regexp.MustCompile("^[a-zA-Z0-9-]{1,16}$")
)

const (
	NamespaceUpsertRequestType MessageType = 100 + iota
	NamespaceDeleteRequestType

	QuotaUpsertRequestType
	QuotaDeleteRequestType
	QuotaUsageUpsertRequestType
)

// NamespaceListStub is used to return a subset of job information
// for the job list
type NamespaceListStub struct {
	ID          string
	Name        string
	Quota       string
	Description string
	Hash        []byte
	CreateIndex uint64
	ModifyIndex uint64
}

// NamespaceListRequest is used to parameterize a list request
type Namespace struct {
	Name        string
	Description string
	Quota       string
	Hash        []byte
	CreateTime  time.Time // Time of creation
	CreateIndex uint64
	ModifyIndex uint64
}

func (a *Namespace) Stub() *NamespaceListStub {
	return &NamespaceListStub{
		ID:          a.Name,
		Name:        a.Name,
		Quota:       a.Quota,
		Description: a.Description,
		Hash:        a.Hash,
		CreateIndex: a.CreateIndex,
		ModifyIndex: a.ModifyIndex,
	}
}

func (c *Namespace) SetHash() []byte {
	// Initialize a 256bit Blake2 hash (32 bytes)
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	// Write all the user set fields
	hash.Write([]byte(c.Name))
	hash.Write([]byte(c.Description))

	// Finalize the hash
	hashVal := hash.Sum(nil)

	// Set and return the hash
	c.Hash = hashVal
	return hashVal
}

func (a *Namespace) Validate() error {
	var mErr multierror.Error
	if !validNamespaceName.MatchString(a.Name) {
		err := fmt.Errorf("invalid Namespace '%s'", a.Name)
		mErr.Errors = append(mErr.Errors, err)
	}

	return mErr.ErrorOrNil()
}

//NamespaceListRequest is used to parameterize a list request
type NamespaceListRequest struct {
	QueryOptions
}

//NodeListResponse is used for a list request
type NamespaceListResponse struct {
	Namespaces []*NamespaceListStub
	QueryMeta
}

//NamespaceSpecificRequest is used when we just need to specify a target node
type NamespaceReadRequest struct {
	ID string
	QueryOptions
}

//SingleNamespaceResponse is used to return a single node
type SingleNamespaceResponse struct {
	Namespace *NamespaceListStub
	QueryMeta
}

type NamespaceDeregisterRequest struct {
	ID string
	WriteRequest
}

//NamespaceUpsertRequest is used to upsert a set of policies
type NamespaceUpsertRequest struct {
	Namespace *Namespace
	WriteRequest
}

//////////////////////////////////////

//QuotaSpec specifies the allowed resource usage across regions.
type QuotaSpec struct {
	// Name is the name for the quota object
	Name string

	// Description is an optional description for the quota object
	Description string

	// Limits is the set of quota limits encapsulated by this quota object. Each
	// limit applies quota in a particular region and in the future over a
	// particular priority range and datacenter set.
	Limits []*QuotaLimit

	// Raft indexes to track creation and modification
	CreateIndex uint64
	ModifyIndex uint64
}

func (c *QuotaLimit) SetHash() []byte {
	// Initialize a 256bit Blake2 hash (32 bytes)
	hash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}

	// Write all the user set fields
	hash.Write([]byte(c.Region))
	hash.Write([]byte(strconv.Itoa(c.RegionLimit.CPU)))
	hash.Write([]byte(strconv.Itoa(c.RegionLimit.DiskMB)))
	hash.Write([]byte(strconv.Itoa(c.RegionLimit.IOPS)))
	hash.Write([]byte(strconv.Itoa(c.RegionLimit.MemoryMB)))

	// Finalize the hash
	hashVal := hash.Sum(nil)

	// Set and return the hash
	c.Hash = hashVal
	return hashVal
}

func (c *QuotaSpec) SetHash() {
	for _, k := range c.Limits {
		k.SetHash()
	}
}

func (a *QuotaSpec) Validate() error {
	var mErr multierror.Error
	if !validNamespaceName.MatchString(a.Name) {
		err := fmt.Errorf("invalid Quota '%s'", a.Name)
		mErr.Errors = append(mErr.Errors, err)
	}

	return mErr.ErrorOrNil()
}

//QuotaLimit describes the resource limit in a particular region.
type QuotaLimit struct {
	// Region is the region in which this limit has affect
	Region string

	// RegionLimit is the quota limit that applies to any allocation within a
	// referencing namespace in the region. A value of zero is treated as
	// unlimited and a negative value is treated as fully disallowed. This is
	// useful for once we support GPUs
	RegionLimit *Resources

	// Hash is the hash of the object and is used to make replication efficient.
	Hash []byte
}

func (c *QuotaUsage) SetHash() {
	for _, k := range c.Used {
		k.SetHash()
	}
}

//QuotaUsage is the resource usage of a Quota
type QuotaUsage struct {
	Name        string
	Used        map[string]*QuotaLimit
	CreateIndex uint64
	ModifyIndex uint64
}

//QuotaListRequest is used to parameterize a list request
type QuotaListRequest struct {
	QueryOptions
}

//QuotaListResponse is used for a list request
type QuotaListResponse struct {
	Quotas []*QuotaSpec
	QueryMeta
}

//QuotaUsageListRequest is used to parameterize a list request
type QuotaUsageListRequest struct {
	QueryOptions
}

//QuotaUsageListResponse is used for a list request
type QuotaUsageListResponse struct {
	Quotas []*QuotaUsage
	QueryMeta
}

//QuotaDeleteRequest is used to delete a set of tokens
type QuotaCommonReadRequest struct {
	ID string
	QueryOptions
}

//SingleQuotaResponse is used to return a single node
type SingleQuotaResponse struct {
	Quota *QuotaSpec
	QueryMeta
}

//SingleQuotaResponse is used to return a single node
type SingleQuotaUsageResponse struct {
	Usage *QuotaUsage
	QueryMeta
}

//QuotaDeleteRequest is used to delete a set of tokens
type QuotaDeleteRequest struct {
	ID string
	WriteRequest
}

//QuotaUpsertRequest is used to upsert a set of policies
type QuotaUpsertRequest struct {
	Quota *QuotaSpec
	WriteRequest
}

//QuotaUpsertRequest is used to upsert a set of policies
type QuotaUsageUpsertRequest struct {
	Quota *QuotaUsage
	WriteRequest
}
