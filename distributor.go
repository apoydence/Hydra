package hydra

type Distributor func(FunctionMap) DistributedFunctionMap

func NewDistributor() Distributor{
	return distribute
}

func distribute(m FunctionMap) DistributedFunctionMap{
	dfm := make(map[string]*distMapper)
	resultsChan := make(chan FunctionInfo)
	for _, fm := range m{
		fi := fm.Info()
		sf := buildSetupFunc(fi.Name(), fi.Function(), resultsChan)
		instances := make([]FunctionInfo, 0)
		instances = append(instances, fi)
		for i:=0; i<fi.Instances()-1; i++{
			go fi.Function()(sf)
			instances = append(instances, <-resultsChan)
		}

		consumers := make([]string, 0)
		for _, f := range fm.Consumers(){
			consumers = append(consumers, f.Name())
		}

		dfm[fi.Name()] = newDistMapper(instances, consumers)
	}

	var funcMap distFunctionMap
	funcMap = dfm
	return funcMap
}
