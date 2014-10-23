package hydra

type FunctionInfo interface {
	Name() string
	Function() func(SetupFunction)
	Parent() string
	FuncType() FunctionType
	ReadChan() chan ReadOnlyChannel
	WriteChan() chan WriteOnlyChannel
}

type functionInfo struct {
	name      string
	f         func(SetupFunction)
	parent    string
	funcType  FunctionType
	readChan  chan ReadOnlyChannel
	writeChan chan WriteOnlyChannel
}

func NewFunctionInfo(name string, f func(SetupFunction), parent string, funcType FunctionType) FunctionInfo {
	return &functionInfo{
		name:      name,
		f:         f,
		parent:    parent,
		funcType:  funcType,
		readChan:  make(chan ReadOnlyChannel),
		writeChan: make(chan WriteOnlyChannel),
	}
}

func (f *functionInfo) Name() string {
	return f.name
}

func (f *functionInfo) Function() func(SetupFunction) {
	return f.f
}

func (f *functionInfo) Parent() string {
	return f.parent
}

func (f *functionInfo) FuncType() FunctionType {
	return f.funcType
}

func (f *functionInfo) ReadChan() chan ReadOnlyChannel {
	return f.readChan
}

func (f *functionInfo) WriteChan() chan WriteOnlyChannel {
	return f.writeChan
}
