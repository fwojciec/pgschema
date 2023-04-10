package pgschema

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type kyselyGenerator struct {
	assumeNonNullableForViews bool
}

type GeneratorOptions struct {
	AssumeNonNullableForViews bool
}

func NewKyselyGenerator(opts GeneratorOptions) Generator {
	return &kyselyGenerator{assumeNonNullableForViews: opts.AssumeNonNullableForViews}
}

var tsTemplate = `{{- $pgTypeToTSType := .pgTypeToTSType -}}
{{- $nullableTSType := .nullableTSType -}}
{{- $assumeNonNullableForViews := .assumeNonNullableForViews -}}
{{- $first := true -}}
{{- range .tables -}}
{{- $table := . -}}
{{ if $first}}{{ $first = false }}{{ else -}}
{{ "\n\n" }}{{ end -}}
export interface {{ .Name | title }} {
{{- range .Columns }}
  {{ .Name }}: {{ call $nullableTSType (call $pgTypeToTSType .Type) .NotNull $table.Type $assumeNonNullableForViews }};
{{- end }}
}
{{- end }}

export interface Database {
{{- range .tables }}
  {{ .Name }}: {{ .Name | title }};
{{- end }}
}`

func (kg *kyselyGenerator) Generate(tables []*Table) (io.Reader, error) {
	tmpl, err := template.New("kysely-schema").Funcs(template.FuncMap{
		"pgTypeToTSType": pgTypeToTSType,
		"nullableTSType": nullableTSType,
		"title":          snakeToPascal,
	}).Parse(tsTemplate)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, map[string]interface{}{
		"tables":                    tables,
		"assumeNonNullableForViews": kg.assumeNonNullableForViews,
		"pgTypeToTSType":            pgTypeToTSType,
		"nullableTSType":            nullableTSType,
		"title":                     snakeToPascal,
	}); err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buffer.Bytes()), nil
}

func snakeToPascal(s string) string {
	return strings.Join(strings.Split(strings.Title(strings.ReplaceAll(s, "_", " ")), " "), "")
}

func pgTypeToTSType(pgType string) string {
	switch pgType {
	case "smallint", "integer", "int4", "bigint", "int8", "serial", "bigserial":
		return "number"
	case "float", "double precision", "numeric", "real", "decimal", "float4", "float8":
		return "number"
	case "boolean":
		return "boolean"
	case "character varying", "text", "varchar", "char", "name":
		return "string"
	case "date", "time", "timestamp", "timestamptz", "timetz":
		return "Date"
	case "json", "jsonb":
		return "unknown"
	default:
		return "unknown"
	}
}

func nullableTSType(tsType string, notNull bool, tableType string, assumeNonNullableForViews bool) string {
	if assumeNonNullableForViews && (tableType == "v" || tableType == "m") {
		return tsType
	}
	if notNull {
		return tsType
	}
	return fmt.Sprintf("%s | null", tsType)
}
