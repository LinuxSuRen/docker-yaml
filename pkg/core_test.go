package pkg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config add command", func() {
	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Context("basic cases", func() {
		It("lack of name", func() {
			Expect(nil).To(HaveOccurred())
		})
	})
})
