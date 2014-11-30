package mapping

import (
	"github.com/apoydence/hydra/types"

	"time"
)

type FunctionMapper func(numOfFunctions int, functionChan <-chan types.FunctionInfo) types.FunctionMap

func NewFunctionMapper() FunctionMapper {
	return mapFunctions
}

func mapFunctions(numOfFunctions int, functionChan <-chan types.FunctionInfo) types.FunctionMap {
	m := types.NewFunctionMapBuilder()
	for i := 0; i < numOfFunctions; i++ {
		funInfo := fetchNextFunctionInfo(functionChan)
		m.Add(funInfo)

		if funInfo.FuncType() != types.PRODUCER {
			m.AddConsumer(funInfo.Parent(), funInfo)
		}
	}

	for _, funcName := range m.FunctionNames() {
		if m.Info(funcName) == nil {
			panic("Unknown function name: " + funcName)
		}
	}

	return m
}

func fetchNextFunctionInfo(c <-chan types.FunctionInfo) types.FunctionInfo {
	t := time.NewTicker(500 * time.Millisecond)
	select {
	case _ = <-t.C:
		panic("Waiting for functions has timed out...")
	case f := <-c:
		return f
	}
}
