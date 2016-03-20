package rpn

import (
	"reflect"
	"unsafe"
)

type Invoker struct {
	rpnValues []*Value
}

type Value struct {
	Value    unsafe.Pointer
	ValueOk  bool
	Variable *Variable
	Callable Callable
}

type Callable struct {
	Function func([]unsafe.Pointer) []unsafe.Pointer
	Inputs   []reflect.Type
	Outputs  []reflect.Type
}

func NewInvoker(rpnValues ...*Value) *Invoker {
	return &Invoker{
		rpnValues: rpnValues,
	}
}

func (r *Invoker) Invoke(inputValue unsafe.Pointer) unsafe.Pointer {
	queue := r.replaceVariables([]unsafe.Pointer{inputValue})

	for len(queue) > 1 || !queue[0].ValueOk {
		for i, value := range queue {
			if !value.ValueOk {
				lenInputs := len(value.Callable.Inputs)
				args := r.buildArgs(i, lenInputs, queue)
				result := value.Callable.Function(args)
				queue[i] = &Value{
					ValueOk: true,
					Value:   result[0],
				}

				queue = append(queue[:i-lenInputs], queue[i:]...)
				break
			}
		}
	}

	return queue[0].Value
}

func (r *Invoker) replaceVariables(values []unsafe.Pointer) []*Value {
	queue := make([]*Value, 0, len(r.rpnValues))
	for _, node := range r.rpnValues {
		if node.Variable != nil {
			queue = append(queue, &Value{
				ValueOk: true,
				Value:   values[node.Variable.Index],
			})
			continue
		}

		queue = append(queue, node)
	}
	return queue
}

func (r *Invoker) buildArgs(i, lenInputs int, queue []*Value) []unsafe.Pointer {
	var args []unsafe.Pointer
	for _, value := range queue[i-lenInputs : i] {
		args = append(args, value.Value)
	}

	return args
}
