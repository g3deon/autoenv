package autoenv

import "reflect"

func getStructTags(i any, tag string) []string {
	tags := make([]string, 0)

	reflectType := reflect.TypeOf(i)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	if reflectType.Kind() != reflect.Struct {
		return tags
	}

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			tags = append(tags, getStructTags(reflect.New(field.Type).Elem().Interface(), tag)...)
			continue
		}

		val := field.Tag.Get(tag)
		if val == "" || val == "-" {
			continue
		}

		tags = append(tags, val)
	}

	return tags
}
