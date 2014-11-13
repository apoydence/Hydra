package hydra

type ChannelMapper func(m DistributedFunctionMap)

func channelMapper(m DistributedFunctionMap) {

	for _, funcName := range m.Functions() {
		instances := m.Instances(funcName)
		cs := createChannels(len(instances))
		setWriteChannels(m.Instances(funcName), cs)

		consumers := m.Consumers(funcName)
		numberOfConsumers := len(consumers)

		if numberOfConsumers == 1 {
			setReadChannels(m.Instances(consumers[0]), cs)
		} else if numberOfConsumers > 1 {
			panic("Not yet implemented...")
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
	go func() {
		for i, fi := range instances {
			fi.WriteChan() <- cs[i]
		}
	}()
}

func setReadChannels(instances []FunctionInfo, cs []chan HashedData) {
	instLength := len(instances)
	chanLength := len(cs)

	if instLength == chanLength {
		go setReadChannelsEqual(instances, cs)
	} else if instLength > chanLength {
		go setReadChannelsLess(instances, cs)
	} else {
		go setReadChannelsGreater(instances, cs)
	}
}

func setReadChannelsEqual(instances []FunctionInfo, cs []chan HashedData) {
	for i, fi := range instances {
		fi.ReadChan() <- cs[i]
	}
}

func setReadChannelsGreater(instances []FunctionInfo, cs []chan HashedData) {
	for i, fi := range instances {
		fi.ReadChan() <- cs[i%len(cs)]
	}
}

func setReadChannelsLess(instances []FunctionInfo, cs []chan HashedData) {
	for i, c := range cs {
		instances[i%len(instances)].ReadChan() <- c
	}
}
