package types

import "encoding"

type stringMarshaler struct {
	str string
}

func NewStringMarshaler(s string) encoding.BinaryMarshaler {
	return &stringMarshaler{
		str: s,
	}
}

func (s *stringMarshaler) MarshalBinary() (data []byte, err error) {
	panic("Not implemented")
}

func ToString(b encoding.BinaryMarshaler) string {
	return b.(*stringMarshaler).str
}
