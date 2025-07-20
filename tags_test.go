package autoenv

import (
	"reflect"
	"testing"
)

func TestGetStructTags(t *testing.T) {
	type Inner struct {
		InnerField string `env:"INNER_FIELD"`
	}

	type Outer struct {
		OuterField string `env:"OUTER_FIELD"`
		Inner      Inner
	}

	type JsonStruct struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Email    string `json:"email,omitempty"`
		Password string `json:"-"`
		Hidden   string
	}

	type MixedTags struct {
		Field1 string `env:"FIELD_1" json:"field1"`
		Field2 string `env:"FIELD_2" json:"field2,omitempty"`
		Field3 string `json:"field3"`
		Field4 string `env:"FIELD_4"`
		Field5 string
	}

	type JsonWithEmbedded struct {
		OuterField string `json:"outer_field"`
		JsonStruct
	}

	tests := []struct {
		input    any
		tag      string
		expected []string
	}{
		{Outer{OuterField: "value1", Inner: Inner{InnerField: "value2"}}, "env", []string{"OUTER_FIELD"}},
		{Inner{InnerField: "value2"}, "env", []string{"INNER_FIELD"}},
		{Outer{}, "env", []string{"OUTER_FIELD"}},

		{JsonStruct{}, "json", []string{"name", "age", "email,omitempty"}},
		{JsonStruct{Name: "John", Age: 30}, "json", []string{"name", "age", "email,omitempty"}},

		{MixedTags{}, "env", []string{"FIELD_1", "FIELD_2", "FIELD_4"}},
		{MixedTags{}, "json", []string{"field1", "field2,omitempty", "field3"}},

		{JsonWithEmbedded{}, "json", []string{"outer_field", "name", "age", "email,omitempty"}},

		{struct{}{}, "json", []string{}},
		{struct{}{}, "env", []string{}},
	}

	for _, test := range tests {
		result := getStructTags(test.input, test.tag)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("getStructTags(%v, %q) = %v; want %v", test.input, test.tag, result, test.expected)
		}
	}
}
