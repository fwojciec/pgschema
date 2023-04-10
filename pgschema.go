package pgschema

import (
	"context"
	"io"
)

// Table is a collection of columns.
type Table struct {
	Name    string    // Name of the table.
	Type    string    // Type of the table ('r' for regular, 'v' for view, and 'm' for materialized view).
	Columns []*Column // Columns in the table.
}

// Column is a column in a table.
type Column struct {
	Name    string // Name of the column.
	Type    string // Type of the column, based on PostgreSQL pg_type.typname.
	NotNull bool   // True if the column is not nullable.
}

// Inspector inspects a database and returns a schema with information about
// specified tables.
type Inspector interface {
	Inspect(ctx context.Context, schema string, tables []string) ([]*Table, error)
}

// Generator generates code from a list of tables.
type Generator interface {
	Generate(tables []*Table) (io.Reader, error)
}
