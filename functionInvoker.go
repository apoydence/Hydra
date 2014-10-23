package hydra

import (
	"reflect"
	"runtime"
)

type FunctionInvoker func(sf SetupFunctionBuilder, fs ...func(SetupFunction)) <-chan FunctionInfo

func functionInvoker(sfb SetupFunctionBuilder, fs ...func(SetupFunction)) <-chan FunctionInfo {
	c := make(chan FunctionInfo)
	for _, f := range fs {
		sf := sfb(getFunctionName(f), f, c)
		go f(sf)
	}

	return c
}

func getFunctionName(f func(SetupFunction)) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
