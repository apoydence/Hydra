package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HashedData", func() {
	Context("NewHashedData", func() {
		It("has the correct data", func() {
			h := NewHashedData(99, 108)

			Expect(h.Hash()).To(Equal(99))
			Expect(h.Data()).To(Equal(108))
		})
	})
})
