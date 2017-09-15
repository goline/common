package utils

import (
	"testing"
)

func TestNewIniLoader(t *testing.T) {
	l := NewIniLoader()
	if l == nil {
		t.Errorf("Expects l is not nil")
	}
}

type iniTest struct {
	DbHost string `ini:"db_host"`
	DbPort int64  `ini:"db_port"`

	RedisHost  string `ini:"host" ini_section:"redis"`
	RedisPort  int64  `ini:"port" ini_section:"redis"`
	RedisDelay int64  `ini:"delay" ini_section:"redis"`
}

func TestFactoryIniLoader_Load(t *testing.T) {
	v := new(iniTest)
	err := NewIniLoader().Load("./fixtures/test.ini", v)
	if err != nil {
		t.Errorf("Expects err is nil")
	} else if v.DbHost != "127.0.0.1" || v.DbPort != 3306 ||
		v.RedisHost != "127.0.0.1" || v.RedisPort != 6379 || v.RedisDelay != -1 {
		t.Errorf("Expects v has correct data. Got %v", v)
	}
}
