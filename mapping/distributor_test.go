package mapping

import (
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
)

var _ = Describe("Distributor", func() {
	Context("Non-Networked", func() {
		It("creates n-1 number of instances", func(done Done) {
			defer close(done)
			m := make(map[string]Mapper)
			var count int32
			count = 0

			fake := func(s SetupFunction) {
				atomic.AddInt32(&count, 1)
				s.AsProducer(1)
			}

			m["a"] = NewMapper(NewFunctionInfo("a", fake, "", 5, PRODUCER))
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

			m := make(map[string]Mapper)
			m["a"] = NewMapper(NewFunctionInfo("a", fakeP, "", 5, PRODUCER))
			m["b"] = NewMapper(NewFunctionInfo("b", fakeC, "", 5, CONSUMER))
			ma := m["a"].(*mapper)
			ma.consumers = append(ma.consumers, m["b"].Info())
			dm := NewDistributor()(m)

			Expect(len(dm.Instances("a"))).To(Equal(5))
			Expect(len(dm.Instances("b"))).To(Equal(5))
			Expect(dm.Consumers("a")).To(ConsistOf("b"))
		}, 1)
	})
})
