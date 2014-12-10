package functionHandlers_test

import (
	. "github.com/apoydence/hydra/functionHandlers"
	"github.com/apoydence/hydra/mocks"
	. "github.com/apoydence/hydra/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("FunctionInvoker", func() {

	Context("when given multiple functions", func() {

		var (
			fakeSetupBuilder SetupFunctionBuilder
			fakeSetup        SetupFunction
			functionInvoker  FunctionInvoker
		)

		BeforeEach(func() {
			functionInvoker = NewFunctionInvoker()
			fakeSetupBuilder = func(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction {
				return fakeSetup
			}

			fakeSetup = mocks.NewMockSetupFunction(nil, nil)
			rand.Seed(time.Now().UnixNano())
			fakeSetup.SetName(strconv.Itoa(rand.Int()))
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
