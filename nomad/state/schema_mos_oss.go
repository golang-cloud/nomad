package state

import memdb "github.com/hashicorp/go-memdb"

func init() {
	// Register all schemas
	RegisterSchemaFactories(namespaceTableSchema)
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

//----add-----
