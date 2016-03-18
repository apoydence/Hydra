package rpn_test

import (
	"reflect"
	"unsafe"

	"github.com/apoydence/hydra/internal/rpn"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Linker", func() {
	var (
		Integer reflect.Type
		String  reflect.Type

		funcs    map[string]rpn.Callable
		rpnNodes []rpn.RawRpnNode
		link     *rpn.Linker
	)

	var callable = func(name string) func([]unsafe.Pointer) []unsafe.Pointer {
		return func([]unsafe.Pointer) []unsafe.Pointer {
			return []unsafe.Pointer{unsafe.Pointer(&name)}
		}
	}

	var callableEqual = func(x, y rpn.Callable) bool {
		return reflect.DeepEqual(x.Function(nil), y.Function(nil))
	}

	var parse = func(query string) []rpn.RawRpnNode {
		parser := rpn.NewParser()
		values, err := parser.Parse(query)
		Expect(err).ToNot(HaveOccurred())
		return values
	}

	var unsafeToInt = func(value unsafe.Pointer) int {
		return *(*int)(value)
	}

	BeforeEach(func() {
		Integer = reflect.TypeOf(0)
		String = reflect.TypeOf("")

		funcs = map[string]rpn.Callable{
			"FuncA": rpn.Callable{
				Function: callable("FuncA"),
				Inputs:   []reflect.Type{Integer},
				Outputs:  []reflect.Type{Integer},
			},
			"FuncB": rpn.Callable{
				Function: callable("FuncB"),
				Inputs:   []reflect.Type{Integer, Integer},
				Outputs:  []reflect.Type{Integer},
			},
			"FuncC": rpn.Callable{
				Function: callable("FuncC"),
				Inputs:   []reflect.Type{String},
				Outputs:  []reflect.Type{Integer},
			},
		}

		link = rpn.New(funcs)
	})

	Describe("Link()", func() {
		Context("single function with constant", func() {
			BeforeEach(func() {
				// FuncA(99)
				rpnNodes = []rpn.RawRpnNode{
					{
						ValueOk: true,
						Name:    "99",
					},
					{
						ValueOk: false,
						Name:    "FuncA",
					},
				}
			})

			It("does not return an error", func() {
				_, err := link.Link(rpnNodes)

				Expect(err).ToNot(HaveOccurred())
			})

			It("returns correct number of values", func() {
				values, _ := link.Link(rpnNodes)

				Expect(values).To(HaveLen(2))
			})

			It("returns the constant as an integer", func() {
				values, _ := link.Link(rpnNodes)

				Expect(values[0].ValueOk).To(BeTrue())
				Expect(unsafeToInt(values[0].Value)).To(Equal(99))
			})

			It("returns the expected function", func() {
				values, _ := link.Link(rpnNodes)

				Expect(values[1].ValueOk).To(BeFalse())
				Expect(callableEqual(values[1].Callable, funcs["FuncA"])).To(BeTrue())
			})

			Context("two functions with constant and variable", func() {
				BeforeEach(func() {
					// FuncB(FuncA(99), $0)
					rpnNodes = append(rpnNodes, []rpn.RawRpnNode{
						{
							ValueOk: true,
							Name:    "$0",
						},
						{
							ValueOk: false,
							Name:    "FuncB",
						},
					}...)
				})

				It("does not return an error", func() {
					_, err := link.Link(rpnNodes)

					Expect(err).ToNot(HaveOccurred())
				})

				It("returns the correct number of values", func() {
					values, _ := link.Link(rpnNodes)

					Expect(values).To(HaveLen(4))
				})

				It("returns the variable with the correct type", func() {
					values, _ := link.Link(rpnNodes)

					Expect(values[2].ValueOk).To(BeFalse())
					Expect(values[2].Variable).To(Equal(&rpn.Variable{
						Index: 0,
						Type:  Integer,
					}))
				})
			})

			Context("two functions with constant and variable with the variable first", func() {
				BeforeEach(func() {
					// FuncB($0, FuncA(99))
					rpnNodes = append([]rpn.RawRpnNode{
						{
							ValueOk: true,
							Name:    "$0",
						},
					},
						rpnNodes...)

					rpnNodes = append(rpnNodes, rpn.RawRpnNode{
						ValueOk: false,
						Name:    "FuncB",
					})
				})

				It("does not return an error", func() {
					_, err := link.Link(rpnNodes)

					Expect(err).ToNot(HaveOccurred())
				})

				It("returns the correct number of values", func() {
					values, _ := link.Link(rpnNodes)

					Expect(values).To(HaveLen(4))
				})

				It("returns the variable with the correct type", func() {
					values, _ := link.Link(rpnNodes)

					Expect(values[0].ValueOk).To(BeFalse())
					Expect(values[0].Variable).To(Equal(&rpn.Variable{
						Index: 0,
						Type:  Integer,
					}))
				})
			})
		})
	})

	DescribeTable("invalid equations", func(query string) {
		rpnNodes = parse(query)
		_, err := link.Link(rpnNodes)

		Expect(err).To(HaveOccurred())
	},
		Entry("invalid return type for input", "FuncC(FuncA(99))"),
		Entry("not enough arguments for a function", "FuncB(99)"),
		Entry("too many args for a function", "FuncA(99, 101)"),
		Entry("no functions", "99"),
		Entry("variables don't start at 0", "FuncA($1)"),
		Entry("variables aren't incremental", "FuncB($0, $5)"),
	)
})
