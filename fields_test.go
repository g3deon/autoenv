package autoenv

import (
	"reflect"
	"testing"
	"time"
)

func TestIsFieldIgnored(t *testing.T) {
	tests := []struct {
		name   string
		target string
		parent string
		ignore string
		want   bool
	}{
		{
			name:   "Ignore matches parent",
			target: "field.child",
			parent: "parent",
			ignore: "parent",
			want:   true,
		},
		{
			name:   "Ignore matches target",
			target: "field.child",
			parent: "parent",
			ignore: "field.child",
			want:   true,
		},
		{
			name:   "Ignore matches target prefix with fieldPathSeparator",
			target: "field.child.subfield",
			parent: "parent",
			ignore: "field.child",
			want:   true,
		},
		{
			name:   "Ignore doesn't match anything",
			target: "field.child",
			parent: "parent",
			ignore: "unrelated",
			want:   false,
		},
		{
			name:   "Empty ignore string",
			target: "field.child",
			parent: "parent",
			ignore: "",
			want:   false,
		},
		{
			name:   "Empty parent string",
			target: "field.child",
			parent: "",
			ignore: "field.child",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFieldIgnored(tt.target, tt.parent, tt.ignore); got != tt.want {
				t.Errorf("isFieldIgnored() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinParent(t *testing.T) {
	tests := []struct {
		name   string
		parent string
		child  string
		want   string
	}{
		{
			name:   "Both parent and child are non-empty",
			parent: "parent",
			child:  "child",
			want:   "parent.child",
		},
		{
			name:   "Parent is empty",
			parent: "",
			child:  "child",
			want:   "child",
		},
		{
			name:   "Child is empty",
			parent: "parent",
			child:  "",
			want:   "parent.",
		},
		{
			name:   "Both parent and child are empty",
			parent: "",
			child:  "",
			want:   "",
		},
		{
			name:   "Parent contains a dot",
			parent: "parent.part",
			child:  "child",
			want:   "parent.part.child",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinParent(tt.parent, tt.child); got != tt.want {
				t.Errorf("joinParent(%q, %q) = %v, want %v", tt.parent, tt.child, got, tt.want)
			}
		})
	}
}

func TestGetStructFields(t *testing.T) {
	tests := []struct {
		name       string
		structType any
		parent     string
		want       []fieldInfo
	}{
		{
			name:       "empty struct",
			structType: struct{}{},
			parent:     "",
			want:       nil,
		},
		{
			name: "struct with exported fields",
			structType: struct {
				Field1 string
				Field2 int
			}{},
			parent: "",
			want: []fieldInfo{
				{name: "Field1"},
				{name: "Field2"},
			},
		},
		{
			name: "struct with unexported fields",
			structType: struct {
				Field1 string
				field2 int
			}{},
			parent: "",
			want: []fieldInfo{
				{name: "Field1"},
			},
		},
		{
			name: "struct with env tag",
			structType: struct {
				Field1 string `env:"custom_field"`
			}{},
			parent: "",
			want: []fieldInfo{
				{name: "custom_field"},
			},
		},
		{
			name: "struct with nested struct",
			structType: struct {
				ParentField struct {
					ChildField string
				}
			}{},
			parent: "",
			want: []fieldInfo{
				{name: "ParentField_ChildField"},
			},
		},
		{
			name: "struct with pointer fields",
			structType: struct {
				Field1 *string
			}{},
			parent: "",
			want: []fieldInfo{
				{name: "Field1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewLoader()
			got := loader.getStructFields(reflect.TypeOf(tt.structType), tt.parent)

			gotNames := make([]string, len(got))
			for i, f := range got {
				gotNames[i] = f.name
			}
			wantNames := make([]string, len(tt.want))
			for i, f := range tt.want {
				wantNames[i] = f.name
			}

			if !reflect.DeepEqual(gotNames, wantNames) {
				t.Errorf("getStructFields() = %v, want %v", gotNames, wantNames)
			}
		})

	}
}

func TestResolveFieldName(t *testing.T) {
	tests := []struct {
		name  string
		field reflect.StructField
		want  string
	}{
		{
			name: "Field with env tag has priority",
			field: reflect.StructField{
				Name: "Field1",
				Tag:  `env:"custom_env" json:"custom_json"`,
			},
			want: "custom_env",
		},
		{
			name: "Field with json tag is chosen if env tag is absent",
			field: reflect.StructField{
				Name: "Field2",
				Tag:  `json:"custom_json"`,
			},
			want: "custom_json",
		},
		{
			name: "Defaults to field name if no tags are present",
			field: reflect.StructField{
				Name: "DefaultField",
				Tag:  ``,
			},
			want: "DefaultField",
		},
		{
			name: "Ignores unrelated tags, defaults to field name",
			field: reflect.StructField{
				Name: "UnrelatedField",
				Tag:  `xml:"custom_xml"`,
			},
			want: "UnrelatedField",
		},
		{
			name: "Empty field name, no tags",
			field: reflect.StructField{
				Name: "",
				Tag:  ``,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewLoader()
			if got := loader.resolveFieldName(tt.field); got != tt.want {
				t.Errorf("resolveFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldIncludeField(t *testing.T) {
	tests := []struct {
		name           string
		field          reflect.StructField
		options        options
		expectedResult bool
	}{
		{
			name: "field with env tag and onlyEnvTag enabled",
			field: reflect.StructField{
				Name: "EnvField",
				Tag:  `env:"ENV_FIELD"`,
			},
			options:        options{onlyEnvTag: true},
			expectedResult: true,
		},
		{
			name: "field with no env tag and onlyEnvTag enabled",
			field: reflect.StructField{
				Name: "NoEnvField",
				Tag:  `json:"no_env_field"`,
			},
			options:        options{onlyEnvTag: true},
			expectedResult: false,
		},
		{
			name: "unexported field with onlyEnvTag disabled",
			field: reflect.StructField{
				Name:    "unexportedField",
				Tag:     "",
				PkgPath: "go.g3deon.com/autoenv",
			},
			options:        options{onlyEnvTag: false},
			expectedResult: false,
		},
		{
			name: "field with no env tag and onlyEnvTag disabled",
			field: reflect.StructField{
				Name: "RegularField",
				Tag:  "",
			},
			options:        options{onlyEnvTag: false},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := &Loader{options: tt.options}
			result := loader.shouldIncludeField(tt.field)

			if result != tt.expectedResult {
				t.Errorf("shouldIncludeField() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	loader := NewLoader()
	now := time.Now().Truncate(time.Second).UTC()

	tests := []struct {
		name     string
		field    reflect.Value
		value    string
		expected any
	}{
		{
			name:     "time.Duration field",
			field:    reflect.ValueOf(new(time.Duration)).Elem(),
			value:    "1h30m",
			expected: time.Hour + 30*time.Minute,
		},
		{
			name:     "time.Time field",
			field:    reflect.ValueOf(new(time.Time)).Elem(),
			value:    now.Format(time.RFC3339),
			expected: now,
		},
		{
			name:     "slice of time.Duration",
			field:    reflect.ValueOf([]time.Duration{}),
			value:    "1s,2m,3h",
			expected: []time.Duration{time.Second, 2 * time.Minute, 3 * time.Hour},
		},
		{
			name:     "slice of time.Time",
			field:    reflect.ValueOf([]time.Time{}),
			value:    now.Format(time.RFC3339) + "," + now.Add(time.Hour).Format(time.RFC3339),
			expected: []time.Time{now, now.Add(time.Hour)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := tt.field
			if fv.Kind() == reflect.Slice && fv.Len() == 0 {
				fv = reflect.New(tt.field.Type()).Elem()
			}
			if err := loader.setFieldValue(fv, tt.value); err != nil {
				t.Fatalf("failed to set field: %v", err)
			}

			switch expected := tt.expected.(type) {
			case time.Duration:
				got := fv.Interface().(time.Duration)
				if got != expected {
					t.Errorf("expected %v, got %v", expected, got)
				}
			case time.Time:
				got := fv.Interface().(time.Time)
				if !got.Equal(expected) {
					t.Errorf("expected %v, got %v", expected, got)
				}
			case []time.Duration:
				got := fv.Interface().([]time.Duration)
				for i := range expected {
					if got[i] != expected[i] {
						t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
					}
				}
			case []time.Time:
				got := fv.Interface().([]time.Time)
				for i := range expected {
					if !got[i].Equal(expected[i]) {
						t.Errorf("index %d: expected %v, got %v", i, expected[i], got[i])
					}
				}
			}
		})
	}
}
