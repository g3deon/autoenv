package autoenv

import (
	"os"
	"reflect"
)

func assignValuesToStruct(s reflect.Value, values ...reflect.Value) {
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	maxFields := s.NumField()
	if len(values) < maxFields {
		maxFields = len(values)
	}

	for i := 0; i < maxFields; i++ {
		field := s.Field(i)
		if !field.CanSet() {
			continue
		}

		value := values[i]
		if !value.IsValid() || !value.Type().AssignableTo(field.Type()) {
			continue
		}

		field.Set(value)
	}
}

func getTagsValuesFromEnv(tags []string) []reflect.Value {
	values := make([]reflect.Value, len(tags))

	for i, tag := range tags {
		value := getEnvValue(tag)
		if value == "" {
			values[i] = reflect.Zero(reflect.TypeOf(""))
			continue
		}

		values[i] = reflect.ValueOf(value)
	}

	return values
}

func getEnvValue(tag string) string {
	return os.Getenv(toSnakeCase(tag))
}
