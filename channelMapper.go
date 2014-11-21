package hydra

type ChannelMapper func(m DistributedFunctionMap)

func channelMapper(m DistributedFunctionMap) {
	for _, funcName := range m.Functions() {
		instances := m.Instances(funcName)
		if instances[0].FuncType() == CONSUMER{
			return;
		}

		cs := createChannels(len(instances))
		go setWriteChannels(m.Instances(funcName), cs)

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
	for i, fi := range instances {
		setWriteChannel(fi, cs[i])
	}
}

func setWriteChannel(instance FunctionInfo, c chan HashedData){
	go func(){
		instance.WriteChan() <- c;
	}()
}

func setReadChannels(consumerInstances []FunctionInfo, cs []chan HashedData) {
	consumerLength := len(consumerInstances)
	producerLength := len(cs)

	if consumerLength == producerLength {
		go setReadChannelsEqual(consumerInstances, cs)
	} else if consumerLength > producerLength {
		go setReadChannelsGreater(consumerInstances, cs)
	} else {
		panic("Not yet implemented...")
		//go setReadChannelsLess(channelCombiner(consumerInstances), cs)
	}
}

func setSingleReadChan(fi FunctionInfo, c chan HashedData){
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
	for index, consumer := range instances{
		consumer.ReadChan() <- cs[index % producerLen]
	}
}

func setReadChannelsLess(outCs []chan chan HashedData, cs []chan HashedData) {
	defer func(){
		for _, c := range outCs{
			close(c)
		}
	}()

	for i, c := range cs {
		outCs[i%len(outCs)] <- c
	}
}
