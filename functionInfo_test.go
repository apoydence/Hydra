package hydra

import (
	"reflect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionInfo", func() {
	Context("NewFunctionInfo", func(){
		It("has the correct information", func(){
			f := func(s SetupFunction){}
			info := NewFunctionInfo("someName", f, "someParent", PRODUCER)

			Expect(info.Name()).To(Equal("someName"))
			Expect(reflect.ValueOf(info.Function()).Pointer()).To(Equal(reflect.ValueOf(f).Pointer()))
			Expect(info.Parent()).To(Equal("someParent"))
			Expect(info.FuncType()).To(Equal(PRODUCER))
		})
	})
})
