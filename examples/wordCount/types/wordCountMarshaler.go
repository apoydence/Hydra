package types

import "encoding"

type wordCountMarshaler map[string]uint32

func NewWordCountMarshaler(m map[string]uint32) encoding.BinaryMarshaler{
	var w wordCountMarshaler = m
	return w 
}

func (w wordCountMarshaler) MarshalBinary() (data []byte, err error){
	panic("Not implemented")
}

func ToMap(bm encoding.BinaryMarshaler) map[string]uint32{
	var w wordCountMarshaler = bm.(wordCountMarshaler)
	return w
}
