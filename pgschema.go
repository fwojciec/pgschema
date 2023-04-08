package pgschema

import "context"

// Table is a collection of columns.
type Table struct {
	Name    string
	Type    string
	Columns []*Column
}

// Column is a column in a table.
type Column struct {
	Name    string
	Type    string
	NotNull bool
}

// Inspector inspects a database and returns a schema with information about
// specified tables.
type Inspector interface {
	Inspect(ctx context.Context, schema string, tables []string) ([]*Table, error)
}
