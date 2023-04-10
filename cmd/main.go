package main

import (
	"bytes"
	"fmt"

	"github.com/fwojciec/pgschema"
)

func main() {
	tables := []*pgschema.Table{
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

	generator := pgschema.NewKyselyGenerator(pgschema.GeneratorOptions{AssumeNonNullableForViews: true})
	reader, err := generator.Generate(tables)
	if err != nil {
		panic(err)
	}

	// Read and print the generated TypeScript code
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
