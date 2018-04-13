package state

import memdb "github.com/hashicorp/go-memdb"

func init() {
	// Register all schemas
	RegisterSchemaFactories(namespaceTableSchema)
	RegisterSchemaFactories(quotaTableSchema)
	RegisterSchemaFactories(quotaUsageTableSchema)
}

//----add-----
func namespaceTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "namespace",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func quotaTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "quota",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func quotaUsageTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "quota_usage",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

//----add-----
