package types

type HashedData interface {
	Hash() int
	Data() interface{}
}

type hashedData struct {
	hash int
	data interface{}
}

func NewHashedData(hash int, data interface{}) HashedData {
	return &hashedData{
		hash: hash,
		data: data,
	}
}

func (h *hashedData) Hash() int {
	return h.hash
}

func (h *hashedData) Data() interface{} {
	return h.data
}
