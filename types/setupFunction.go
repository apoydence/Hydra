package types

type FunctionType int

type SetupFunction interface {
	AsProducer() ProducerBuilder
	AsFilter(parent string) FilterBuilder
	AsConsumer(parent string) ConsumerBuilder
	Instances(count int) SetupFunction
	WriteBufferSize(count int) SetupFunction
}

type ProducerBuilder interface{
	Build() WriteOnlyChannel
}

type FilterBuilder interface{
	Build() (in ReadOnlyChannel, out WriteOnlyChannel)
}

type ConsumerBuilder interface{
	Build() ReadOnlyChannel
}

type SetupFunctionBuilder func(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction

const (
	PRODUCER FunctionType = iota
	FILTER
	CONSUMER
)

type setup struct{
	name string
	fs func(SetupFunction)
	funcInfoChan chan FunctionInfo
	instances int
	bufferSize int
}

type setupProducer struct{
	s *setup
}

type setupFilter struct{
	s *setup
	parent string
}

type setupConsumer struct{
	s *setup
	parent string
}

func NewSetupFunctionBuilder(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction{
	return &setup{
		name: name,
		fs: f,
		funcInfoChan: c,
	}
}
/*
func NewSetupFunctionBuilder(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction {
	var setupF setupFunction
	setupF = func(parent string, instances int, funcType FunctionType) (in ReadOnlyChannel, out WriteOnlyChannel) {

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
	return setupF
}
*/
func (s *setup) AsProducer() ProducerBuilder{
	return &setupProducer{
		s: s,
	}
}

func (s *setup) AsFilter(parent string) FilterBuilder{
	return &setupFilter{
		s: s,
		parent: parent,
	}
}

func (s *setup) AsConsumer(parent string) ConsumerBuilder{
	return &setupConsumer{
		s: s,
		parent: parent,
	}
}

func (s *setup) Instances(count int) SetupFunction{
	s.instances = count
	return s
}

func (s *setup) WriteBufferSize(count int) SetupFunction{
	s.bufferSize = count
	return s
}

func submitFuncInfo(s *setup, parent string, funcType FunctionType) FunctionInfo {
	fi := NewFunctionInfo(s.name, s.fs, parent, s.instances, funcType)
	s.funcInfoChan <- fi
	return fi
}

func (sp *setupProducer) Build() WriteOnlyChannel{
	fi := submitFuncInfo(sp.s, "", PRODUCER)
	return <-fi.WriteChan()
}

func (sp *setupFilter) Build() (ReadOnlyChannel, WriteOnlyChannel){
	fi := submitFuncInfo(sp.s, sp.parent, FILTER)
	return <-fi.ReadChan(), <-fi.WriteChan()
}

func (sp *setupConsumer) Build() ReadOnlyChannel{
	fi := submitFuncInfo(sp.s, sp.parent, CONSUMER)
	return <-fi.ReadChan()
}
