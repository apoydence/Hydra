package types

import (
	"sync"
	"sync/atomic"
)

type FunctionType int

type SetupFunction interface {
	AsProducer() ProducerBuilder
	AsFilter(parent string) FilterBuilder
	AsConsumer(parent string) ConsumerBuilder
	SetInstances(count int) SetupFunction
	Instances() int
	SetWriteBufferSize(count int) SetupFunction
	WriteBufferSize() int
	SetName(name string) SetupFunction
	Name() string
}

type ProducerBuilder interface {
	Build() WriteOnlyChannel
}

type FilterBuilder interface {
	Build() (in ReadOnlyChannel, out WriteOnlyChannel)
}

type ConsumerBuilder interface {
	Build() ReadOnlyChannel
}

type SetupFunctionBuilder func(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction

const (
	PRODUCER FunctionType = iota
	FILTER
	CONSUMER
)

type setup struct {
	name         string
	fs           func(SetupFunction)
	funcInfoChan chan FunctionInfo
	instances    int32
	bufferSize   int32
	rwLock       *sync.RWMutex
}

type setupProducer struct {
	s *setup
}

type setupFilter struct {
	s      *setup
	parent string
}

type setupConsumer struct {
	s      *setup
	parent string
}

func NewSetupFunctionBuilder(name string, f func(SetupFunction), c chan FunctionInfo) SetupFunction {
	return &setup{
		name:         name,
		fs:           f,
		funcInfoChan: c,
		instances:    1,
		bufferSize:   0,
		rwLock:       &sync.RWMutex{},
	}
}

func (s *setup) AsProducer() ProducerBuilder {
	return &setupProducer{
		s: s,
	}
}

func (s *setup) AsFilter(parent string) FilterBuilder {
	return &setupFilter{
		s:      s,
		parent: parent,
	}
}

func (s *setup) AsConsumer(parent string) ConsumerBuilder {
	return &setupConsumer{
		s:      s,
		parent: parent,
	}
}

func (s *setup) SetInstances(count int) SetupFunction {
	atomic.StoreInt32(&s.instances, int32(count))
	return s
}

func (s *setup) Instances() int {
	return int(atomic.LoadInt32(&s.instances))
}

func (s *setup) SetWriteBufferSize(count int) SetupFunction {
	atomic.StoreInt32(&s.bufferSize, int32(count))
	return s
}

func (s *setup) WriteBufferSize() int {
	return int(atomic.LoadInt32(&s.bufferSize))
}

func (s *setup) Name() string {
	defer s.rwLock.RUnlock()
	s.rwLock.RLock()
	return s.name
}

func (s *setup) SetName(name string) SetupFunction {
	defer s.rwLock.RUnlock()
	s.rwLock.RLock()
	s.name = name
	return s
}

func submitFuncInfo(s *setup, parent string, funcType FunctionType) FunctionInfo {
	fi := NewFunctionInfo(s.Name(), s.fs, parent, s.Instances(), s.WriteBufferSize(), funcType)
	s.funcInfoChan <- fi
	return fi
}

func (sp *setupProducer) Build() WriteOnlyChannel {
	fi := submitFuncInfo(sp.s, "", PRODUCER)
	return <-fi.WriteChan()
}

func (sp *setupFilter) Build() (ReadOnlyChannel, WriteOnlyChannel) {
	fi := submitFuncInfo(sp.s, sp.parent, FILTER)
	return <-fi.ReadChan(), <-fi.WriteChan()
}

func (sp *setupConsumer) Build() ReadOnlyChannel {
	fi := submitFuncInfo(sp.s, sp.parent, CONSUMER)
	return <-fi.ReadChan()
}
