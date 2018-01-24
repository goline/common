package tools_test

import (
	"github.com/goline/errors"
	. "github.com/goline/tools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type sampleCopyFromInput1 struct {
	Age  uint64
	Name string
}

type sampleCopyToInput2 struct {
	Age   uint64
	Name  string
	Email string
}

type sampleCopyFromInput3 struct {
	I      sampleCopyFromInput1
	Gender bool
}

type sampleCopyToInput4 struct {
	Age    uint64
	Name   string
	Gender bool
}

type sampleCopyMapInput5 struct {
	Age  int64 `json:"age"`
	Info struct {
		Name string `json:"name"`
		Sex  bool   `json:"sex"`
	} `json:"info"`
}

var _ = Describe("Tools", func() {
	It("Copy must return error code ERR_TOOLS_COPY_TO_NON_POINTER", func() {
		f := &sampleCopyFromInput1{}
		t := sampleCopyToInput2{}
		err := Copy(f, t)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_TOOLS_COPY_TO_NON_POINTER))
	})

	It("Copy must return error code ERR_TOOLS_COPY_FROM_INVALID_TYPE", func() {
		f := "a string"
		t := &sampleCopyToInput2{}
		err := Copy(&f, t)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_TOOLS_COPY_FROM_INVALID_TYPE))
	})

	It("Copy must return nil", func() {
		f := &sampleCopyFromInput1{Age: 10, Name: "John"}
		t := &sampleCopyToInput2{Email: "e@mail.com"}
		err := Copy(f, t)
		Expect(err).To(BeNil())
		Expect(t.Age).To(Equal(uint64(10)))
		Expect(t.Name).To(Equal("John"))
		Expect(t.Email).To(Equal("e@mail.com"))
	})

	It("Copy must return nil", func() {
		f := &sampleCopyFromInput3{I: sampleCopyFromInput1{10, "John"}, Gender: true}
		t := &sampleCopyToInput4{}
		err := Copy(f, t)
		Expect(err).To(BeNil())
		Expect(t.Age).To(Equal(uint64(10)))
		Expect(t.Name).To(Equal("John"))
		Expect(t.Gender).To(BeTrue())
	})

	It("Copy must return error code ERR_TOOLS_COPY_TO_NON_POINTER", func() {
		f := make(map[string]interface{})
		t := sampleCopyToInput2{}
		err := CopyMap(f, t)
		Expect(err).NotTo(BeNil())
		Expect(err.(errors.Error).Code()).To(Equal(ERR_TOOLS_COPY_TO_NON_POINTER))
	})

	It("CopyMap must return nil", func() {
		f := make(map[string]interface{})
		f["age"] = int64(10)
		f["info"] = map[string]interface{}{"name": "John", "sex": true}
		t := &sampleCopyMapInput5{}
		err := CopyMap(f, t)
		Expect(err).To(BeNil())
		Expect(t.Age).To(Equal(int64(10)))
		Expect(t.Info.Name).To(Equal("John"))
		Expect(t.Info.Sex).To(BeTrue())
	})
})
