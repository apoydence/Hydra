package rpn

import "reflect"

type Invoker struct {
	rpnValues []Value
}

type Value struct {
	Value    interface{}
	ValueOk  bool
	Callable Callable
}

type Callable struct {
	Function func([]interface{}) []interface{}
	Inputs   []reflect.Type
	Outputs  []reflect.Type
}

func NewInvoker(rpnValues ...Value) *Invoker {
	return &Invoker{
		rpnValues: rpnValues,
	}
}

func (r *Invoker) Invoke(inputValue interface{}) interface{} {
	queue := r.replaceVariables([]interface{}{inputValue})

	for len(queue) > 1 || !queue[0].ValueOk {
		for i, value := range queue {
			if !value.ValueOk {
				lenInputs := len(value.Callable.Inputs)
				args := r.buildArgs(i, lenInputs, queue)
				result := value.Callable.Function(args)
				queue[i] = Value{
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

func (r *Invoker) replaceVariables(values []interface{}) []Value {
	queue := make([]Value, 0, len(r.rpnValues))
	for _, node := range r.rpnValues {
		if variable, ok := node.Value.(Variable); ok {
			queue = append(queue, Value{
				ValueOk: true,
				Value:   values[variable.Index],
			})
			continue
		}

		queue = append(queue, node)
	}
	return queue
}

func (r *Invoker) buildArgs(i, lenInputs int, queue []Value) []interface{} {
	var args []interface{}
	for _, value := range queue[i-lenInputs : i] {
		args = append(args, value.Value)
	}

	return args
}
