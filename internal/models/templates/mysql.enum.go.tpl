{{- $type := .Name -}}
{{- $short := (shortname $type "enumVal" "text" "buf" "ok" "src") -}}
{{- $reverseNames := .ReverseConstNames -}}
// {{ $type }} is the '{{ .Enum.EnumName }}' enum type from schema '{{ .Schema  }}'.
type {{ $type }} uint16

const (
{{- range .Values }}
	// {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }} is the '{{ .Val.EnumValue }}' {{ $type }}.
	{{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }} = {{ $type }}({{ .Val.ConstValue }})
{{ end -}}
)

// {{ $type }}FromString converts a string to a {{ $type }}.
func {{ $type }}FromString(str string) ({{ $type }}, error) {
	var enumVal {{ $type }}

	err := enumVal.UnmarshalText([]byte(str))
	return enumVal, err
}

// String returns the string value of the {{ $type }}.
func ({{ $short }} {{ $type }}) String() string {
	var enumVal string

	switch {{ $short }} {
{{- range .Values }}
	case {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }}:
		enumVal = "{{ .Val.EnumValue }}"
{{ end -}}
	}

	return enumVal
}

// MarshalText marshals {{ $type }} into text.
func ({{ $short }} {{ $type }}) MarshalText() ([]byte, error) {
	return []byte({{ $short }}.String()), nil
}

// UnmarshalText unmarshals {{ $type }} from text.
func ({{ $short }} *{{ $type }}) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
{{- range .Values }}
	case "{{ .Val.EnumValue }}":
		*{{ $short }} = {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }}
{{- end }}
	default:
		return errors.New("invalid {{ $type }}")
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func ({{ $short }} *{{ $type }}) UnmarshalJSON(data []byte) error {
	if data[0] != '"' { // string
		return &json.UnmarshalTypeError{Value: "{{ $type }}", Type: reflect.TypeOf({{ $short }})}
	}

	if err := {{ $short }}.UnmarshalText(data[1 : len(data)-1]); err != nil {
		return &json.UnmarshalTypeError{Value: "{{ $type }}", Type: reflect.TypeOf({{ $short }})}
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for {{ $type }}.
func ({{ $short }} {{ $type }}) Value() (driver.Value, error) {
	return {{ $short }}.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for {{ $type }}.
func ({{ $short }} *{{ $type }}) Scan(src interface{}) error {
	var buf []byte
	switch v := src.(type) {
	case []byte:
		buf = v
	case string:
		buf = []byte(v)
	default:
		return errors.New("invalid {{ $type }}")
	}

	return {{ $short }}.UnmarshalText(buf)
}

