package autoenv

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNilLoader = errors.New("global loader is nil")
	ErrNilInput  = errors.New("input is nil")
)

type errUnsupportedKind struct {
	kind reflect.Kind
}

func (e *errUnsupportedKind) Error() string {
	return fmt.Sprintf("unsupported kind: %s", e.kind)
}

func IsUnsupportedKindError(err error) bool {
	var errUnsupportedKind *errUnsupportedKind
	ok := errors.As(err, &errUnsupportedKind)
	return ok
}
