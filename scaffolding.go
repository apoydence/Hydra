package hydra

type Scaffolding func(fs ...func(SetupFunction))

func setupScaffolding() Scaffolding{
	return func(fs ...func(SetupFunction)){
		buildSetup := buildSetupFunc
		funcInvoker := functionInvoker
		funcMapper := mapFunctions
		chMapper := channelMapper

		chanFuncInfo := funcInvoker(buildSetup, fs...)
		fmap := funcMapper(len(fs), chanFuncInfo)
		chMapper(fmap)
	}
}
