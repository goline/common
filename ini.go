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

func (l *iniLoader) inject(input interface{}) error {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	switch t.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		return errors.New(ERR_TOOLS_LOAD_INI_INVALID_ARGUMENT, fmt.Sprintf("Could not load data to %v", t.Kind()))
	}

	l.injectTV(t, v)
	return nil
}

func (l *iniLoader) injectTV(t reflect.Type, v reflect.Value) {
	n := t.NumField()
	if n == 0 {
		return
	}

	var section, key string
	var ok bool
	for i := 0; i < n; i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		if vf.CanInterface() && (vf.Kind() == reflect.Struct || vf.Kind() == reflect.Ptr) {
			l.injectTV(vf.Type(), vf)
			continue
		}

		key, ok = tf.Tag.Lookup("ini")
		if ok == false {
			continue
		}

		section = ""
		if st, ok := tf.Tag.Lookup("ini_section"); ok == true {
			section = st
		}

		if vf.CanSet() == false {
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

		switch vf.Kind() {
		case reflect.String:
			vf.SetString(value.String())
		case reflect.Bool:
			vb, err := value.Bool()
			if err == nil {
				vf.SetBool(vb)
			}
		case reflect.Int64:
			vi, err := value.Int64()
			if err == nil {
				vf.SetInt(vi)
			}
		case reflect.Float64:
			vfl, err := value.Float64()
			if err == nil {
				vf.SetFloat(vfl)
			}
		}
	}
}
