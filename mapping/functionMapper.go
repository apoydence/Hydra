package mapping

import (
	"github.com/apoydence/hydra/types"

	"time"
)

type Mapper interface {
	Info() types.FunctionInfo
	Consumers() []types.FunctionInfo
}

type mapper struct {
	info      types.FunctionInfo
	consumers []types.FunctionInfo
}

type FunctionMap map[string]Mapper

type FunctionMapper func(numOfFunctions int, functionChan <-chan types.FunctionInfo) FunctionMap

func NewMapper(info types.FunctionInfo) Mapper {
	return &mapper{
		info:      info,
		consumers: make([]types.FunctionInfo, 0),
	}
}

func (m *mapper) Info() types.FunctionInfo {
	return m.info
}

func (m *mapper) Consumers() []types.FunctionInfo {
	return m.consumers
}

func NewFunctionMapper() FunctionMapper{
	return mapFunctions
}

func mapFunctions(numOfFunctions int, functionChan <-chan types.FunctionInfo) FunctionMap {
	m := make(FunctionMap)
	for i := 0; i < numOfFunctions; i++ {
		funInfo := fetchNextFunctionInfo(functionChan)

		addToMap(funInfo, m)

		if funInfo.FuncType() != types.PRODUCER {
			parentInfo := fetchParent(funInfo.Parent(), m)
			parentInfo.consumers = append(parentInfo.consumers, funInfo)
		}
	}

	for k, v := range m {
		if v.Info() == nil {
			panic("Unknown function name: " + k)
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

func addToMap(info types.FunctionInfo, m FunctionMap) {
	var mapInfo *mapper
	i, ok := m[info.Name()]
	if ok {
		mapInfo = i.(*mapper)
		if i.Info() != nil {
			panic(info.Name() + " (function name) is being used twice")
		}

		mapInfo.info = info
	} else {
		m[info.Name()] = NewMapper(info)
	}
}

func fetchParent(parent string, m FunctionMap) *mapper {
	info, ok := m[parent]
	if ok {
		return info.(*mapper)
	}

	parentInfo := NewMapper(nil).(*mapper)
	m[parent] = parentInfo
	return parentInfo
}
