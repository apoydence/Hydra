package rpn

type Stack struct {
	values []string
}

func NewStack() *Stack {
	return new(Stack)
}

func (s *Stack) Push(value string) {
	s.values = append(s.values, value)
}

func (s *Stack) Pop() (string, bool) {
	if len(s.values) <= 0 {
		return "", false
	}

	l := len(s.values) - 1
	value := s.values[l]

	s.values = s.values[:l]

	return value, true
}
