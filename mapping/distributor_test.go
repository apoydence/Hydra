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
				s.AsProducer(1)
			}
			fia := NewFunctionInfo("a", fake, "", 5, PRODUCER)
			m.Add(fia)
			NewDistributor()(m)
			Eventually(func() int32 { return count }).Should(BeEquivalentTo(4))
		}, 1)

		It("stores each FunctionInfo for each instance with the proper mapping", func(done Done) {
			defer close(done)
			fakeP := func(s SetupFunction) {
				s.AsProducer(1)
			}

			fakeC := func(s SetupFunction) {
				s.AsConsumer("b", 1)
			}

			m := NewFunctionMapBuilder()
			fia := NewFunctionInfo("a", fakeP, "", 5, PRODUCER)
			fib := NewFunctionInfo("b", fakeC, "", 5, CONSUMER)
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
