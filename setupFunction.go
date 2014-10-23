package hydra

type FunctionType int

type ReadOnlyChannel <-chan HashedData
type WriteOnlyChannel chan<- HashedData

type SetupFunction interface {
	AsProducer() WriteOnlyChannel
	AsFilter(parent string) (in ReadOnlyChannel, out WriteOnlyChannel)
	AsConsumer(parent string) ReadOnlyChannel
}

type setupFunction func(parent string, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel)

type SetupFunctionBuilder func(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction

const (
	PRODUCER FunctionType = iota
	FILTER
	CONSUMER
)

func buildSetupFunc(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction {
	return func(parent string, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel) {

		fi := NewFunctionInfo(name, f, parent, funcType)
		c <- fi

		switch funcType {
		case PRODUCER:
			return nil, <-fi.WriteChan()
		case FILTER:
			return <-fi.ReadChan(), <-fi.WriteChan()
		case CONSUMER:
			return <-fi.ReadChan(), nil
		default:
			panic("Invalid type: " + string(funcType))
		}
	}
}

func (sf setupFunction) AsProducer() WriteOnlyChannel {
	_, out := sf("", PRODUCER)
	return out
}

func (sf setupFunction) AsFilter(parent string) (in ReadOnlyChannel, out WriteOnlyChannel) {
	in, out = sf(parent, FILTER)
	return
}

func (sf setupFunction) AsConsumer(parent string) ReadOnlyChannel {
	in, _ := sf(parent, CONSUMER)
	return in
}
