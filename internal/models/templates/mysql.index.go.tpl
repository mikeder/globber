{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" .Fields) -}}
{{- $table := (schema .Schema .Type.Table.TableName) -}}
// {{ .FuncName }} retrieves a row from '{{ $table }}' as a {{ .Type.Name }}.
//
// Generated from index '{{ .Index.IndexName }}'.
func {{ .FuncName }}(ctx context.Context, db XODB{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}[]{{ end }}*{{ .Type.Name }}, error) {
	var err error

	// sql query
	const sqlstr = "SELECT " +
		"{{ colnames .Type.Fields }} " +
		"FROM {{ $table }} " +
		"WHERE {{ colnamesquery .Fields " AND " }}"

	// run query
{{- if .Index.IsUnique }}
	{{ $short }} := {{ .Type.Name }}{
	{{- if .Type.PrimaryKey }}
		_exists: true,
	{{ end -}}
	}

	err = db.QueryRowContext(ctx, sqlstr{{ goparamlist .Fields true false }}).Scan({{ fieldnames .Type.Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}

	return &{{ $short }}, nil
{{- else }}
	q, err := db.QueryContext(ctx, sqlstr{{ goparamlist .Fields true false }})
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*{{ .Type.Name }}{}
	for q.Next() {
		{{ $short }} := {{ .Type.Name }}{
		{{- if .Type.PrimaryKey }}
			_exists: true,
		{{ end -}}
		}

		// scan
		err = q.Scan({{ fieldnames .Type.Fields (print "&" $short) }})
		if err != nil {
			return nil, err
		}

		res = append(res, &{{ $short }})
	}
	if err := q.Err(); err != nil {
		return nil, err
	}

	return res, nil
{{- end }}
}

