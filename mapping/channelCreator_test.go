package mapping_test

import (
	"encoding"
	. "github.com/apoydence/hydra/mapping"
	. "github.com/apoydence/hydra/testing_helpers"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChannelCreator", func() {
	Context("Creating a channel", func() {
		It("should use the correct buffer size", func(done Done) {
			defer close(done)
			ch := NewChannelCreator()(5)
			for i := 0; i < 5; i++ {
				ch <- NewIntMarshaler(i)
			}

			for i := 0; i < 5; i++ {
				<-ch
			}
			Expect(dataOnChannel(ch)).To(Equal(false))
		}, 1)
	})
})

func dataOnChannel(ch chan encoding.BinaryMarshaler) bool {
	timer := time.NewTimer(time.Millisecond * 500).C
	select {
	case <-timer:
		return false
	case <-ch:
		return true
	}
}
