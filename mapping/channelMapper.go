package mapping

import (
	"encoding"
	"github.com/apoydence/hydra/types"
)

type ChannelMapper func(m types.DistributedFunctionMap)

func NewChannelMapper(chanCreator ChannelCreator) ChannelMapper {
	cm := &chMapper{
		chanCreator: chanCreator,
	}
	return cm.channelMapper
}

type chMapper struct{
	chanCreator ChannelCreator
}

func (cm *chMapper) channelMapper(m types.DistributedFunctionMap) {
	for _, funcName := range m.Functions() {
		instances := m.Instances(funcName)
		if instances[0].FuncType() == types.CONSUMER {
			continue
		}

		cs := cm.createChannels(len(instances))
		go setWriteChannels(m.Instances(funcName), cs)

		consumers := m.Consumers(funcName)
		numberOfConsumers := len(consumers)

		if numberOfConsumers == 1 {
			cm.setReadChannels(m.Instances(consumers[0]), cs)
		} else if numberOfConsumers > 1 {
			cloneMatrix := cm.cloneProducerChannels(numberOfConsumers, cs)
			for i, cloneCs := range cloneMatrix {
				cm.setReadChannels(m.Instances(consumers[i]), cloneCs)
			}
		}
	}
}

func (cm *chMapper) createChannels(count int) []chan encoding.BinaryMarshaler {
	results := make([]chan encoding.BinaryMarshaler, 0)
	for i := 0; i < count; i++ {
		results = append(results, cm.chanCreator(0))
	}
	return results
}

func setWriteChannels(instances []types.FunctionInfo, cs []chan encoding.BinaryMarshaler) {
	for i, fi := range instances {
		go func(instance types.FunctionInfo, c chan encoding.BinaryMarshaler) {
			instance.WriteChan() <- c
		}(fi, cs[i])
	}
}

func (cm *chMapper) setReadChannels(consumerInstances []types.FunctionInfo, cs []chan encoding.BinaryMarshaler) {
	consumerLength := len(consumerInstances)
	producerLength := len(cs)

	if consumerLength == producerLength {
		go setReadChannelsEqual(consumerInstances, cs)
	} else if consumerLength > producerLength {
		go setReadChannelsGreater(consumerInstances, cs)
	} else {
		combinedCs := cm.channelCombiner(consumerLength, cs)
		go setReadChannelsEqual(consumerInstances, combinedCs)
	}
}

func (cm *chMapper) cloneProducerChannels(numOfConsumers int, producerCh []chan encoding.BinaryMarshaler) [][]chan encoding.BinaryMarshaler {
	result := make([][]chan encoding.BinaryMarshaler, 0)

	for i := 0; i < numOfConsumers; i++ {
		result = append(result, make([]chan encoding.BinaryMarshaler, 0))
		for _ = range producerCh {
			clonedCh := cm.chanCreator(0)
			result[i] = append(result[i], clonedCh)
		}
	}

	for i, c := range producerCh {
		go cloneAcrossChannels(i, c, result)
	}

	return result
}

func cloneAcrossChannels(col int, ch chan encoding.BinaryMarshaler, matrix [][]chan encoding.BinaryMarshaler) {
	defer func() {
		for _, row := range matrix {
			close(row[col])
		}
	}()

	for data := range ch {
		for _, row := range matrix {
			row[col] <- data
		}
	}
}

func setSingleReadChan(fi types.FunctionInfo, c chan encoding.BinaryMarshaler) {
	fi.ReadChan() <- c
}

func setReadChannelsEqual(instances []types.FunctionInfo, cs []chan encoding.BinaryMarshaler) {
	for i, fi := range instances {
		fi.ReadChan() <- cs[i]
		close(fi.ReadChan())
	}
}

func setReadChannelsGreater(instances []types.FunctionInfo, cs []chan encoding.BinaryMarshaler) {
	producerLen := len(cs)
	for index, consumer := range instances {
		consumer.ReadChan() <- cs[index%producerLen]
	}
}

func (cm *chMapper) channelCombiner(consumerCount int, cs []chan encoding.BinaryMarshaler) []chan encoding.BinaryMarshaler {
	result := make([]chan encoding.BinaryMarshaler, 0)
	doneChs := make([]chan interface{}, 0)
	counts := make([]int, 0)

	for i := 0; i < consumerCount; i++ {
		result = append(result, cm.chanCreator(0))
		doneChs = append(doneChs, make(chan interface{}))
		counts = append(counts, 0)
	}

	for i, c := range cs {
		adjustedIndex := i % consumerCount
		go dataCombiner(result[adjustedIndex], c, doneChs[adjustedIndex])
		counts[adjustedIndex]++
	}

	for i, c := range result {
		go closeCombinedChannels(counts[i], doneChs[i], c)
	}

	return result
}

func dataCombiner(consumerCh, producerCh chan encoding.BinaryMarshaler, doneCh chan interface{}) {
	for data := range producerCh {
		consumerCh <- data
	}
	doneCh <- nil
}

func closeCombinedChannels(count int, doneCh chan interface{}, ch chan encoding.BinaryMarshaler) {
	for i := 0; i < count; i++ {
		<-doneCh
	}
	close(ch)
}
