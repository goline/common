package tools

import (
	"fmt"
	"reflect"

	"github.com/go-ini/ini"
	"github.com/goline/errors"
	"strconv"
	"strings"
)

func LoadIni(file string, v interface{}) errors.Error {
	return new(IniReadWriter).Read(file, v)
}

func SaveIni(file string, v interface{}) errors.Error {
	return new(IniReadWriter).Write(file, v)
}

type IniReadWriter struct {
	ini *ini.File
}

func (l *IniReadWriter) Read(file string, v interface{}) errors.Error {
	var err error
	l.ini, err = ini.InsensitiveLoad(file)
	if err != nil {
		return errors.New(ERR_TOOLS_READ_INI_FAILED, "Failed to load INI").WithDebug(err.Error())
	}

	return l.inject(v)
}

func (l *IniReadWriter) Write(file string, v interface{}) errors.Error {
	l.ini = ini.Empty()

	if err := l.scan(v); err != nil {
		return err
	}

	if err := l.ini.SaveTo(file); err != nil {
		return errors.New(ERR_TOOLS_WRITE_INI_SAVE_FAILED, fmt.Sprintf("Unable to save input to file %s", file)).
			WithDebug(err.Error())
	}

	return nil
}

func (l *IniReadWriter) inject(input interface{}) errors.Error {
	return new(StructReader).ReadTag(input, l.injectField, "ini", "ini_section")
}

func (l *IniReadWriter) injectField(sf reflect.StructField, v reflect.Value) errors.Error {
	if v.CanInterface() && (v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr) {
		return new(StructReader).ReadTag(v.Interface(), func(tf reflect.StructField, vf reflect.Value) errors.Error {
			return l.injectField(tf, v.FieldByName(tf.Name))
		}, "ini", "ini_section")
	}

	if !v.CanSet() {
		return nil
	}

	sec := ""
	if st := sf.Tag.Get("ini_section"); st != "" {
		sec = st
	}

	section := l.ini.Section(sec)
	if section == nil {
		return errors.New(ERR_TOOLS_READ_INI_INVALID_SECTION, fmt.Sprintf("Section %s could not be found", sec))
	}

	key := sf.Tag.Get("ini")
	if key == "" {
		return nil
	}

	value, err := section.GetKey(key)
	if err != nil {
		// load default value
		value = section.Key(key)
		value.SetValue(sf.Tag.Get("default"))
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

	return nil
}

func (l *IniReadWriter) scan(input interface{}) errors.Error {
	return new(StructReader).ReadTag(input, l.scanField, "ini", "ini_section")
}

func (l *IniReadWriter) scanField(sf reflect.StructField, v reflect.Value) errors.Error {
	if v.CanInterface() && (v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr) {
		return new(StructReader).ReadTag(v.Interface(), func(tf reflect.StructField, vf reflect.Value) errors.Error {
			return l.scanField(tf, v.FieldByName(tf.Name))
		}, "ini", "ini_section")
	}

	sec := ""
	if st := sf.Tag.Get("ini_section"); st != "" {
		sec = st
	}

	var err error
	section := l.ini.Section(sec)
	if section == nil {
		section, err = l.ini.NewSection(sec)
		if err != nil {
			return errors.New(ERR_TOOLS_READ_INI_INVALID_SECTION, fmt.Sprintf("Section %s could not be created", sec)).
				WithDebug(err.Error())
		}
	}

	key := sf.Tag.Get("ini")
	if key == "" {
		return nil
	}

	var value string
	switch v.Kind() {
	case reflect.String:
		value = v.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = strconv.Itoa(int(v.Uint()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = strconv.Itoa(int(v.Int()))
	case reflect.Bool:
		if v.Bool() {
			value = "true"
		} else {
			value = "false"
		}
	default:
		return errors.New(ERR_TOOLS_WRITE_INI_INVALID_TYPE, fmt.Sprintf("Type %s is not supported", v.Kind().String()))
	}

	if value == "" || value == "0" {
		// ignore when empty, zero
		// TODO: please improve this
		if sec != "" {
			value = fmt.Sprintf("%s_%s", strings.ToUpper(sec), strings.ToUpper(key))
		} else {
			value = strings.ToUpper(key)
		}
	}

	k, err := section.NewKey(key, value)
	if err != nil {
		return errors.New(ERR_TOOLS_WRITE_INI_KEY_FAILED, fmt.Sprintf("Unable to write key %s with value %s", key, value)).
			WithDebug(err.Error())
	}

	comment := sf.Tag.Get("comment")
	if comment != "" {
		k.Comment = comment
	}

	return nil
}
