package tools_test

import (
	"reflect"

	"github.com/goline/errors"
	"github.com/goline/tools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StructReader", func() {
	It("Read allows to scan through struct", func() {
		type X struct {
			Name string `json:"name"`
		}

		r := new(tools.StructReader)
		r.Read(&X{Name: "John"}, func(sf reflect.StructField, v reflect.Value) errors.Error {
			if j := sf.Tag.Get("json"); j != "" {
				Expect(v.String()).To(Equal("John"))
			}

			return nil
		})
		r.Read(X{Name: "Marry"}, func(sf reflect.StructField, v reflect.Value) errors.Error {
			if j := sf.Tag.Get("json"); j != "" {
				Expect(v.String()).To(Equal("Marry"))
			}

			return nil
		})
	})

	It("ReadTag allows to read specific tag(s)", func() {
		type X struct {
			Name  string `json:"name" comment:"This is name"`
			Email string `json:"email"`
		}

		r := new(tools.StructReader)
		r.ReadTag(&X{Name: "John", Email: "e@mail.com"}, func(sf reflect.StructField, v reflect.Value) errors.Error {
			Expect(sf.Name).NotTo(Equal("Email"))
			Expect(sf.Tag.Get("comment")).To(Equal("This is name"))

			return nil
		}, "comment")
		r.ReadTag(X{Name: "Marry"}, func(sf reflect.StructField, v reflect.Value) errors.Error {
			Expect(sf.Name).NotTo(Equal("Email"))
			Expect(sf.Tag.Get("comment")).To(Equal("This is name"))

			return nil
		}, "comment")
	})
})
