package hydra

import (
	"github.com/apoydence/hydra/functionHandlers"
	"github.com/apoydence/hydra/mapping"
	"github.com/apoydence/hydra/types"
)

type Scaffolding func(fs ...func(types.SetupFunction)) types.Canceller

func NewSetupScaffolding() Scaffolding {
	return func(fs ...func(types.SetupFunction)) types.Canceller {
		buildSetup := types.NewSetupFunctionBuilder
		funcInvoker := functionHandlers.NewFunctionInvoker()
		funcMapper := mapping.NewFunctionMapper()
		distributor := mapping.NewDistributor()
		chCreator := mapping.NewChannelCreator()
		chMapper := mapping.NewChannelMapper(chCreator)

		chanFuncInfo := funcInvoker(buildSetup, fs...)
		fmap := funcMapper(len(fs), chanFuncInfo)
		dfMap := distributor(fmap)
		chMapper(dfMap)

		cancellers := buildCancellerSlice(fmap)

		return func() {
			for _, c := range cancellers {
				c()
			}
		}
	}
}

func buildCancellerSlice(fmap types.FunctionMap) []types.Canceller {
	c := make([]types.Canceller, 0)
	for _, f := range fmap.FunctionNames() {
		c = append(c, fmap.Info(f).Cancel)
	}

	return c
}
