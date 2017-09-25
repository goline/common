package tools

import (
	"fmt"
	"reflect"

	"github.com/go-ini/ini"
	"github.com/goline/errors"
)

func LoadIni(file string, v interface{}) error {
	return new(iniLoader).Load(file, v)
}

type iniLoader struct {
	ini *ini.File
}

func (l *iniLoader) Load(file string, v interface{}) error {
	var err error
	l.ini, err = ini.InsensitiveLoad(file)
	if err != nil {
		return errors.New(ERR_TOOLS_LOAD_INI_FAILED, fmt.Sprintf("Failed to load INI. Got %s", err.Error()))
	}

	return l.inject(v)
}

func (l *iniLoader) inject(v interface{}) error {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Ptr:
	default:
		return errors.New(ERR_TOOLS_LOAD_INI_INVALID_ARGUMENT, fmt.Sprintf("Could not load data to %v", t.Kind()))
	}

	s := t.Elem()
	n := s.NumField()
	if n == 0 {
		return nil
	}

	var section, key string
	var ok bool
	vv := reflect.ValueOf(v).Elem()
	for i := 0; i < n; i++ {
		sf := s.Field(i)

		key, ok = sf.Tag.Lookup("ini")
		if ok == false {
			continue
		}

		section = ""
		if st, ok := sf.Tag.Lookup("ini_section"); ok == true {
			section = st
		}

		f := vv.Field(i)
		if f.CanSet() == false {
			continue
		}

		sec := l.ini.Section(section)
		if sec == nil {
			continue
		}

		value, err := sec.GetKey(key)
		if err != nil {
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(value.String())
		case reflect.Bool:
			vb, err := value.Bool()
			if err == nil {
				f.SetBool(vb)
			}
		case reflect.Int64:
			vi, err := value.Int64()
			if err == nil {
				f.SetInt(vi)
			}
		case reflect.Float64:
			vf, err := value.Float64()
			if err == nil {
				f.SetFloat(vf)
			}
		}
	}

	return nil
}
