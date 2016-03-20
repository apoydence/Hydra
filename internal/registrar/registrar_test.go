package registrar_test

import (
	"reflect"
	"unsafe"

	"github.com/apoydence/hydra/internal/registrar"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registrar", func() {

	var (
		integer reflect.Type
		reg     *registrar.Registrar

		invoker       chan string
		invokerArgs   chan []unsafe.Pointer
		invokerReturn chan []unsafe.Pointer
	)

	var buildInvoker = func(name string) registrar.Invoker {
		return registrar.Invoker(func(function interface{}) func([]unsafe.Pointer) []unsafe.Pointer {
			invoker <- name
			return func(args []unsafe.Pointer) []unsafe.Pointer {
				invokerArgs <- args
				return <-invokerReturn
			}
		})
	}

	BeforeEach(func() {
		invoker = make(chan string, 100)
		invokerArgs = make(chan []unsafe.Pointer, 100)
		invokerReturn = make(chan []unsafe.Pointer, 100)

		integer = reflect.TypeOf(0)
		reg = registrar.New()
	})

	Context("invoker has been registered", func() {
		BeforeEach(func() {
			ok := reg.RegisterInvoker([]reflect.Type{integer, integer}, []reflect.Type{integer}, buildInvoker("I,I->I"))
			Expect(ok).To(BeTrue())
		})

		Context("single function added", func() {
			var (
				expectedFuncName string
			)

			BeforeEach(func() {
				expectedFuncName = "SomeFunc1"
				close(invokerReturn)
			})

			It("maps the function to its name", func() {
				reg.Register(expectedFuncName, someFunc1)

				Expect(reg.Funcs()).To(HaveKey(expectedFuncName))
			})

			It("reads the proper input types of the function", func() {
				reg.Register(expectedFuncName, someFunc1)
				callable := reg.Funcs()[expectedFuncName]

				Expect(callable.Inputs).To(Equal([]reflect.Type{integer, integer}))
			})

			It("reads the proper output types of the function", func() {
				reg.Register(expectedFuncName, someFunc1)
				callable := reg.Funcs()[expectedFuncName]

				Expect(callable.Outputs).To(Equal([]reflect.Type{integer}))
			})

			It("uses the correct invoker", func() {
				reg.Register(expectedFuncName, someFunc1)
				callable := reg.Funcs()[expectedFuncName]
				callable.Function(nil)

				Expect(invoker).To(Receive(Equal("I,I->I")))
			})

			Context("invoker is registered twice", func() {
				Describe("RegisterInvoker()", func() {
					It("returns false", func() {
						ok := reg.RegisterInvoker([]reflect.Type{integer, integer}, []reflect.Type{integer}, buildInvoker("I,I->I"))
						Expect(ok).To(BeFalse())
					})
				})
			})

		})
	})

	DescribeTable("Register() panic cases", func(name string, function interface{}) {
		ok := reg.RegisterInvoker([]reflect.Type{integer, integer}, []reflect.Type{integer}, buildInvoker("I,I->I"))
		Expect(ok).To(BeTrue())
		f := func() {
			reg.Register(name, function)
		}
		Expect(f).To(Panic())
	},
		Entry("invoker not registered", "SomeName2", someFunc2),
		Entry("not a function", "invalid", 99),
		Entry("invalid function name", "-invalid", someFunc1),
	)
})

func someFunc1(a, b int) int {
	return 1
}

func someFunc2(a int) int {
	return 2
}
