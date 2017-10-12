package tools

import (
	"fmt"
	"reflect"

	"github.com/goline/errors"
)

type scanFn func(sf reflect.StructField, v reflect.Value)
type StructReader struct{}

// Read scans all fields inside target then apply each function
func (s *StructReader) Read(target interface{}, each scanFn) errors.Error {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	case reflect.Struct:
	default:
		return errors.New(ERR_TOOLS_STRUCT_READ_INVALID_TYPE, fmt.Sprintf("Reading invalid type. Got %s", t.Kind())).
			WithLevel(errors.LEVEL_WARN)
	}

	n := v.NumField()
	for i := 0; i < n; i++ {
		each(t.Field(i), v.Field(i))
	}

	return nil
}

// ReadTag scans through target for specific tags then apply each function
func (s *StructReader) ReadTag(target interface{}, each scanFn, tags ...string) errors.Error {
	if len(tags) == 0 {
		return nil
	}

	return s.Read(target, func(sf reflect.StructField, v reflect.Value) {
		for _, tag := range tags {
			if t, ok := sf.Tag.Lookup(tag); ok && t != "" {
				each(sf, v)
			}
		}
	})
}
