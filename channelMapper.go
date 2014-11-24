package hydra

type ChannelMapper func(m DistributedFunctionMap)

func channelMapper(m DistributedFunctionMap) {
	for _, funcName := range m.Functions() {
		instances := m.Instances(funcName)
		if instances[0].FuncType() == CONSUMER {
			continue
		}

		cs := createChannels(len(instances))
		go setWriteChannels(m.Instances(funcName), cs)

		consumers := m.Consumers(funcName)
		numberOfConsumers := len(consumers)

		if numberOfConsumers == 1 {
			setReadChannels(m.Instances(consumers[0]), cs)
		} else if numberOfConsumers > 1 {
			cloneMatrix := cloneProducerChannels(numberOfConsumers, cs)
			for i, cloneCs := range cloneMatrix {
				setReadChannels(m.Instances(consumers[i]), cloneCs)
			}
		}
	}
}

func createChannels(count int) []chan HashedData {
	results := make([]chan HashedData, 0)
	for i := 0; i < count; i++ {
		results = append(results, make(chan HashedData))
	}
	return results
}

func setWriteChannels(instances []FunctionInfo, cs []chan HashedData) {
	for i, fi := range instances {
		go func(instance FunctionInfo, c chan HashedData) {
			instance.WriteChan() <- c
		}(fi, cs[i])
	}
}

func setReadChannels(consumerInstances []FunctionInfo, cs []chan HashedData) {
	consumerLength := len(consumerInstances)
	producerLength := len(cs)

	if consumerLength == producerLength {
		go setReadChannelsEqual(consumerInstances, cs)
	} else if consumerLength > producerLength {
		go setReadChannelsGreater(consumerInstances, cs)
	} else {
		combinedCs := channelCombiner(consumerLength, cs)
		go setReadChannelsEqual(consumerInstances, combinedCs)
	}
}

func cloneProducerChannels(numOfConsumers int, producerCh []chan HashedData) [][]chan HashedData {
	result := make([][]chan HashedData, 0)

	for i := 0; i < numOfConsumers; i++ {
		result = append(result, make([]chan HashedData, 0))
		for _ = range producerCh {
			clonedCh := make(chan HashedData)
			result[i] = append(result[i], clonedCh)
		}
	}

	for i, c := range producerCh {
		go cloneAcrossChannels(i, c, result)
	}

	return result
}

func cloneAcrossChannels(col int, ch chan HashedData, matrix [][]chan HashedData) {
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

func setSingleReadChan(fi FunctionInfo, c chan HashedData) {
	fi.ReadChan() <- c
}

func setReadChannelsEqual(instances []FunctionInfo, cs []chan HashedData) {
	for i, fi := range instances {
		fi.ReadChan() <- cs[i]
		close(fi.ReadChan())
	}
}

func setReadChannelsGreater(instances []FunctionInfo, cs []chan HashedData) {
	producerLen := len(cs)
	for index, consumer := range instances {
		consumer.ReadChan() <- cs[index%producerLen]
	}
}

func channelCombiner(consumerCount int, cs []chan HashedData) []chan HashedData {
	result := make([]chan HashedData, 0)
	doneChs := make([]chan interface{}, 0)
	counts := make([]int, 0)

	for i := 0; i < consumerCount; i++ {
		result = append(result, make(chan HashedData))
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

func dataCombiner(consumerCh, producerCh chan HashedData, doneCh chan interface{}) {
	for data := range producerCh {
		consumerCh <- data
	}
	doneCh <- nil
}

func closeCombinedChannels(count int, doneCh chan interface{}, ch chan HashedData) {
	for i := 0; i < count; i++ {
		<-doneCh
	}
	close(ch)
}
