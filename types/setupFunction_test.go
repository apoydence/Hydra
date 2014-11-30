package types_test

import (
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"encoding"
)

var _ = Describe("SetupFunction", func() {
	Context("Builder", func() {
		var resultsChan chan FunctionInfo
		BeforeEach(func() {
			resultsChan = make(chan FunctionInfo, 1)
		})

		It("writes the correct information to the channel", func(done Done) {
			defer close(done)

			fake := func(s SetupFunction) {}
			f := NewSetupFunctionBuilder("someName", fake, resultsChan)
			fakeIn := make(chan encoding.BinaryMarshaler)
			fakeOut := make(chan encoding.BinaryMarshaler)

			go func() {
				defer GinkgoRecover()
				in, out := f.AsFilter("someParent", 5)
				var fin ReadOnlyChannel
				var fout WriteOnlyChannel
				fin = fakeIn
				fout = fakeOut
				Expect(in).To(BeEquivalentTo(fin))
				Expect(out).To(BeEquivalentTo(fout))
			}()

			fi := <-resultsChan
			fi.ReadChan() <- fakeIn
			fi.WriteChan() <- fakeOut

			Expect(fi.Name()).To(BeEquivalentTo("someName"))
			Expect(reflect.ValueOf(fi.Function()).Pointer()).To(Equal(reflect.ValueOf(fake).Pointer()))
			Expect(fi.Parent()).To(BeEquivalentTo("someParent"))
			Expect(fi.FuncType()).To(Equal(FILTER))
			Expect(fi.Instances()).To(Equal(5))
		}, 1)

		It("PRODUCER doesn't read from ReadChan", func(done Done) {
			defer close(done)

			fake := func(s SetupFunction) {}
			f := NewSetupFunctionBuilder("someName", fake, resultsChan)

			var fout WriteOnlyChannel
			fakeOut := make(chan encoding.BinaryMarshaler)
			fout = fakeOut

			go func() {
				fi := <-resultsChan
				fi.WriteChan() <- fakeOut
			}()

			out := f.AsProducer(5)

			Expect(out).To(BeEquivalentTo(fout))
		}, 1)

		It("CONSUMER doesn't read from WriteChan", func(done Done) {
			defer close(done)

			fake := func(s SetupFunction) {}
			f := NewSetupFunctionBuilder("someName", fake, resultsChan)

			var fin ReadOnlyChannel
			fakeIn := make(chan encoding.BinaryMarshaler)
			fin = fakeIn

			go func() {
				fi := <-resultsChan
				fi.ReadChan() <- fakeIn
			}()

			in := f.AsConsumer("someParent", 5)

			Expect(in).To(BeEquivalentTo(fin))
		}, 1)
	})
})
