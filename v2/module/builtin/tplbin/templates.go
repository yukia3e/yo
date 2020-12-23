package tplbin

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets5350a404af545b0c50752d1a1e9e1f589fa59a1f = "// Code generated by yo. DO NOT EDIT.\n\n{{if .BuildTag -}}\n// +build {{ .BuildTag}}\n{{- end -}}\n\n// Package {{ .Package }} contains the types.\npackage {{ .Package }}\n\nimport (\n\t\"context\"\n\t\"errors\"\n\t\"fmt\"\n\n\t\"google.golang.org/grpc/codes\"\n\t\"google.golang.org/api/iterator\"\n\t\"cloud.google.com/go/spanner\"\n)\n"
var _Assets55414a05c7f828ef2c56f0399578e07612ca1970 = "{{- range .Indexes }}\n{{- $short := (shortname .Type.Name \"err\" \"sqlstr\" \"db\" \"q\" \"res\" \"YOLog\" .Fields) -}}\n{{- $table := (.Type.Table.TableName) -}}\n\n{{- if not .Index.IsUnique }}\n// Find{{ .FuncName }} retrieves multiple rows from '{{ $table }}' as a slice of {{ .Type.Name }}.\n//\n// Generated from index '{{ .Index.IndexName }}'.\nfunc Find{{ .FuncName }}(ctx context.Context, db YORODB{{ gocustomparamlist .Fields true true }}) ([]*{{ .Type.Name }}, error) {\n{{- else }}\n// Find{{ .FuncName }} retrieves a row from '{{ $table }}' as a {{ .Type.Name }}.\n//\n// If no row is present with the given key, then ReadRow returns an error where\n// spanner.ErrCode(err) is codes.NotFound.\n//\n// Generated from unique index '{{ .Index.IndexName }}'.\nfunc Find{{ .FuncName }}(ctx context.Context, db YORODB{{ gocustomparamlist .Fields true true }}) (*{{ .Type.Name }}, error) {\n{{- end }}\n\t{{- if not .NullableFields }}\n\tconst sqlstr = \"SELECT \" +\n\t\t\"{{ escapedcolnames .Type.Fields }} \" +\n\t\t\"FROM {{ $table }}@{FORCE_INDEX={{ .Index.IndexName }}} \" +\n\t\t\"WHERE {{ colnamesquery .Fields \" AND \" }}\"\n\t{{- else }}\n\tvar sqlstr = \"SELECT \" +\n\t\t\"{{ escapedcolnames .Type.Fields }} \" +\n\t\t\"FROM {{ $table }}@{FORCE_INDEX={{ .Index.IndexName }}} \"\n\n\tconds := make([]string, {{ columncount .Fields }})\n\t{{- range $i, $f := .Fields }}\n\t{{- if $f.Col.NotNull }}\n\t\tconds[{{ $i }}] = \"{{ escapedcolname $f.Col }} = @param{{ $i }}\"\n\t{{- else }}\n\tif {{ nullcheck $f }} {\n\t\tconds[{{ $i }}] = \"{{ escapedcolname $f.Col }} IS NULL\"\n\t} else {\n\t\tconds[{{ $i }}] = \"{{ escapedcolname $f.Col }} = @param{{ $i }}\"\n\t}\n\t{{- end }}\n\t{{- end }}\n\tsqlstr += \"WHERE \" + strings.Join(conds, \" AND \")\n\t{{- end }}\n\n\tstmt := spanner.NewStatement(sqlstr)\n\t{{- range $i, $f := .Fields }}\n\t\t{{- if $f.CustomType }}\n\t\t\tstmt.Params[\"param{{ $i }}\"] = {{ $f.Type }}({{ goparamname $f.Name }})\n\t\t{{- else }}\n\t\t\tstmt.Params[\"param{{ $i }}\"] = {{ goparamname $f.Name }}\n\t\t{{- end }}\n\t{{- end}}\n\n\n\tdecoder := new{{ .Type.Name }}_Decoder({{ .Type.Name }}Columns())\n\n\t// run query\n\tYOLog(ctx, sqlstr{{ goparamlist .Fields true false }})\n{{- if .Index.IsUnique }}\n\titer := db.Query(ctx, stmt)\n\tdefer iter.Stop()\n\n\trow, err := iter.Next()\n\tif err != nil {\n\t\tif err == iterator.Done {\n\t\t\treturn nil, newErrorWithCode(codes.NotFound, \"Find{{ .FuncName }}\", \"{{ $table }}\", err)\n\t\t}\n\t\treturn nil, newError(\"Find{{ .FuncName }}\", \"{{ $table }}\", err)\n\t}\n\n\t{{ $short }}, err := decoder(row)\n\tif err != nil {\n\t\treturn nil, newErrorWithCode(codes.Internal, \"Find{{ .FuncName }}\", \"{{ $table }}\", err)\n\t}\n\n\treturn {{ $short }}, nil\n{{- else }}\n\titer := db.Query(ctx, stmt)\n\tdefer iter.Stop()\n\n\t// load results\n\tres := []*{{ .Type.Name }}{}\n\tfor {\n\t\trow, err := iter.Next()\n\t\tif err != nil {\n\t\t\tif err == iterator.Done {\n\t\t\t\tbreak\n\t\t\t}\n\t\t\treturn nil, newError(\"Find{{ .FuncName }}\", \"{{ $table }}\", err)\n\t\t}\n\n\t\t{{ $short }}, err := decoder(row)\n        if err != nil {\n            return nil, newErrorWithCode(codes.Internal, \"Find{{ .FuncName }}\", \"{{ $table }}\", err)\n        }\n\n\t\tres = append(res, {{ $short }})\n\t}\n\n\treturn res, nil\n{{- end }}\n}\n\n\n// Read{{ .FuncName }} retrieves multiples rows from '{{ $table }}' by KeySet as a slice.\n//\n// This does not retrives all columns of '{{ $table }}' because an index has only columns\n// used for primary key, index key and storing columns. If you need more columns, add storing\n// columns or Read by primary key or Query with join.\n//\n// Generated from unique index '{{ .Index.IndexName }}'.\nfunc Read{{ .FuncName }}(ctx context.Context, db YORODB, keys spanner.KeySet) ([]*{{ .Type.Name }}, error) {\n\tvar res []*{{ .Type.Name }}\n    columns := []string{\n{{- range .Type.PrimaryKeyFields }}\n\t\t\"{{ colname .Col }}\",\n{{- end }}\n{{- range .Fields }}\n\t\t\"{{ colname .Col }}\",\n{{- end }}\n{{- range .StoringFields }}\n\t\t\"{{ colname .Col }}\",\n{{- end }}\n}\n\n\tdecoder := new{{ .Type.Name }}_Decoder(columns)\n\n\trows := db.ReadUsingIndex(ctx, \"{{ $table }}\", \"{{ .Index.IndexName }}\", keys, columns)\n\terr := rows.Do(func(row *spanner.Row) error {\n\t\t{{ $short }}, err := decoder(row)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tres = append(res, {{ $short }})\n\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\treturn nil, newErrorWithCode(codes.Internal, \"Read{{ .FuncName }}\", \"{{ $table }}\", err)\n\t}\n\n    return res, nil\n}\n{{- end -}}\n"
var _Assetsa6d3263a736fc023cdf294784a762b9a5817e9ef = "{{- $short := (shortname .Name \"err\" \"res\" \"sqlstr\" \"db\" \"YOLog\") -}}\n{{- $table := (.Table.TableName) -}}\n\n// {{ .Name }} represents a row from '{{ $table }}'.\ntype {{ .Name }} struct {\n{{- range .Fields }}\n{{- if eq (.Col.DataType) (.Col.ColumnName) }}\n\t{{ .Name }} string `spanner:\"{{ .Col.ColumnName }}\" json:\"{{ .Col.ColumnName }}\"` // {{ .Col.ColumnName }} enum\n{{- else if .CustomType }}\n\t{{ .Name }} {{ retype .CustomType }} `spanner:\"{{ .Col.ColumnName }}\" json:\"{{ .Col.ColumnName }}\"` // {{ .Col.ColumnName }}\n{{- else }}\n\t{{ .Name }} {{ .Type }} `spanner:\"{{ .Col.ColumnName }}\" json:\"{{ .Col.ColumnName }}\"` // {{ .Col.ColumnName }}\n{{- end }}\n{{- end }}\n}\n\nfunc {{ .Name }}PrimaryKeys() []string {\n     return []string{\n{{- range .PrimaryKeyFields }}\n\t\t\"{{ colname .Col }}\",\n{{- end }}\n\t}\n}\n\nfunc {{ .Name }}Columns() []string {\n\treturn []string{\n{{- range .Fields }}\n\t\t\"{{ colname .Col }}\",\n{{- end }}\n\t}\n}\n\nfunc ({{ $short }} *{{ .Name }}) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {\n\tret := make([]interface{}, 0, len(cols))\n\tfor _, col := range cols {\n\t\tif val, ok := customPtrs[col]; ok {\n\t\t\tret = append(ret, val)\n\t\t\tcontinue\n\t\t}\n\n\t\tswitch col {\n{{- range .Fields }}\n\t\tcase \"{{ colname .Col }}\":\n\t\t\tret = append(ret, &{{ $short }}.{{ .Name }})\n{{- end }}\n\t\tdefault:\n\t\t\treturn nil, fmt.Errorf(\"unknown column: %s\", col)\n\t\t}\n\t}\n\treturn ret, nil\n}\n\nfunc ({{ $short }} *{{ .Name }}) columnsToValues(cols []string) ([]interface{}, error) {\n\tret := make([]interface{}, 0, len(cols))\n\tfor _, col := range cols {\n\t\tswitch col {\n{{- range .Fields }}\n\t\tcase \"{{ colname .Col }}\":\n\t\t\t{{- if .CustomType }}\n\t\t\tret = append(ret, {{ .Type }}({{ $short }}.{{ .Name }}))\n\t\t\t{{- else }}\n\t\t\tret = append(ret, {{ $short }}.{{ .Name }})\n\t\t\t{{- end }}\n{{- end }}\n\t\tdefault:\n\t\t\treturn nil, fmt.Errorf(\"unknown column: %s\", col)\n\t\t}\n\t}\n\n\treturn ret, nil\n}\n\n// new{{ .Name }}_Decoder returns a decoder which reads a row from *spanner.Row\n// into {{ .Name }}. The decoder is not goroutine-safe. Don't use it concurrently.\nfunc new{{ .Name }}_Decoder(cols []string) func(*spanner.Row) (*{{ .Name }}, error) {\n\t{{- range .Fields }}\n\t\t{{- if .CustomType }}\n\t\t\tvar {{ customtypeparam .Name }} {{ .Type }}\n\t\t{{- end }}\n\t{{- end }}\n\tcustomPtrs := map[string]interface{}{\n\t\t{{- range .Fields }}\n\t\t\t{{- if .CustomType }}\n\t\t\t\t\"{{ colname .Col }}\": &{{ customtypeparam .Name }},\n\t\t\t{{- end }}\n\t{{- end }}\n\t}\n\n\treturn func(row *spanner.Row) (*{{ .Name }}, error) {\n        var {{ $short }} {{ .Name }}\n        ptrs, err := {{ $short }}.columnsToPtrs(cols, customPtrs)\n        if err != nil {\n            return nil, err\n        }\n\n        if err := row.Columns(ptrs...); err != nil {\n            return nil, err\n        }\n        {{- range .Fields }}\n            {{- if .CustomType }}\n                {{ $short }}.{{ .Name }} = {{ retype .CustomType }}({{ customtypeparam .Name }})\n            {{- end }}\n        {{- end }}\n\n\n\t\treturn &{{ $short }}, nil\n\t}\n}\n"
var _Assets185d1c2aed0961d2cf30ea0fd83b17ca89757c11 = "// YODB is the common interface for database operations.\ntype YODB interface {\n\tYORODB\n}\n\n// YORODB is the common interface for database operations.\ntype YORODB interface {\n\tReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)\n\tRead(ctx context.Context, table string, keys spanner.KeySet, columns []string) *spanner.RowIterator\n\tReadUsingIndex(ctx context.Context, table, index string, keys spanner.KeySet, columns []string) (ri *spanner.RowIterator)\n\tQuery(ctx context.Context, statement spanner.Statement) *spanner.RowIterator\n}\n\n// YOLog provides the log func used by generated queries.\nvar YOLog = func(context.Context, string, ...interface{}) { }\n\nfunc newError(method, table string, err error) error {\n\tcode := spanner.ErrCode(err)\n\treturn newErrorWithCode(code, method, table, err)\n}\n\nfunc newErrorWithCode(code codes.Code, method, table string, err error) error {\n\treturn &yoError{\n\t\tmethod: method,\n\t\ttable:  table,\n\t\terr:    err,\n\t\tcode:   code,\n\t}\n}\n\ntype yoError struct {\n\terr    error\n\tmethod string\n\ttable  string\n\tcode   codes.Code\n}\n\nfunc (e yoError) Error() string {\n\treturn fmt.Sprintf(\"yo error in %s(%s): %v\", e.method, e.table, e.err)\n}\n\nfunc (e yoError) Unwrap() error {\n\treturn e.err\n}\n\nfunc (e yoError) DBTableName() string {\n\treturn e.table\n}\n\n// GRPCStatus implements a conversion to a gRPC status using `status.Convert(error)`.\n// If the error is originated from the Spanner library, this returns a gRPC status of\n// the original error. It may contain details of the status such as RetryInfo.\nfunc (e yoError) GRPCStatus() *status.Status {\n\tvar se *spanner.Error\n\tif errors.As(e.err, &se) {\n\t\treturn status.Convert(se.Unwrap())\n\t}\n\n\treturn status.New(e.code, e.Error())\n}\n\nfunc (e yoError) Timeout() bool { return e.code == codes.DeadlineExceeded }\nfunc (e yoError) Temporary() bool { return e.code == codes.DeadlineExceeded }\nfunc (e yoError) NotFound() bool { return e.code == codes.NotFound }\n"
var _Assets9831b3c5674e85d920191e8a99a167a7c9080158 = "{{- $short := (shortname .Name \"err\" \"res\" \"sqlstr\" \"db\" \"YOLog\") -}}\n{{- $table := (.Table.TableName) -}}\n\n// Insert returns a Mutation to insert a row into a table. If the row already\n// exists, the write or transaction fails.\nfunc ({{ $short }} *{{ .Name }}) Insert(ctx context.Context) *spanner.Mutation {\n\treturn spanner.Insert(\"{{ $table }}\", {{ .Name }}Columns(), []interface{}{\n\t\t{{ fieldnames .Fields $short }},\n\t})\n}\n\n{{ if ne (fieldnames .Fields $short .PrimaryKeyFields) \"\" }}\n// Update returns a Mutation to update a row in a table. If the row does not\n// already exist, the write or transaction fails.\nfunc ({{ $short }} *{{ .Name }}) Update(ctx context.Context) *spanner.Mutation {\n\treturn spanner.Update(\"{{ $table }}\", {{ .Name }}Columns(), []interface{}{\n\t\t{{ fieldnames .Fields $short }},\n\t})\n}\n\n// InsertOrUpdate returns a Mutation to insert a row into a table. If the row\n// already exists, it updates it instead. Any column values not explicitly\n// written are preserved.\nfunc ({{ $short }} *{{ .Name }}) InsertOrUpdate(ctx context.Context) *spanner.Mutation {\n\treturn spanner.InsertOrUpdate(\"{{ $table }}\", {{ .Name }}Columns(), []interface{}{\n\t\t{{ fieldnames .Fields $short }},\n\t})\n}\n\n// UpdateColumns returns a Mutation to update specified columns of a row in a table.\nfunc ({{ $short }} *{{ .Name }}) UpdateColumns(ctx context.Context, cols ...string) (*spanner.Mutation, error) {\n\t// add primary keys to columns to update by primary keys\n\tcolsWithPKeys := append(cols, {{ .Name }}PrimaryKeys()...)\n\n\tvalues, err := {{ $short }}.columnsToValues(colsWithPKeys)\n\tif err != nil {\n\t\treturn nil, newErrorWithCode(codes.InvalidArgument, \"{{ .Name }}.UpdateColumns\", \"{{ $table }}\", err)\n\t}\n\n\treturn spanner.Update(\"{{ $table }}\", colsWithPKeys, values), nil\n}\n\n// Find{{ .Name }} gets a {{ .Name }} by primary key\nfunc Find{{ .Name }}(ctx context.Context, db YORODB{{ gocustomparamlist .PrimaryKeyFields true true }}) (*{{ .Name }}, error) {\n\tkey := spanner.Key{ {{ gocustomparamlist .PrimaryKeyFields false false }} }\n\trow, err := db.ReadRow(ctx, \"{{ $table }}\", key, {{ .Name }}Columns())\n\tif err != nil {\n\t\treturn nil, newError(\"Find{{ .Name }}\", \"{{ $table }}\", err)\n\t}\n\n\tdecoder := new{{ .Name }}_Decoder({{ .Name}}Columns())\n\t{{ $short }}, err := decoder(row)\n\tif err != nil {\n\t\treturn nil, newErrorWithCode(codes.Internal, \"Find{{ .Name }}\", \"{{ $table }}\", err)\n\t}\n\n\treturn {{ $short }}, nil\n}\n\n// Read{{ .Name }} retrieves multiples rows from {{ .Name }} by KeySet as a slice.\nfunc Read{{ .Name }}(ctx context.Context, db YORODB, keys spanner.KeySet) ([]*{{ .Name }}, error) {\n\tvar res []*{{ .Name }}\n\n\tdecoder := new{{ .Name }}_Decoder({{ .Name}}Columns())\n\n\trows := db.Read(ctx, \"{{ $table }}\", keys, {{ .Name }}Columns())\n\terr := rows.Do(func(row *spanner.Row) error {\n\t\t{{ $short }}, err := decoder(row)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tres = append(res, {{ $short }})\n\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\treturn nil, newErrorWithCode(codes.Internal, \"Read{{ .Name }}\", \"{{ $table }}\", err)\n\t}\n\n\treturn res, nil\n}\n{{ end }}\n\n// Delete deletes the {{ .Name }} from the database.\nfunc ({{ $short }} *{{ .Name }}) Delete(ctx context.Context) *spanner.Mutation {\n\tvalues, _ := {{ $short }}.columnsToValues({{ .Name }}PrimaryKeys())\n\treturn spanner.Delete(\"{{ $table }}\", spanner.Key(values))\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"index.go.tpl": &assets.File{
		Path:     "index.go.tpl",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1608701581, 1608701581333731623),
		Data:     []byte(_Assets55414a05c7f828ef2c56f0399578e07612ca1970),
	}, "type.go.tpl": &assets.File{
		Path:     "type.go.tpl",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1608700373, 1608700373784628030),
		Data:     []byte(_Assetsa6d3263a736fc023cdf294784a762b9a5817e9ef),
	}, "header.go.tpl": &assets.File{
		Path:     "header.go.tpl",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1608617403, 1608617403644983290),
		Data:     []byte(_Assets5350a404af545b0c50752d1a1e9e1f589fa59a1f),
	}, "operation.go.tpl": &assets.File{
		Path:     "operation.go.tpl",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1608700367, 1608700367740519711),
		Data:     []byte(_Assets9831b3c5674e85d920191e8a99a167a7c9080158),
	}, "yo_db.go.tpl": &assets.File{
		Path:     "yo_db.go.tpl",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1608610929, 1608610929145130051),
		Data:     []byte(_Assets185d1c2aed0961d2cf30ea0fd83b17ca89757c11),
	}}, "")
