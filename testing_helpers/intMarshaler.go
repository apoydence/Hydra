package testing_helpers

import "encoding"

type IntMarshaler struct{
	Number int
}

func (i IntMarshaler) MarshalBinary()(data []byte, err error){
	panic("Not implemented...")
}

func NewIntMarshaler(i int) encoding.BinaryMarshaler{
	return IntMarshaler{
		Number: i,
	}
}
