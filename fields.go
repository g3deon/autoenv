package autoenv

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	envTag             = "env"
	jsonTag            = "json"
	fieldPathSeparator = "."
)

type fieldInfo struct {
	field reflect.StructField
	name  string
}

func isFieldIgnored(target, parent, ignore string) bool {
	if parent == ignore {
		return true
	}

	if target == ignore {
		return true
	}

	if strings.HasPrefix(target, ignore+fieldPathSeparator) {
		return true
	}

	return false
}

func joinParent(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + fieldPathSeparator + child
}

func (l *Loader) getStructFields(structType reflect.Type, parent string) []fieldInfo {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return nil
	}

	infos := make([]fieldInfo, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if !l.shouldIncludeField(field) {
			continue
		}

		fieldInfo := l.processStructField(field, parent)
		if fieldInfo != nil {
			infos = append(infos, *fieldInfo...)
		}
	}
	return infos
}

func (l *Loader) processStructField(field reflect.StructField, parent string) *[]fieldInfo {
	fieldType := l.getFieldType(field.Type)
	name := l.resolveFieldName(field)

	isIgnored := l.isIgnoring(name, parent)
	if isIgnored {
		return nil
	}

	if fieldType.Kind() == reflect.Struct {
		return l.processNestedStruct(field, fieldType, name, parent)
	}

	return &[]fieldInfo{{
		field: field,
		name:  name,
	}}
}

func (l *Loader) getFieldType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func (l *Loader) processNestedStruct(field reflect.StructField, fieldType reflect.Type, name, parent string) *[]fieldInfo {
	subFields := l.getStructFields(fieldType, joinParent(parent, name))
	if len(subFields) == 0 {
		return nil
	}

	result := make([]fieldInfo, 0, len(subFields))
	for _, sf := range subFields {
		sf.field.Index = append(field.Index, sf.field.Index...)
		sf.name = name + "_" + sf.name
		result = append(result, sf)
	}
	return &result
}

func (l *Loader) shouldIncludeField(f reflect.StructField) bool {
	if !f.IsExported() {
		if l.isVerbose() {
			l.options.logger.DebugF("%s: excluding (not exported)", f.Name)
		}
		return false
	}

	if l.isOnlyEnvTag() {
		_, ok := f.Tag.Lookup(envTag)
		if l.isVerbose() && !ok {
			l.options.logger.DebugF("%s: excluding (onlyEnvTag without %s)", f.Name, envTag)
		}
		return ok
	}

	return true
}

func (l *Loader) resolveFieldName(f reflect.StructField) string {
	if v, ok := f.Tag.Lookup(envTag); ok {
		return v
	}

	if v, ok := f.Tag.Lookup(jsonTag); ok {
		return v
	}

	return f.Name
}

func (l *Loader) setFieldValue(fv reflect.Value, val string) error {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}
		return l.setFieldValue(fv.Elem(), val)
	}

	switch fv.Kind() {
	case reflect.String:
		fv.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(val)
			if err != nil {
				return err
			}
			fv.SetInt(int64(d))
			return nil
		}
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Struct:
		if fv.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(t))
			return nil
		}
	case reflect.Slice:
		parts := strings.Split(val, ",")
		slice := reflect.MakeSlice(fv.Type(), len(parts), len(parts))
		for i, s := range parts {
			s = strings.TrimSpace(s)
			if err := l.setFieldValue(slice.Index(i), s); err != nil {
				return err
			}
		}
		fv.Set(slice)
	default:
		return &errUnsupportedKind{fv.Kind()}
	}
	return nil
}
