package rpn_test

import (
	"reflect"

	"github.com/apoydence/hydra/internal/rpn"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Invoker", func() {
	var (
		Integer reflect.Type

		invoker   *rpn.Invoker
		rpnValues []rpn.Value

		expectedInput   int
		expectedOutputA int
		expectedOutputB int

		calledName chan string
		calledIn   chan []interface{}
		calledOut  chan []interface{}
	)

	var funcBuilder = func(name string) func([]interface{}) []interface{} {
		return func(args []interface{}) []interface{} {
			calledName <- name
			calledIn <- args
			return <-calledOut
		}
	}

	BeforeEach(func() {
		Integer = reflect.TypeOf(0)

		expectedInput = 97
		expectedOutputA = 99
		expectedOutputB = 101

		calledName = make(chan string, 100)
		calledIn = make(chan []interface{}, 100)
		calledOut = make(chan []interface{}, 100)
		rpnValues = nil

	})

	JustBeforeEach(func() {
		invoker = rpn.NewInvoker(rpnValues...)
	})

	Describe("Invoke()", func() {
		Context("single function", func() {
			BeforeEach(func() {

				// Value FuncA => FuncA(Value)
				rpnValues = []rpn.Value{
					{
						ValueOk: true,
						Value:   rpn.Placeholder,
					},
					{
						Callable: rpn.Callable{
							Function: funcBuilder("FuncA"),
							Inputs:   []reflect.Type{Integer},
							Outputs:  []reflect.Type{Integer},
						},
					},
				}

				calledOut <- []interface{}{expectedOutputA}
			})

			It("invokes the function", func() {
				invoker.Invoke(expectedInput)

				Expect(calledName).Should(Receive(Equal("FuncA")))
			})

			It("passes the expected arguments", func() {
				invoker.Invoke(expectedInput)

				Expect(calledIn).Should(Receive(Equal([]interface{}{expectedInput})))
			})

			It("returns the expected values", func() {
				returnValue := invoker.Invoke(expectedInput)

				Expect(returnValue).Should(Equal(expectedOutputA))
			})

			Context("two functions", func() {
				BeforeEach(func() {

					// Value FuncA FuncB => FuncB(FuncA(Value))
					rpnValues = append(rpnValues,
						rpn.Value{
							Callable: rpn.Callable{
								Function: funcBuilder("FuncB"),
								Inputs:   []reflect.Type{Integer},
								Outputs:  []reflect.Type{Integer},
							},
						})

					calledOut <- []interface{}{expectedOutputB}
				})

				It("invokes the functions in order", func(done Done) {
					defer close(done)
					invoker.Invoke(expectedInput)

					Expect(calledName).Should(Receive(Equal("FuncA")))
					Expect(calledName).Should(Receive(Equal("FuncB")))
				})

				It("passes the expected arguments", func() {
					invoker.Invoke(expectedInput)

					Expect(calledIn).Should(Receive(Equal([]interface{}{expectedInput})))
					Expect(calledIn).Should(Receive(Equal([]interface{}{expectedOutputA})))
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
						[]rpn.Value{
							{
								ValueOk: true,
								Value:   rpn.Placeholder,
							},
							{
								Callable: rpn.Callable{
									Function: funcBuilder("FuncB"),
									Inputs:   []reflect.Type{Integer, Integer},
									Outputs:  []reflect.Type{Integer},
								},
							},
						}...)

					calledOut <- []interface{}{expectedOutputB}
				})

				It("invokes the functions in order", func(done Done) {
					defer close(done)
					invoker.Invoke(expectedInput)

					Expect(calledName).Should(Receive(Equal("FuncA")))
					Expect(calledName).Should(Receive(Equal("FuncB")))
				})

				It("passes the expected arguments", func() {
					invoker.Invoke(expectedInput)

					Expect(calledIn).Should(Receive(Equal([]interface{}{expectedInput})))
					Expect(calledIn).Should(Receive(Equal([]interface{}{expectedOutputA, expectedInput})))
				})
			})
		})
	})

})
