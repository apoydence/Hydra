package functionHandlers

import (
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionInvoker", func() {

	Context("when given multiple functions", func() {
		fakeSetupChan := make(chan SetupFunction)
		fakeSetupChanResult := make(chan bool)

		var (
			fakeSetupBuilder SetupFunctionBuilder
			fakeSetup        SetupFunction
			functionInvoker	 FunctionInvoker
		)

		fakeSetupComparer := func(s SetupFunction) bool {
			fakeSetupChan <- s.(SetupFunction)
			return <-fakeSetupChanResult
		}

		BeforeEach(func() {

			functionInvoker = NewFunctionInvoker()
			fakeSetupBuilder = func(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction {
				return fakeSetup
			}
			fakeSetup = &fakeSetupFunction{}

			go func(fakeSetup SetupFunction) {
				for sf := range fakeSetupChan {
					fakeSetupChanResult <- sf == fakeSetup
				}
			}(fakeSetup)
		})

		It("invokes each function once initially", func(done Done) {
			defer close(done)

			countChan := make(chan interface{})
			fake := func(sf SetupFunction) {
				countChan <- nil
			}
	
			functionInvoker(fakeSetupBuilder, fake, fake, fake)
			for i := 0; i < 3; i++ {
				<-countChan
			}

			Expect(tryRead(countChan)).ToNot(BeTrue())
		}, 1)

		It("invokes each function on their own go routine", func(done Done) {
			defer close(done)
			fake := func(sf SetupFunction) {
				x := make(chan struct{})
				<-x
			}

			functionInvoker(fakeSetupBuilder, fake, fake, fake)
		}, 1)

		It("passes the SetupFunction to each function", func(done Done) {
			defer close(done)

			fake := func(sf SetupFunction) {
				defer GinkgoRecover()
				Expect(fakeSetupComparer(sf)).To(Equal(true))
			}

			functionInvoker(fakeSetupBuilder, fake, fake, fake)
		}, 1)

		It("returns the same channel (non-nil) as the functions receive", func(done Done) {
			defer close(done)
			resultChan := make(chan chan FunctionInfo, 1)
			fsb := func(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction {
				resultChan <- c
				return fakeSetup
			}
			fake := func(sf SetupFunction) {}

			result := functionInvoker(fsb, fake)
			c := <-resultChan
			Expect(c).ToNot(BeNil())
			Expect(c).To(BeEquivalentTo(result))
		}, 1)
	})
})

func tryRead(c chan interface{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}

type fakeSetupFunction struct{
}

func (f *fakeSetupFunction) AsProducer(instances int) WriteOnlyChannel{
	panic("Not intended to be called")
}

func (f *fakeSetupFunction) AsFilter(parent string, instances int) (ReadOnlyChannel, WriteOnlyChannel){
	panic("Not intended to be called")
}

func (f *fakeSetupFunction) AsConsumer(parent string, instances int) ReadOnlyChannel{
	panic("Not intended to be called")
}
