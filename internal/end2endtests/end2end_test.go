package end2end_test

import (
	"reflect"
	"unsafe"

	"github.com/apoydence/hydra"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("End2end", func() {
	var invokerIntIntInt = func(function interface{}) func([]unsafe.Pointer) []unsafe.Pointer {
		f := function.(func(int, int) int)
		return func(args []unsafe.Pointer) []unsafe.Pointer {
			arg1 := *(*int)(args[0])
			arg2 := *(*int)(args[1])
			output := f(arg1, arg2)
			return []unsafe.Pointer{unsafe.Pointer(&output)}
		}
	}

	var registerInvoker = func(h *hydra.Hydra) {
		integer := reflect.TypeOf(0)
		inputs := []reflect.Type{integer, integer}
		outputs := []reflect.Type{integer}
		h.RegisterInvoker(inputs, outputs, invokerIntIntInt)
	}

	var registerFunctions = func(h *hydra.Hydra) {
		h.Register("Add", Add)
		h.Register("Sub", Sub)
		h.Register("Mult", Mult)
		h.Register("Div", Div)
	}

	DescribeTable("matches normal invokation results", func(query string, variable, expectedOutput int) {
		h := hydra.New()
		registerInvoker(h)
		registerFunctions(h)
		function, err := h.BuildFunction(query)
		Expect(err).ToNot(HaveOccurred())

		output := function([]unsafe.Pointer{unsafe.Pointer(&variable)})
		Expect(output).To(HaveLen(1))
		outputInt := *(*int)(output[0])
		Expect(outputInt).To(Equal(expectedOutput))
	},
		Entry("Add($0, 5)", "Add($0, 5)", 10, Add(10, 5)),
		Entry("Sub(Add($0, 5), $0)", "Sub(Add($0, 5), $0)", 10, Sub(Add(10, 5), 10)),
		Entry("Sub(Add($0, Mult(5, -4)), $0)", "Sub(Add($0, Mult(5, -4)), $0)", 10, Sub(Add(10, Mult(5, -4)), 10)),
	)
})

func Add(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}

func Mult(a, b int) int {
	return a * b
}

func Div(a, b int) int {
	return a / b
}
