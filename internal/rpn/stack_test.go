package rpn_test

import (
	"github.com/poy/hydra/internal/rpn"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stack", func() {
	var (
		stack *rpn.Stack
	)

	BeforeEach(func() {
		stack = rpn.NewStack()
	})

	Describe("Pop()", func() {
		Context("empty", func() {
			It("returns false", func() {
				_, ok := stack.Pop()

				Expect(ok).To(BeFalse())
			})

			Context("single entry", func() {
				BeforeEach(func() {
					stack.Push("a")
				})

				It("returns true", func() {
					_, ok := stack.Pop()

					Expect(ok).To(BeTrue())
				})

				It("returns the only value", func() {
					value, _ := stack.Pop()

					Expect(value).To(Equal("a"))
				})

				It("removes the only value", func() {
					stack.Pop()
					_, ok := stack.Pop()

					Expect(ok).To(BeFalse())
				})

				Context("two entries", func() {
					BeforeEach(func() {
						stack.Push("b")
					})

					It("returns the last value", func() {
						value, _ := stack.Pop()

						Expect(value).To(Equal("b"))
					})

					It("returns the second value on second Pop()", func() {
						stack.Pop()
						value, _ := stack.Pop()

						Expect(value).To(Equal("a"))
					})
				})
			})
		})
	})
})
