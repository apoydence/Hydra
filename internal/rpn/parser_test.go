package rpn_test

import (
	"fmt"

	"github.com/apoydence/hydra/internal/rpn"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	var (
		expectedQuery string
		parser        *rpn.Parser
	)

	BeforeEach(func() {
		parser = rpn.NewParser()
	})

	Describe("Parse()", func() {
		Context("single arguments", func() {
			Context("single function", func() {
				BeforeEach(func() {
					expectedQuery = "FuncA(%v)"
				})

				Context("number used as argument", func() {
					BeforeEach(func() {
						expectedQuery = fmt.Sprintf(expectedQuery, 99)
					})

					It("does not return an error", func() {
						_, err := parser.Parse(expectedQuery)

						Expect(err).ToNot(HaveOccurred())
					})

					It("returns 2 nodes", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result).To(HaveLen(2))
					})

					It("returns the value first", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result[0].ValueOk).To(BeTrue())
						Expect(result[0].Name).To(Equal("99"))
					})

					It("returns the function second", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result[1].ValueOk).To(BeFalse())
						Expect(result[1].Name).To(Equal("FuncA"))
					})

					Context("extra (but valid) parenthesis", func() {
						BeforeEach(func() {
							expectedQuery = fmt.Sprintf("(%s)", expectedQuery)
						})

						It("does not return an error", func() {
							_, err := parser.Parse(expectedQuery)

							Expect(err).ToNot(HaveOccurred())
						})
					})
				})

				Context("variable used as argument", func() {
					BeforeEach(func() {
						expectedQuery = fmt.Sprintf(expectedQuery, "$1")
					})

					It("does not return an error", func() {
						_, err := parser.Parse(expectedQuery)

						Expect(err).ToNot(HaveOccurred())
					})

					It("returns 2 nodes", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result).To(HaveLen(2))
					})

					It("returns the value first", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result[0].ValueOk).To(BeTrue())
						Expect(result[0].Name).To(Equal("$1"))
					})
				})

				Context("string used as argument", func() {
					BeforeEach(func() {
						expectedQuery = fmt.Sprintf(expectedQuery, `"some-string"`)
					})

					It("does not return an error", func() {
						_, err := parser.Parse(expectedQuery)

						Expect(err).ToNot(HaveOccurred())
					})

					It("returns 2 nodes", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result).To(HaveLen(2))
					})

					It("returns the value first", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result[0].ValueOk).To(BeTrue())
						Expect(result[0].Name).To(Equal("some-string"))
					})
				})

				Context("boolean used as argument", func() {
					BeforeEach(func() {
						expectedQuery = fmt.Sprintf(expectedQuery, "true")
					})

					It("does not return an error", func() {
						_, err := parser.Parse(expectedQuery)

						Expect(err).ToNot(HaveOccurred())
					})

					It("returns 2 nodes", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result).To(HaveLen(2))
					})

					It("returns the value first", func() {
						result, _ := parser.Parse(expectedQuery)

						Expect(result[0].ValueOk).To(BeTrue())
						Expect(result[0].Name).To(Equal("true"))
					})
				})

				DescribeTable("invalid syntax", func(query string) {
					_, err := parser.Parse(query)

					Expect(err).To(HaveOccurred())
				},
					Entry("extra right parenthesis", "(99))"),
					Entry("extra left parenthesis", "((99)"),
					Entry("invalid function name", "9invalid(99)"),
					Entry("use of ^", "FuncA(^99)"),
					Entry(`non-matching '"'`, `FuncString("d)`),
				)
			})
		})

		Context("multiple arguments", func() {
			Context("single function", func() {
				BeforeEach(func() {
					expectedQuery = "FuncA(99, -101)"
				})

				It("does not return an error", func() {
					_, err := parser.Parse(expectedQuery)

					Expect(err).ToNot(HaveOccurred())
				})

				It("returns 3 nodes", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result).To(HaveLen(3))
				})

				It("returns the first argument first", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[0].ValueOk).To(BeTrue())
					Expect(result[0].Name).To(Equal("99"))
				})

				It("returns the second argument second", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[1].ValueOk).To(BeTrue())
					Expect(result[1].Name).To(Equal("-101"))
				})

				It("returns the func argument last", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[2].ValueOk).To(BeFalse())
					Expect(result[2].Name).To(Equal("FuncA"))
				})
			})

			Context("multiple functions", func() {
				BeforeEach(func() {
					expectedQuery = "FuncA(99, FuncB(101, 103))"
				})

				It("does not return an error", func() {
					_, err := parser.Parse(expectedQuery)

					Expect(err).ToNot(HaveOccurred())
				})

				It("returns 5 nodes", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result).To(HaveLen(5))
				})

				It("returns the first argument first", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[0].ValueOk).To(BeTrue())
					Expect(result[0].Name).To(Equal("99"))
				})

				It("returns the first argument from the second func second", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[1].ValueOk).To(BeTrue())
					Expect(result[1].Name).To(Equal("101"))
				})

				It("returns the second argument from the third func third", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[2].ValueOk).To(BeTrue())
					Expect(result[2].Name).To(Equal("103"))
				})

				It("returns the inner func argument fourth", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[3].ValueOk).To(BeFalse())
					Expect(result[3].Name).To(Equal("FuncB"))
				})

				It("returns the outer func argument last", func() {
					result, _ := parser.Parse(expectedQuery)

					Expect(result[4].ValueOk).To(BeFalse())
					Expect(result[4].Name).To(Equal("FuncA"))
				})
			})

			DescribeTable("misplaced comma (',')", func(query string) {
				_, err := parser.Parse(query)

				Expect(err).To(HaveOccurred())
			},
				Entry("outside of parenthesis", "FuncA(9),"),
				Entry("outside of function", "(9,5)"),
				Entry("no arg after", "FuncA(9,)"),
				Entry("no arg before", "FuncA(,9)"),
			)
		})
	})

})
