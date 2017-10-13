package tools_test

import (
	"github.com/goline/errors"
	"github.com/goline/tools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

type iniInputTest struct {
	DbHost string `ini:"db_host"`
	DbPort int64  `ini:"db_port"`

	RedisHost  string `ini:"host" ini_section:"redis"`
	RedisPort  int64  `ini:"port" ini_section:"redis"`
	RedisDelay int64  `ini:"delay" ini_section:"redis" default:"-1"`

	System iniSubInputTest `ini:"-"`

	JWT struct {
		Algorithm string `ini:"jwt_algorithm"`
	} `ini:"-"`

	Other struct {
		DevMode bool `ini:"dev_mode" ini_section:"other"`
	} `ini:"-" ini_section:"other" comment:"Other stuffs come here"`
}

type iniSubInputTest struct {
	Env       string `ini:"env" comment:"Available value are: prod, dev and test"`
	EnableLog bool   `ini:"enable_log"`
}

var _ = Describe("IniLoader", func() {
	It("LoadIni should return error code ERR_TOOLS_READ_INI_FAILED", func() {
		err := tools.LoadIni(".", nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(tools.ERR_TOOLS_READ_INI_FAILED))
	})

	It("LoadIni should load to input", func() {
		input := new(iniInputTest)
		err := tools.LoadIni("./fixtures/test.ini", input)
		Expect(err).To(BeNil())
		Expect(input.DbHost).To(Equal("127.0.0.1"))
		Expect(input.DbPort).To(Equal(int64(3306)))
		Expect(input.RedisHost).To(Equal("127.0.0.1"))
		Expect(input.RedisPort).To(Equal(int64(6379)))
		Expect(input.RedisDelay).To(Equal(int64(-1)))
		Expect(input.System.Env).To(Equal("test"))
		Expect(input.System.EnableLog).To(BeTrue())
		Expect(input.JWT.Algorithm).To(Equal("hs256"))
	})

	It("SaveIni should save input to INI file", func() {
		input := &iniInputTest{
			DbHost:    "db",
			DbPort:    3300,
			RedisHost: "redis.io",
			System: iniSubInputTest{
				Env:       "dev",
				EnableLog: false,
			},
			JWT: struct {
				Algorithm string `ini:"jwt_algorithm"`
			}{
				Algorithm: "es256",
			},
		}

		var err error
		err = tools.SaveIni("./fixtures/output.ini", input)
		Expect(err).To(BeNil())

		b, err := ioutil.ReadFile("./fixtures/output.ini")
		Expect(err).To(BeNil())
		Expect(string(b)).To(Equal(`db_host       = db
db_port       = 3300
; Available value are: prod, dev and test
env           = dev
enable_log    = false
jwt_algorithm = es256

[redis]
host  = redis.io
port  = %REDIS_PORT%
delay = %REDIS_DELAY%

; Other stuffs come here
[other]
dev_mode = false

`))
	})
})
