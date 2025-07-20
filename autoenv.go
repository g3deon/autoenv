package autoenv

import (
	"reflect"
)

const (
	defaultTag = "json"
	envTag     = "env"
)

type Loader struct {
	options *options
}

func New(i any, opts ...Option) (*Loader, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	loader := &Loader{
		options: options,
	}

	if options.withEnvFile {
		if err := loader.loadEnvFile(); err != nil {
			return nil, err
		}
	}

	return loader, nil
}

func MustNew(i any, opts ...Option) *Loader {
	loader, err := New(i, opts...)
	if err != nil {
		panic(err)
	}

	return loader
}

func (l *Loader) Load(i any) error {
	if i == nil {
		return ErrNilInput
	}

	structValue := reflect.ValueOf(i)
	if structValue.Kind() != reflect.Ptr || structValue.IsNil() {
		return ErrInvalidInput
	}

	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	return l.loadStructFromEnv(structValue)
}

func (l *Loader) loadStructFromEnv(s reflect.Value) error {
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	tags := getStructTags(s.Interface(), defaultTag)
	values := getTagsValuesFromEnv(tags)

	assignValuesToStruct(s, values...)
	return nil
}

func (l *Loader) loadEnvFile() error {
	if l.options.verbose {
		l.options.logger.Debugf("Loading environment variables from file:", l.options.filePath)
	}

	if err := loadEnvFile(l.options.filePath); err != nil {
		l.options.logger.Errorf("Failed to load environment file %s: %v", l.options.filePath, err)
		return err
	}

	return nil
}
