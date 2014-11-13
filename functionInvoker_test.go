package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("FunctionInvoker", func() {

	Context("when given multiple functions", func() {
		fakeSetupChan := make(chan setupFunction)
		fakeSetupChanResult := make(chan bool)

		var (
			fakeSetupBuilder SetupFunctionBuilder
			fakeSetup        setupFunction
		)

		fakeSetupComparer := func(s SetupFunction) bool {
			fakeSetupChan <- s.(setupFunction)
			return <-fakeSetupChanResult
		}

		BeforeEach(func() {

			fakeSetupBuilder = func(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction {
				return fakeSetup
			}

			fakeSetup = func(parent string, instances int, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel) {
				panic("Not intended to be called")
			}

			pointer := reflect.ValueOf(fakeSetup).Pointer()
			go func() {
				for sf := range fakeSetupChan {
					fakeSetupChanResult <- reflect.ValueOf(sf).Pointer() == pointer
				}
			}()
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
				Expect(fakeSetupComparer(sf)).To(Equal(true))
			}

			functionInvoker(fakeSetupBuilder, fake, fake, fake)
			f := <-c
			Expect(reflect.ValueOf(f).Pointer()).To(Equal(reflect.ValueOf(fakeSetup).Pointer()))
		}, 1)

		It("returns the same channel (non-nil) as the functions receive", func(done Done) {
			defer close(done)
			resultChan := make(chan chan FunctionInfo, 1)
			fsb := func(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction {
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
