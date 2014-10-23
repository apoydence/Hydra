package hydra

type ChannelMapper func(m FunctionMap)

func channelMapper(m FunctionMap) {
	for _, v := range m {
		c := make(chan HashedData)
		setWriteChannel(v.Info().WriteChan(), c)

		if len(v.Consumers()) == 1 {
			setReadChannel(v.Consumers()[0].ReadChan(), c)
		} else {
			cs := assembleChannels(v.Consumers())
			go distributeData(c, cs)
		}
	}
}

func setWriteChannel(c chan WriteOnlyChannel, d chan HashedData) {
	go func() {
		c <- d
	}()
}

func setReadChannel(c chan ReadOnlyChannel, d chan HashedData) {
	go func() {
		c <- d
	}()
}

func assembleChannels(fi []FunctionInfo) []chan HashedData {
	results := make([]chan HashedData, 0)
	for _, f := range fi {
		c := make(chan HashedData)
		results = append(results, c)
		setReadChannel(f.ReadChan(), c)
	}
	return results
}

func distributeData(dc chan HashedData, cs []chan HashedData) {
	defer func() {
		for _, c := range cs {
			close(c)
		}
	}()

	for data := range dc {
		for _, c := range cs {
			c <- data
		}
	}
}
