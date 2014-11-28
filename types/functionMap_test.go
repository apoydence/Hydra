package types_test

import (
	. "github.com/apoydence/hydra/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionMap", func() {
	Context("When adding", func(){
		It("Adds a FunctionInfo with an empty consumer slice", func(){
			m := NewFunctionMapBuilder()
			fi := NewFunctionInfo("a", nil, "", 5, PRODUCER)
			m.Add(fi)
			Expect(m.Info("a")).To(Equal(fi))
			Expect(len(m.Consumers("a"))).To(Equal(0))
		})
		It("Adds a consumer with a nil FunctionInfo", func(){
			m := NewFunctionMapBuilder()
			fi := NewFunctionInfo("a", nil, "", 5, CONSUMER)
			m.AddConsumer("producer", fi)
			Expect(m.Info("producer")).To(BeNil())
			Expect(len(m.Consumers("producer"))).To(Equal(1))
			Expect(m.Consumers("producer")[0]).To(Equal(fi))
		})
		It("Panics if the same function name is added more than once", func(){
			m := NewFunctionMapBuilder()
			fi := NewFunctionInfo("a", nil, "", 5, PRODUCER)
			m.Add(fi)
			Expect(func(){m.Add(fi)}).Should(Panic())
		})
	})

})
