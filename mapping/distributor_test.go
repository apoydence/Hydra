package mapping_test

import (
	. "github.com/apoydence/hydra/mapping"
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
)

var _ = Describe("Distributor", func() {
	Context("Non-Networked", func() {
		It("creates n-1 number of instances", func(done Done) {
			defer close(done)
			m := NewFunctionMapBuilder()
			var count int32
			count = 0

			fake := func(s SetupFunction) {
				atomic.AddInt32(&count, 1)
				s.AsProducer().Build()
			}
			fia := NewFunctionInfo("a", fake, "", 5, 0, PRODUCER, nil)
			m.Add(fia)
			NewDistributor()(m)
			Eventually(func() int32 { return count }).Should(BeEquivalentTo(4))
		}, 1)

		It("stores each FunctionInfo for each instance with the proper mapping", func(done Done) {
			defer close(done)
			fakeP := func(s SetupFunction) {
				s.AsProducer().Build()
			}

			fakeC := func(s SetupFunction) {
				s.AsConsumer("b").Build()
			}

			m := NewFunctionMapBuilder()
			fia := NewFunctionInfo("a", fakeP, "", 5, 0, PRODUCER, nil)
			fib := NewFunctionInfo("b", fakeC, "", 5, 0, CONSUMER, nil)
			m.Add(fia)
			m.Add(fib)
			m.AddConsumer("a", fib)
			dm := NewDistributor()(m)

			Expect(len(dm.Instances("a"))).To(Equal(5))
			Expect(len(dm.Instances("b"))).To(Equal(5))
			Expect(dm.Consumers("a")).To(ConsistOf("b"))
		}, 1)
	})
})
