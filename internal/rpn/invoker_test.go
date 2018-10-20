package rpn_test

import (
	"reflect"
	"unsafe"

	"github.com/poy/hydra/internal/rpn"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Invoker", func() {
	var (
		Integer reflect.Type

		invoker   *rpn.Invoker
		rpnValues []*rpn.Value

		expectedInput   unsafe.Pointer
		expectedOutputA unsafe.Pointer
		expectedOutputB unsafe.Pointer

		calledName chan string
		calledIn   chan []unsafe.Pointer
		calledOut  chan []unsafe.Pointer
	)

	var funcBuilder = func(name string) func([]unsafe.Pointer) []unsafe.Pointer {
		return func(args []unsafe.Pointer) []unsafe.Pointer {
			calledName <- name
			calledIn <- args
			return <-calledOut
		}
	}

	var unsafeInt = func(i int) unsafe.Pointer {
		return unsafe.Pointer(&i)
	}

	BeforeEach(func() {
		Integer = reflect.TypeOf(0)

		expectedInput = unsafeInt(97)
		expectedOutputA = unsafeInt(99)
		expectedOutputB = unsafeInt(101)

		calledName = make(chan string, 100)
		calledIn = make(chan []unsafe.Pointer, 100)
		calledOut = make(chan []unsafe.Pointer, 100)
		rpnValues = nil

	})

	JustBeforeEach(func() {
		invoker = rpn.NewInvoker(rpnValues...)
	})

	Describe("Invoke()", func() {
		Context("single function", func() {
			BeforeEach(func() {

				// Value FuncA => FuncA(Value)
				rpnValues = []*rpn.Value{
					{
						Variable: &rpn.Variable{
							Index: 0,
							Type:  Integer,
						},
					},
					{
						Callable: rpn.Callable{
							Function: funcBuilder("FuncA"),
							Inputs:   []reflect.Type{Integer},
							Outputs:  []reflect.Type{Integer},
						},
					},
				}

				calledOut <- []unsafe.Pointer{expectedOutputA}
			})

			It("invokes the function", func() {
				invoker.Invoke(expectedInput)

				Expect(calledName).Should(Receive(Equal("FuncA")))
			})

			It("passes the expected arguments", func() {
				invoker.Invoke(expectedInput)

				Expect(calledIn).Should(Receive(Equal([]unsafe.Pointer{expectedInput})))
			})

			It("returns the expected values", func() {
				returnValue := invoker.Invoke(expectedInput)

				Expect(returnValue).Should(Equal(expectedOutputA))
			})

			Context("two functions", func() {
				BeforeEach(func() {

					// Value FuncA FuncB => FuncB(FuncA(Value))
					rpnValues = append(rpnValues,
						&rpn.Value{
							Callable: rpn.Callable{
								Function: funcBuilder("FuncB"),
								Inputs:   []reflect.Type{Integer},
								Outputs:  []reflect.Type{Integer},
							},
						})

					calledOut <- []unsafe.Pointer{expectedOutputB}
				})

				It("invokes the functions in order", func(done Done) {
					defer close(done)
					invoker.Invoke(expectedInput)

					Expect(calledName).Should(Receive(Equal("FuncA")))
					Expect(calledName).Should(Receive(Equal("FuncB")))
				})

				It("passes the expected arguments", func() {
					invoker.Invoke(expectedInput)

					Expect(calledIn).Should(Receive(Equal([]unsafe.Pointer{expectedInput})))
					Expect(calledIn).Should(Receive(Equal([]unsafe.Pointer{expectedOutputA})))
				})

				It("returns the expected values", func() {
					returnValue := invoker.Invoke(expectedInput)

					Expect(returnValue).Should(Equal(expectedOutputB))
				})
			})

			Context("function with two arguments", func() {
				BeforeEach(func() {

					// Value FuncA Value FuncB => FuncB(FuncA(Value), Value)
					rpnValues = append(rpnValues,
						[]*rpn.Value{
							{
								Variable: &rpn.Variable{
									Index: 0,
									Type:  Integer,
								},
							},
							{
								Callable: rpn.Callable{
									Function: funcBuilder("FuncB"),
									Inputs:   []reflect.Type{Integer, Integer},
									Outputs:  []reflect.Type{Integer},
								},
							},
						}...)

					calledOut <- []unsafe.Pointer{expectedOutputB}
				})

				It("invokes the functions in order", func(done Done) {
					defer close(done)
					invoker.Invoke(expectedInput)

					Expect(calledName).Should(Receive(Equal("FuncA")))
					Expect(calledName).Should(Receive(Equal("FuncB")))
				})

				It("passes the expected arguments", func() {
					invoker.Invoke(expectedInput)

					Expect(calledIn).Should(Receive(Equal([]unsafe.Pointer{expectedInput})))
					Expect(calledIn).Should(Receive(Equal([]unsafe.Pointer{expectedOutputA, expectedInput})))
				})
			})
		})
	})

})
