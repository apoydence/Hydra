package hydra

type DistributedFunctionMap interface{
	Instances(name string)[]FunctionInfo
	Consumers(name string)[]string
}

type distMapper struct{
	instances []FunctionInfo
	consumers []string
}

type distFunctionMap map[string]*distMapper

func newDistMapper(instances []FunctionInfo, consumers []string) *distMapper{
	return &distMapper{
		instances : instances,
		consumers : consumers,
	}
}

func NewDistributedMap() DistributedFunctionMap{
	var dfm distFunctionMap
	dfm = make(map[string]*distMapper)
	return dfm
}

func (dm distFunctionMap) Instances(name string)[]FunctionInfo{
	if m, ok := dm[name]; ok{
		return m.instances
	}

	return nil
}

func (dm distFunctionMap) Consumers(name string)[]string{
	if m, ok := dm[name]; ok{
		return m.consumers
	}

	return nil
}
