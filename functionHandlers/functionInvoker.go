package functionHandlers

import (
	"github.com/apoydence/hydra/types"

	"reflect"
	"runtime"
)

type FunctionInvoker func(sf types.SetupFunctionBuilder, fs ...func(types.SetupFunction)) <-chan types.FunctionInfo

func NewFunctionInvoker() FunctionInvoker{
	return funcInvoker
}

func funcInvoker(sfb types.SetupFunctionBuilder, fs ...func(types.SetupFunction)) <-chan types.FunctionInfo {
	c := make(chan types.FunctionInfo)
	for _, f := range fs {
		sf := sfb(getFunctionName(f), f, c)
		go f(sf)
	}

	return c
}

func getFunctionName(f func(types.SetupFunction)) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
