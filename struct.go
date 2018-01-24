package tools

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/goline/errors"
)

type scanFn func(sf reflect.StructField, v reflect.Value) errors.Error
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
		if err := each(t.Field(i), v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

// ReadTag scans through target for specific tags then apply each function
func (s *StructReader) ReadTag(target interface{}, each scanFn, tags ...string) errors.Error {
	if len(tags) == 0 {
		return nil
	}

	return s.Read(target, func(sf reflect.StructField, v reflect.Value) errors.Error {
		for _, tag := range tags {
			if t, ok := sf.Tag.Lookup(tag); ok && t != "" {
				if err := each(sf, v); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func Copy(from interface{}, to interface{}) errors.Error {
	ft := reflect.TypeOf(from)
	tt := reflect.TypeOf(to)
	if tt.Kind() != reflect.Ptr {
		return errors.New(ERR_TOOLS_COPY_TO_NON_POINTER, "Destination of Copy must be a pointer")
	}
	tt = tt.Elem()

	var fv reflect.Value
	switch ft.Kind() {
	case reflect.Struct:
		fv = reflect.ValueOf(from)
	case reflect.Ptr:
		ft = ft.Elem()
		if ft.Kind() != reflect.Struct {
			return errors.New(ERR_TOOLS_COPY_FROM_INVALID_TYPE, "Source of Copy must be either a struct or a pointer")
		}
		fv = reflect.ValueOf(from).Elem()
	default:
		return errors.New(ERR_TOOLS_COPY_FROM_INVALID_TYPE, "Source of Copy must be either a struct or a pointer")
	}

	n, m := ft.NumField(), tt.NumField()
	if m == 0 || n == 0 {
		return nil
	}

	// Code could be shorter if I use v.FieldByName() method.
	// However, after do some benchmark, I see that performance of changes
	// is worse than this code. So I accept long code for better performance
	tv := reflect.ValueOf(to).Elem()
	for i := 0; i < n; i++ {
		fn := ft.Field(i).Name
		ff := fv.Field(i)
		if ff.CanInterface() && (ff.Kind() == reflect.Struct || ff.Kind() == reflect.Ptr) {
			Copy(ff.Interface(), to)
			continue
		}
		for j := 0; j < m; j++ {
			tn := tt.Field(j).Name
			tf := tv.Field(j)
			if fn != tn || ff.Kind() != tf.Kind() || !tf.CanSet() {
				continue
			}
			tf.Set(ff)
			break
		}
	}

	return nil
}

func CopyMap(from map[string]interface{}, to interface{}) errors.Error {
	tt := reflect.TypeOf(to)
	if tt.Kind() != reflect.Ptr {
		return errors.New(ERR_TOOLS_COPY_TO_NON_POINTER, "Destination of CopyMap must be a pointer")
	}

	return copyMapTV(from, tt.Elem(), reflect.ValueOf(to).Elem())
}

func copyMapTV(from map[string]interface{}, t reflect.Type, v reflect.Value) errors.Error {
	n := t.NumField()
	if n == 0 {
		return nil
	}

	rg, _ := regexp.Compile(`^(\w+)`)
	for key, value := range from {
		for i := 0; i < n; i++ {
			tf := t.Field(i)
			vf := v.Field(i)

			tag := tf.Tag.Get("json")
			if tag == "" || !rg.MatchString(tag) {
				continue
			}

			if vf.CanInterface() && (vf.Kind() == reflect.Struct || vf.Kind() == reflect.Ptr) {
				vm, ok := value.(map[string]interface{})
				if !ok {
					continue
				}

				copyMapTV(vm, vf.Type(), vf)
				continue
			}

			matches := rg.FindStringSubmatch(tag)
			if strings.Compare(key, matches[1]) != 0 {
				continue
			}

			if !vf.CanSet() {
				continue
			}

			// This custom allow to fix float64 -> integer value
			// It detects a number in $from but its mapping in $to is int64.
			// However, we expect int64 is final value, so make a tweak to do it
			// TODO: Find better solution for this case
			vv := reflect.ValueOf(value)
			switch vv.Kind() {
			case reflect.Float32, reflect.Float64:
				if vf.Kind() == reflect.Int64 {
					vf.SetInt(int64(vv.Float()))
				}
			default:
				vf.Set(vv)
			}

			break
		}
	}

	return nil
}
