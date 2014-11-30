package types

type FunctionMap interface {
	FunctionNames() []string
	Info(funcName string) FunctionInfo
	Consumers(funcName string) []FunctionInfo
}

type FunctionMapBuilder interface {
	FunctionMap
	Add(info FunctionInfo)
	AddConsumer(funcName string, parent FunctionInfo)
}

type mapper struct {
	info      FunctionInfo
	consumers []FunctionInfo
}

type functionMap map[string]*mapper

func NewFunctionMapBuilder() FunctionMapBuilder {
	var fm functionMap
	fm = make(map[string]*mapper)
	return fm
}

func (fm functionMap) FunctionNames() []string {
	names := make([]string, 0)
	for k := range fm {
		names = append(names, k)
	}
	return names
}

func (fm functionMap) Info(funcName string) FunctionInfo {
	if m, ok := fm[funcName]; ok {
		return m.info
	}
	return nil
}

func (fm functionMap) Consumers(funcName string) []FunctionInfo {
	if m, ok := fm[funcName]; ok {
		return m.consumers
	}
	return nil
}

func (fm functionMap) Add(info FunctionInfo) {
	if m, ok := fm[info.Name()]; ok {
		if m.info != nil {
			panic(info.Name() + " (function name) is being used twice")
		}

		m.info = info
	} else {
		fm[info.Name()] = &mapper{
			info:      info,
			consumers: make([]FunctionInfo, 0),
		}
	}

}
func (fm functionMap) AddConsumer(funcName string, parent FunctionInfo) {
	m, ok := fm[funcName]

	if !ok {
		m = &mapper{
			info:      nil,
			consumers: make([]FunctionInfo, 0),
		}
		fm[funcName] = m
	}

	m.consumers = append(m.consumers, parent)
}
