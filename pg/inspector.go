package pg

import (
	"context"

	"github.com/fwojciec/pgschema"
	"github.com/lib/pq"
)

// -- Get possible types in a PostgreSQL database
// SELECT
// c.relname AS table_name,
// a.attname AS column_name,
// t.typname AS data_type,
// a.attnotnull AS not_null,
// c.relkind
// FROM pg_catalog.pg_class c
// JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
// JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid
// JOIN pg_catalog.pg_type t ON a.atttypid = t.oid
// WHERE
// c.relname = ANY($1)
// AND n.nspname = $2
// AND a.attnum > 0
// AND NOT a.attisdropped
// ORDER BY c.relname, a.attnum;

type inspector struct {
	db *DB
}

type inspectionResult struct {
	TableName  string `db:"table_name"`
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
	NotNull    bool   `db:"not_null"`
	TableKind  string `db:"relkind"`
}

func (s *inspector) Inspect(ctx context.Context, schema string, tables []string) ([]*pgschema.Table, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const q = `
SELECT
    c.relname AS table_name,
    a.attname AS column_name,
    t.typname AS data_type,
    a.attnotnull AS not_null,
    c.relkind
FROM pg_catalog.pg_class c
JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
JOIN pg_catalog.pg_attribute a ON a.attrelid = c.oid
JOIN pg_catalog.pg_type t ON a.atttypid = t.oid
WHERE
    c.relname = ANY($1)
    AND n.nspname = $2
    AND a.attnum > 0
    AND NOT a.attisdropped
ORDER BY c.relname, a.attnum;
`

	var data []inspectionResult

	if err := tx.SelectContext(ctx, &data, q, pq.Array(tables), schema); err != nil {
		return nil, err
	}

	tableMap := make(map[string]*pgschema.Table)
	tableOrder := make([]string, 0)

	for _, r := range data {
		table, ok := tableMap[r.TableName]
		if !ok {
			table = &pgschema.Table{
				Name:    r.TableName,
				Type:    r.TableKind,
				Columns: []*pgschema.Column{},
			}
			tableMap[r.TableName] = table
			tableOrder = append(tableOrder, r.TableName)
		}
		table.Columns = append(table.Columns, &pgschema.Column{
			Name:    r.ColumnName,
			Type:    r.DataType,
			NotNull: r.NotNull,
		})
	}

	res := make([]*pgschema.Table, len(tableOrder))
	for i, tableName := range tableOrder {
		res[i] = tableMap[tableName]
	}

	return res, nil
}

func NewInspector(db *DB) pgschema.Inspector {
	return &inspector{db: db}
}
