package structs

import (
	"fmt"
	"regexp"
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

// NamespaceListRequest is used to parameterize a list request
type NamespaceListRequest struct {
	QueryOptions
}

// NodeListResponse is used for a list request
type NamespaceListResponse struct {
	Namespaces []*NamespaceListStub
	QueryMeta
}

// NamespaceSpecificRequest is used when we just need to specify a target node
type NamespaceReadRequest struct {
	ID string
	QueryOptions
}

// SingleNodeResponse is used to return a single node
type SingleNamespaceResponse struct {
	Namespace *NamespaceListStub
	QueryMeta
}

type NamespaceDeregisterRequest struct {
	ID string
	WriteRequest
}

// ACLPolicyUpsertRequest is used to upsert a set of policies
type NamespaceUpsertRequest struct {
	Namespace *Namespace
	WriteRequest
}
