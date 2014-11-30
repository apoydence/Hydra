package types

import "encoding"

type ReadOnlyChannel <-chan encoding.BinaryMarshaler
type WriteOnlyChannel chan<- encoding.BinaryMarshaler
