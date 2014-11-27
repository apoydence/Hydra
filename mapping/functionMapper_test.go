package mapping

import (
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionMapper", func() {
	Context("FunctionMapper", func() {
		It("maps the correct key to the function", func() {
			fc := make(chan FunctionInfo, 3)
			fc <- NewFunctionInfo("a", nil, "", 1, PRODUCER)
			fc <- NewFunctionInfo("b", nil, "", 1, PRODUCER)
			fc <- NewFunctionInfo("c", nil, "", 1, PRODUCER)

			m := NewFunctionMapper()(3, fc)

			for k, v := range m {
				Expect(k).To(Equal(v.Info().Name()))
			}
		})

		It("maps each part to its consumer", func() {
			a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
			b := NewFunctionInfo("b", nil, "a", 1, FILTER)
			c := NewFunctionInfo("c", nil, "b", 1, CONSUMER)

			fc := make(chan FunctionInfo, 3)
			fc <- a
			fc <- b
			fc <- c

			m := NewFunctionMapper()(3, fc)

			Expect(len(m["a"].Consumers())).To(Equal(1))
			Expect(m["a"].Consumers()).To(ContainElement(b))
			Expect(len(m["b"].Consumers())).To(Equal(1))
			Expect(m["b"].Consumers()).To(ContainElement(c))
		})

		It("handles multiple consumers", func() {
			fc := make(chan FunctionInfo, 3)
			a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
			b := NewFunctionInfo("b", nil, "a", 1, CONSUMER)
			c := NewFunctionInfo("c", nil, "a", 1, CONSUMER)

			fc <- a
			fc <- b
			fc <- c

			m := NewFunctionMapper()(3, fc)

			Expect(len(m["a"].Consumers())).To(Equal(2))
			Expect(m["a"].Consumers()).To(ContainElement(b))
			Expect(m["a"].Consumers()).To(ContainElement(c))
		})

		It("doesn't just keep the producers/filters that have a consumer", func() {
			fc := make(chan FunctionInfo, 4)
			a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
			b := NewFunctionInfo("b", nil, "a", 1, CONSUMER)
			c := NewFunctionInfo("c", nil, "a", 1, FILTER)
			d := NewFunctionInfo("d", nil, "", 1, PRODUCER)

			fc <- a
			fc <- b
			fc <- c
			fc <- d

			m := NewFunctionMapper()(4, fc)

			Expect(len(m)).To(Equal(4))
			Expect(m["a"]).ToNot(BeNil())
			Expect(m["b"]).ToNot(BeNil())
			Expect(m["c"]).ToNot(BeNil())
			Expect(m["d"]).ToNot(BeNil())
		})

		It("panics with a name mismatch", func() {
			fc := make(chan FunctionInfo, 2)
			a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
			b := NewFunctionInfo("b", nil, "wrong", 1, CONSUMER)

			fc <- a
			fc <- b

			Expect(func() { NewFunctionMapper()(2, fc) }).To(Panic())
		})

		It("detects that the number of functions is wrong", func(done Done) {
			fc := make(chan FunctionInfo, 2)
			a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
			b := NewFunctionInfo("b", nil, "wrong", 1, CONSUMER)

			fc <- a
			fc <- b

			Expect(func() { NewFunctionMapper()(9, fc) }).To(Panic())
			close(done)
		}, 1)
	})
})
