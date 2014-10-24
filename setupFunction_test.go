package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
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
			f := buildSetupFunc("someName", fake, resultsChan)
			fakeIn := make(chan HashedData)
			fakeOut := make(chan HashedData)

			go func() {
				defer GinkgoRecover()
				in, out := f("someParent", 5, FILTER)
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
			f := buildSetupFunc("someName", fake, resultsChan)

			var fout WriteOnlyChannel
			fakeOut := make(chan HashedData)
			fout = fakeOut

			go func() {
				fi := <-resultsChan
				fi.WriteChan() <- fakeOut
			}()

			_, out := f("someParent", 5, PRODUCER)
			Expect(out).To(BeEquivalentTo(fout))
		}, 1)

		It("CONSUMER doesn't read from WriteChan", func(done Done) {
			defer close(done)

			fake := func(s SetupFunction) {}
			f := buildSetupFunc("someName", fake, resultsChan)

			var fin ReadOnlyChannel
			fakeIn := make(chan HashedData)
			fin = fakeIn

			go func() {
				fi := <-resultsChan
				fi.ReadChan() <- fakeIn
			}()

			in, _ := f("someParent", 5, CONSUMER)
			
			Expect(in).To(BeEquivalentTo(fin))
		}, 1)
	})

	Context("Interface Implementation", func(){
		var(
			fake *fakeSetupFunction
			fakeSetup SetupFunction
		)

		BeforeEach(func() {
			fake = NewFakeSetupFunction()
			fakeSetup = setupFunction(fake.setup)
		})

		Context("AsProducer", func(){
			It("Returns the correct channel and FunctionType", func(){
				Expect(fakeSetup.AsProducer(5)).To(Equal(fake.out))
				Expect(fake.funcType).To(Equal(PRODUCER))
				Expect(fake.instances).To(Equal(5))
			})
		})

		Context("AsFilter", func(){
			It("Returns the correct channels, FunctionType and parent", func(){
				in, out := fakeSetup.AsFilter("fakeParent", 5)
				Expect(in).To(Equal(fake.in))
				Expect(out).To(Equal(fake.out))
				Expect(fake.funcType).To(Equal(FILTER))
				Expect(fake.parent).To(Equal("fakeParent"))
				Expect(fake.instances).To(Equal(5))
			})
		})

		Context("AsConsumer", func(){
			It("Returns the correct channel, FunctionType, and parent", func(){
				Expect(fakeSetup.AsConsumer("fakeParent", 5)).To(Equal(fake.in))
				Expect(fake.funcType).To(Equal(CONSUMER))
				Expect(fake.parent).To(Equal("fakeParent"))
				Expect(fake.instances).To(Equal(5))
			})
		})
	})
})

type fakeSetupFunction struct{
	parent string
	instances int
	funcType FunctionType
	in       ReadOnlyChannel
	out      WriteOnlyChannel
}

func NewFakeSetupFunction() *fakeSetupFunction {
	return &fakeSetupFunction{
		in:  make(chan HashedData),
		out: make(chan HashedData),
	}
}

func (f *fakeSetupFunction) setup (parent string, instances int, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel){
	f.parent = parent
	f.instances = instances
	f.funcType = funcType
	return f.in, f.out
}
