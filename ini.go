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

func (l *iniLoader) inject(input interface{}) errors.Error {
	return new(StructReader).ReadTag(input, l.injectField, "ini", "ini_section")
}

func (l *iniLoader) injectField(sf reflect.StructField, v reflect.Value) errors.Error {
	if v.CanInterface() && (v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr) {
		return new(StructReader).ReadTag(v.Interface(), func(tf reflect.StructField, vf reflect.Value) errors.Error {
			return l.injectField(tf, v.FieldByName(tf.Name))
		}, "ini", "ini_section")
	} else {
		if !v.CanSet() {
			return nil
		}

		sec := ""
		if st := sf.Tag.Get("ini_section"); st != "" {
			sec = st
		}

		section := l.ini.Section(sec)
		if section == nil {
			return errors.New(ERR_TOOLS_LOAD_INI_INVALID_SECTION, fmt.Sprintf("Section %s could not be found", sec))
		}

		key := sf.Tag.Get("ini")
		if key == "" {
			return nil
		}

		value, err := section.GetKey(key)
		if err != nil {
			return nil
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(value.String())
		case reflect.Bool:
			vb, err := value.Bool()
			if err == nil {
				v.SetBool(vb)
			}
		case reflect.Int64:
			vi, err := value.Int64()
			if err == nil {
				v.SetInt(vi)
			}
		case reflect.Float64:
			vfl, err := value.Float64()
			if err == nil {
				v.SetFloat(vfl)
			}
		}
	}

	return nil
}
