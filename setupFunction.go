package hydra

type FunctionType int

type ReadOnlyChannel <-chan HashedData
type WriteOnlyChannel chan<- HashedData

type SetupFunction interface {
	AsProducer(instances int) WriteOnlyChannel
	AsFilter(parent string, instances int) (in ReadOnlyChannel, out WriteOnlyChannel)
	AsConsumer(parent string, instances int) ReadOnlyChannel
}

type setupFunction func(parent string, instances int, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel)

type SetupFunctionBuilder func(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction

const (
	PRODUCER FunctionType = iota
	FILTER
	CONSUMER
)

func buildSetupFunc(name string, f func(SetupFunction), c chan FunctionInfo) setupFunction {
	return func(parent string, instances int, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel) {

		fi := NewFunctionInfo(name, f, parent, instances, funcType)
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

func (sf setupFunction) AsProducer(instances int) WriteOnlyChannel {
	_, out := sf("", instances, PRODUCER)
	return out
}

func (sf setupFunction) AsFilter(parent string, instances int) (in ReadOnlyChannel, out WriteOnlyChannel) {
	in, out = sf(parent, instances, FILTER)
	return
}

func (sf setupFunction) AsConsumer(parent string, instances int) ReadOnlyChannel {
	in, _ := sf(parent, instances, CONSUMER)
	return in
}
