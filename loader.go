package autoenv

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
)

type Loader struct {
	options options
}

func NewLoader(options ...Option) *Loader {
	opts := newOptions()
	opts.apply(options...)

	return &Loader{options: opts}
}

func (l *Loader) Load(i any) error {
	if l == nil {
		return nil
	}

	if i == nil {
		return ErrNilInput
	}

	if l.isVerbose() {
		l.options.logger.DebugF("loading struct %T", i)
	}

	l.loadEnvFiles()

	t := reflect.TypeOf(i)
	fields := l.getStructFields(t, "")
	return l.mapEnvValues(reflect.ValueOf(i), fields)
}

func (l *Loader) isVerbose() bool {
	if l == nil {
		return false
	}

	return l.options.verbose
}

func (l *Loader) isOnlyEnvTag() bool {
	if l == nil {
		return false
	}

	return l.options.onlyEnvTag
}

func (l *Loader) isIgnoring(field, parent string) bool {
	if l == nil {
		return false
	}

	target := joinParent(parent, field)
	target = strings.ToLower(target)
	parent = strings.ToLower(parent)

	return slices.ContainsFunc(l.options.ignores, func(ignore string) bool {
		return isFieldIgnored(target, parent, strings.ToLower(ignore))
	})
}

func (l *Loader) loadEnvFiles() {
	if !l.options.withFiles {
		return
	}

	for _, fileName := range l.options.filesPaths {
		if err := l.loadEnvFile(fileName); err != nil {
			l.options.logger.ErrorF("failed to load file: %s", err)
			return
		}

		if l.isVerbose() {
			l.options.logger.DebugF("loaded file: %s", fileName)
		}
	}
}

func (l *Loader) mapEnvValues(target reflect.Value, fields []fieldInfo) error {
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	for _, fi := range fields {
		key := l.getEnvKey(fi.name)
		if key == "" {
			continue
		}

		val := os.Getenv(key)
		if val == "" {
			continue
		}

		fv := target.FieldByIndex(fi.field.Index)
		if !fv.CanSet() {
			continue
		}

		if err := l.setFieldValue(fv, val); err != nil {
			return err
		}

		if l.isVerbose() {
			l.options.logger.DebugF("loaded %s as %s (%s)", fi.name, key, fi.field.Type.String())
		}
	}

	return nil
}

func (l *Loader) getEnvKey(name string) string {
	if l.options.prefix == "" {
		return strings.ToUpper(toSnakeCase(name))
	}

	return strings.ToUpper(fmt.Sprintf("%s_%s", l.options.prefix, toSnakeCase(name)))
}
