package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChannelMapper", func() {
	Context("With a single instance of each function", func() {
		Context("With a linear path", func() {
			var (
				m map[string]Mapper
			)

			BeforeEach(func() {
				m = make(map[string]Mapper)
			})

			It("passes a linked channel for PRODUCER/FILTER to FILTER/CONSUMER", func(done Done) {
				defer close(done)

				a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
				b := NewFunctionInfo("b", nil, "", 1, CONSUMER)
				ma := NewMapper(a).(*mapper)
				m["a"] = ma
				ma.consumers = append(ma.consumers, b)

				go channelMapper(m)

				Expect(areStreamsLinked(<-a.WriteChan(), <-b.ReadChan())).To(BeTrue())
			}, 1)
		})

		Context("With a non-linear path", func() {
			var (
				m map[string]Mapper
			)

			BeforeEach(func() {
				m = make(map[string]Mapper)
			})

			It("passes a linked channel for PRODUCER/FILTER to FILTER/CONSUMER", func(done Done) {
				defer close(done)

				a := NewFunctionInfo("a", nil, "", 1, PRODUCER)
				a1 := NewFunctionInfo("a1", nil, "", 1, FILTER)
				a2 := NewFunctionInfo("a2", nil, "", 1, CONSUMER)
				a3 := NewFunctionInfo("a3", nil, "", 1, CONSUMER)
				ma := NewMapper(a).(*mapper)
				ma1 := NewMapper(a1).(*mapper)
				m["a"] = ma
				m["a1"] = ma1
				ma.consumers = append(ma.consumers, a1, a2)
				ma1.consumers = append(ma1.consumers, a3)

				go channelMapper(m)

				Expect(areStreamsLinked(<-a.WriteChan(), <-a1.ReadChan(), <-a2.ReadChan())).To(BeTrue())
				Expect(areStreamsLinked(<-a1.WriteChan(), <-a3.ReadChan())).To(BeTrue())
			})
		})
	})
})

func areStreamsLinked(out WriteOnlyChannel, ins ...ReadOnlyChannel) bool {
	for i := 0; i < 10; i++ {
		go func() {
			out <- NewHashedData(i, i)
		}()

		for _, in := range ins {
			test := <-in
			if test.Hash() != i || test.Data().(int) != i {
				return false
			}
		}
	}

	return true
}
