package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionMapper", func() {
	Context("FunctionMapper", func(){
		It("maps correct key to function", func(){
			fc := make(chan FunctionInfo, 3)
			fc <- NewFunctionInfo("a", nil, "", PRODUCER)
			fc <- NewFunctionInfo("b", nil, "", PRODUCER)
			fc <- NewFunctionInfo("c", nil, "", PRODUCER)

			m := mapFunctions(3, fc)

			for k, v := range m{
				Expect(k).To(Equal(v.Info().Name()))
			}
		})

		It("maps all the parts to their consumers", func(){
			a := NewFunctionInfo("a", nil, "", PRODUCER)
			b := NewFunctionInfo("b", nil, "a", FILTER)
			c := NewFunctionInfo("c", nil, "b", CONSUMER)

			fc := make(chan FunctionInfo, 3)
			fc <- a
			fc <- b
			fc <- c

			m := mapFunctions(3, fc)

			Expect(len(m["a"].Consumers())).To(Equal(1))
			Expect(m["a"].Consumers()).To(ContainElement(b))
			Expect(len(m["b"].Consumers())).To(Equal(1))
			Expect(m["b"].Consumers()).To(ContainElement(c))
		})

		It("handles multiple consumers", func(){
			fc := make(chan FunctionInfo, 3)
			a := NewFunctionInfo("a", nil, "", PRODUCER)
			b := NewFunctionInfo("b", nil, "a", CONSUMER)
			c := NewFunctionInfo("c", nil, "a", CONSUMER)
		
			fc <- a
			fc <- b
			fc <- c

			m := mapFunctions(3, fc)

			Expect(len(m["a"].Consumers())).To(Equal(2))
			Expect(m["a"].Consumers()).To(ContainElement(b))
			Expect(m["a"].Consumers()).To(ContainElement(c))
		})

		It("keeps only producers/filters that have a consumer", func(){
			fc := make(chan FunctionInfo, 4)
			a := NewFunctionInfo("a", nil, "", PRODUCER)
			b := NewFunctionInfo("b", nil, "a", CONSUMER)
			c := NewFunctionInfo("c", nil, "a", FILTER)
			d := NewFunctionInfo("d", nil, "", PRODUCER)
		
			fc <- a
			fc <- b
			fc <- c
			fc <- d

			m := mapFunctions(4, fc)

			Expect(len(m)).To(Equal(1))
			Expect(m["a"]).ToNot(BeNil())
		})

		It("panics with a name mismatch", func(){
			fc := make(chan FunctionInfo, 2)
			a := NewFunctionInfo("a", nil, "", PRODUCER)
			b := NewFunctionInfo("b", nil, "wrong", CONSUMER)

			fc <- a
			fc <- b
			
			Expect(func(){mapFunctions(2, fc)}).To(Panic())
		})

		It("detects that the number of functions is wrong", func(done Done){
			fc := make(chan FunctionInfo, 2)
			a := NewFunctionInfo("a", nil, "", PRODUCER)
			b := NewFunctionInfo("b", nil, "wrong", CONSUMER)

			fc <- a
			fc <- b
			
			Expect(func(){mapFunctions(9, fc)}).To(Panic())
			close(done)
		}, 1)
	})
})
