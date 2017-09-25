package tools_test

import (
	"unicode/utf8"

	"github.com/goline/tools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String", func() {
	It("Random should return a string with 8 characters", func() {
		s := tools.Random(8)
		Expect(utf8.RuneCountInString(s)).To(Equal(8))
	})
})
