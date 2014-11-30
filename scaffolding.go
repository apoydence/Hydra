package hydra

import (
	"github.com/apoydence/hydra/functionHandlers"
	"github.com/apoydence/hydra/mapping"
	"github.com/apoydence/hydra/types"
)

type Scaffolding func(fs ...func(types.SetupFunction))

func NewSetupScaffolding() Scaffolding {
	return func(fs ...func(types.SetupFunction)) {
		buildSetup := types.NewSetupFunctionBuilder
		funcInvoker := functionHandlers.NewFunctionInvoker()
		funcMapper := mapping.NewFunctionMapper()
		distributor := mapping.NewDistributor()
		chMapper := mapping.NewChannelMapper()

		chanFuncInfo := funcInvoker(buildSetup, fs...)
		fmap := funcMapper(len(fs), chanFuncInfo)
		dfMap := distributor(fmap)
		chMapper(dfMap)
	}
}
