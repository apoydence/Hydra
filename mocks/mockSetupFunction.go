package mocks

import (
	"encoding"
	. "github.com/apoydence/hydra/types"
)

func NewMockSetupFunction(in, out chan encoding.BinaryMarshaler) SetupFunction {
	c := make(chan FunctionInfo)
	go rxFuncInfo(in, out, c)
	return NewSetupFunctionBuilder("", nil, c)
}

func rxFuncInfo(in, out chan encoding.BinaryMarshaler, c chan FunctionInfo) {
	fi := <-c
	switch fi.FuncType() {
	case PRODUCER:
		fi.WriteChan() <- out
		break
	case FILTER:
		fi.ReadChan() <- in
		fi.WriteChan() <- out
		break
	case CONSUMER:
		fi.ReadChan() <- in
		break
	}
}
