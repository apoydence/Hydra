package mapping

import(
	"github.com/apoydence/hydra/types"
)

type Distributor func(types.FunctionMap) types.DistributedFunctionMap

func NewDistributor() Distributor {
	return distribute
}

func distribute(m types.FunctionMap) types.DistributedFunctionMap {
	dfm := types.NewDistributedMap()
	resultsChan := make(chan types.FunctionInfo)
	for _, funcName := range m.FunctionNames() {
		fi := m.Info(funcName)
		sf := types.NewSetupFunctionBuilder(fi.Name(), fi.Function(), resultsChan)
		instances := make([]types.FunctionInfo, 0)
		instances = append(instances, fi)
		for i := 0; i < fi.Instances()-1; i++ {
			go fi.Function()(sf)
			instances = append(instances, <-resultsChan)
		}

		consumers := make([]string, 0)
		for _, f := range m.Consumers(funcName) {
			consumers = append(consumers, f.Name())
		}

		dfm.Add(fi.Name(), instances, consumers)
	}

	return dfm
}
