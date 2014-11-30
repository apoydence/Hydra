package mapping

import "encoding"

type ChannelCreator func(bufferSize int) chan encoding.BinaryMarshaler

func NewChannelCreator() ChannelCreator {
	return channelCreator
}

func channelCreator(bufferSize int) chan encoding.BinaryMarshaler {
	return make(chan encoding.BinaryMarshaler, bufferSize)
}
