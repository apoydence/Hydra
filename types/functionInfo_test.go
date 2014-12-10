package types

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("FunctionInfo", func() {
	Context("NewFunctionInfo", func() {
		It("has the correct information", func() {
			f := func(s SetupFunction) {}
			called := false
			c := func() {
				called = true
			}
			info := NewFunctionInfo("someName", f, "someParent", 1, 2, PRODUCER, c)
			info.Cancel()

			Expect(info.Name()).To(Equal("someName"))
			Expect(reflect.ValueOf(info.Function()).Pointer()).To(Equal(reflect.ValueOf(f).Pointer()))
			Expect(info.Parent()).To(Equal("someParent"))
			Expect(info.FuncType()).To(Equal(PRODUCER))
			Expect(info.Instances()).To(Equal(1))
			Expect(info.WriteBufferSize()).To(Equal(2))
			Expect(called).To(BeTrue())
		})
	})
})
