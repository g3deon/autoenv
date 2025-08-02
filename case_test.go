package autoenv

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "empty",
			str:  "",
			want: "",
		},
		{
			name: "single lowercase",
			str:  "simple",
			want: "simple",
		},
		{
			name: "all uppercase",
			str:  "HTML",
			want: "html",
		},
		{
			name: "camelCase",
			str:  "simpleCamelCase",
			want: "simple_camel_case",
		},
		{
			name: "PascalCase",
			str:  "SimpleCamelCase",
			want: "simple_camel_case",
		},
		{
			name: "leading uppercase",
			str:  "Camel",
			want: "camel",
		},
		{
			name: "single letter start",
			str:  "aTest",
			want: "a_test",
		},
		{
			name: "HTTP prefix",
			str:  "HTTPServer",
			want: "http_server",
		},
		{
			name: "mixed acronym",
			str:  "getXMLParser",
			want: "get_xml_parser",
		},
		{
			name: "XMLHttpRequest",
			str:  "XMLHttpRequest",
			want: "xml_http_request",
		},
		{
			name: "trailing numbers",
			str:  "version2",
			want: "version2",
		},
		{
			name: "numbers in middle",
			str:  "test123Case",
			want: "test123_case",
		},
		{
			name: "number acronym",
			str:  "IPv4Config",
			want: "ipv4_config",
		},
		{
			name: "already snake",
			str:  "already_snake",
			want: "already_snake",
		},
		{
			name: "mixed underscores",
			str:  "mixed_SNakeCase",
			want: "mixed_sn_ake_case",
		},
		{
			name: "single char",
			str:  "X",
			want: "x",
		},
		{
			name: "multiple capitols",
			str:  "ABC",
			want: "abc",
		},
		{
			name: "ends with uppercase",
			str:  "endsWithX",
			want: "ends_with_x",
		},
		{
			name: "starts with underscore",
			str:  "_privateField",
			want: "_private_field",
		},
		{
			name: "trailing underscore",
			str:  "field_",
			want: "field_",
		},
		{
			name: "special characters",
			str:  "field@name",
			want: "field@name",
		},
		{
			name: "numeric prefix",
			str:  "123Field",
			want: "123_field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSnakeCase(tt.str)
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q; want %q", tt.str, got, tt.want)
			}
		})
	}
}

func TestGetNextChars(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		index int
		want1 byte
		want2 byte
	}{
		{
			name:  "empty string",
			str:   "",
			index: 0,
			want1: 0,
			want2: 0,
		},
		{
			name:  "single character string",
			str:   "a",
			index: 0,
			want1: 0,
			want2: 0,
		},
		{
			name:  "normal case",
			str:   "abc",
			index: 1,
			want1: 'c',
			want2: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := getNextChars(tt.str, tt.index)
			if got1 != tt.want1 || got2 != tt.want2 {
				t.Errorf("getNextChars(%q, %d) = (%q, %q); want (%q, %q)", tt.str, tt.index, got1, got2, tt.want1, tt.want2)
			}
		})
	}
}

func TestIsUppercase(t *testing.T) {
	tests := []struct {
		name string
		char byte
		want bool
	}{
		{
			name: "uppercase letter",
			char: 'A',
			want: true,
		},
		{
			name: "lowercase letter",
			char: 'a',
			want: false,
		},
		{
			name: "digit",
			char: '1',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUppercase(tt.char)
			if got != tt.want {
				t.Errorf("isUppercase(%q) = %v; want %v", tt.char, got, tt.want)
			}
		})
	}
}
