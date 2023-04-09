package pg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/fwojciec/pgschema"
	"github.com/fwojciec/pgschema/pg"
	"github.com/fwojciec/ta"
)

func TestInspector(t *testing.T) {
	t.Parallel()

	tdb := pg.NewDB(TEST_DSN)
	if err := tdb.Open(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { tdb.Close() })

	isolate(t, tdb, func(t *testing.T, db *pg.DB, schema string) {
		ctx := context.Background()
		createTestTable(t, ctx, db, "test_table", []string{
			`id SERIAL PRIMARY KEY`,
			`name VARCHAR(255) `,
			`created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		})
		createTestView(t, ctx, db, "test_view", `SELECT * FROM test_table`)
		createTestMaterializedView(t, ctx, db, "test_materialized_view", `SELECT * FROM test_table`)
		inspector := pg.NewInspector(db)

		tables, err := inspector.Inspect(ctx, schema, []string{"test_table", "test_view", "test_materialized_view"})

		ta.OK(t, err)
		expected := []*pgschema.Table{
			{
				Name: "test_materialized_view",
				Type: "m",
				Columns: []*pgschema.Column{
					{
						Name:    "id",
						Type:    "int4",
						NotNull: false,
					},
					{
						Name:    "name",
						Type:    "varchar",
						NotNull: false,
					},
					{
						Name:    "created_at",
						Type:    "timestamptz",
						NotNull: false,
					},
				},
			},
			{
				Name: "test_table",
				Type: "r",
				Columns: []*pgschema.Column{
					{
						Name:    "id",
						Type:    "int4",
						NotNull: true,
					},
					{
						Name:    "name",
						Type:    "varchar",
						NotNull: false,
					},
					{
						Name:    "created_at",
						Type:    "timestamptz",
						NotNull: true,
					},
				},
			},
			{
				Name: "test_view",
				Type: "v",
				Columns: []*pgschema.Column{
					{
						Name:    "id",
						Type:    "int4",
						NotNull: false,
					},
					{
						Name:    "name",
						Type:    "varchar",
						NotNull: false,
					},
					{
						Name:    "created_at",
						Type:    "timestamptz",
						NotNull: false,
					},
				},
			},
		}

		ta.Equals(t, expected, tables)
	})

}

func createTestTable(tb testing.TB, ctx context.Context, db *pg.DB, name string, columns []string) {
	tb.Helper()

	tx, err := db.BeginTx(ctx)
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(fmt.Sprintf(`CREATE TABLE %s (%s)`, name, strings.Join(columns, ","))); err != nil {
		tb.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
}

func createTestView(tb testing.TB, ctx context.Context, db *pg.DB, name string, query string) {
	tb.Helper()

	tx, err := db.BeginTx(ctx)
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(fmt.Sprintf(`CREATE VIEW %s AS %s`, name, query)); err != nil {
		tb.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
}

func createTestMaterializedView(tb testing.TB, ctx context.Context, db *pg.DB, name string, query string) {
	tb.Helper()

	tx, err := db.BeginTx(ctx)
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(fmt.Sprintf(`CREATE MATERIALIZED VIEW %s AS %s`, name, query)); err != nil {
		tb.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
}
