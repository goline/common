package tools_test

import (
	"github.com/goline/errors"
	"github.com/goline/tools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type iniInputTest struct {
	DbHost string `ini:"db_host"`
	DbPort int64  `ini:"db_port"`

	RedisHost  string `ini:"host" ini_section:"redis"`
	RedisPort  int64  `ini:"port" ini_section:"redis"`
	RedisDelay int64  `ini:"delay" ini_section:"redis"`

	System iniSubInputTest

	JWT struct {
		Algorithm string `ini:"jwt_algorithm"`
	}
}

type iniSubInputTest struct {
	Env       string `ini:"env"`
	EnableLog bool   `ini:"enable_log"`
}

var _ = Describe("IniLoader", func() {
	It("LoadIni should return error code ERR_TOOLS_LOAD_INI_FAILED", func() {
		err := tools.LoadIni(".", nil)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(tools.ERR_TOOLS_LOAD_INI_FAILED))
	})

	It("LoadIni should return error code ERR_TOOLS_LOAD_INI_INVALID_ARGUMENT", func() {
		err := tools.LoadIni("./fixtures/test.ini", "string")
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(tools.ERR_TOOLS_LOAD_INI_INVALID_ARGUMENT))
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
})
