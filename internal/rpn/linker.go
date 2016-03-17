package rpn

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

type Linker struct {
	funcs map[string]Callable
}

type Variable struct {
	Index int
	Type  reflect.Type
}

type indexedNode struct {
	RawRpnNode
	index     int
	traversed bool
}

func New(funcs map[string]Callable) *Linker {
	return &Linker{
		funcs: funcs,
	}
}

func (l *Linker) Link(nodes []RawRpnNode) ([]*Value, error) {
	decoratedNodes := l.decorateNodes(nodes)
	values := make([]*Value, len(nodes))

	for len(decoratedNodes) > 1 || l.keepGoing(decoratedNodes) {
		if l.onlyArgs(decoratedNodes) {
			return nil, fmt.Errorf("Extra arguments")
		}

		for i, value := range decoratedNodes {
			if !value.ValueOk && !value.traversed {
				value.traversed = true
				callable, _ := l.funcs[value.Name]
				values[value.index] = &Value{
					Callable: callable,
				}
				lenInputs := len(callable.Inputs)

				args, ok := l.buildArgs(i, lenInputs, decoratedNodes)
				if !ok {
					return nil, fmt.Errorf("Expected %d args for %s", lenInputs, value.Name)
				}

				if err := l.convertArgs(args, callable.Inputs, values); err != nil {
					return nil, err
				}

				decoratedNodes = append(decoratedNodes[:i-lenInputs], decoratedNodes[i:]...)
				break
			}
		}
	}

	return l.validateResult(values)
}

func (l *Linker) validateResult(values []*Value) ([]*Value, error) {
	if err := l.hasFunction(values); err != nil {
		return nil, err
	}

	return values, l.correctVariableIndexes(values)
}

func (l *Linker) hasFunction(values []*Value) error {
	for _, value := range values {
		if !value.ValueOk {
			return nil
		}
	}
	return fmt.Errorf("No functions")
}

func (l *Linker) correctVariableIndexes(values []*Value) error {
	var vars []int

	for _, v := range values {
		if !v.ValueOk {
			continue
		}

		variable, ok := v.Value.(Variable)
		if !ok {
			continue
		}

		vars = append(vars, variable.Index)
	}

	if len(vars) == 0 {
		return nil
	}

	sort.Sort(sort.IntSlice(vars))
	for i, j := range vars {
		if i != j {
			return fmt.Errorf("variable numbers aren't incremental")
		}
	}

	return nil
}

func (l *Linker) buildArgs(i, lenInputs int, nodes []*indexedNode) ([]*indexedNode, bool) {
	if i-lenInputs < 0 {
		return nil, false
	}

	return nodes[i-lenInputs : i], true
}

func (l *Linker) onlyArgs(nodes []*indexedNode) bool {
	for _, node := range nodes {
		if !node.traversed && !node.ValueOk {
			return false
		}
	}
	return true
}

func (l *Linker) keepGoing(nodes []*indexedNode) bool {
	if len(nodes) == 0 {
		return false
	}

	node := nodes[0]
	return !node.RawRpnNode.ValueOk && !node.traversed
}

func (l *Linker) decorateNodes(nodes []RawRpnNode) []*indexedNode {
	result := make([]*indexedNode, 0, len(nodes))
	for i, n := range nodes {
		result = append(result, &indexedNode{
			RawRpnNode: n,
			index:      i,
		})
	}
	return result
}

func (l *Linker) convertArgs(nodes []*indexedNode, inputs []reflect.Type, values []*Value) error {
	for i, node := range nodes {
		input := inputs[i]
		if err := l.convertValue(node, input, values); err != nil {
			return err
		}

		if err := l.validateOutputType(node, input); err != nil {
			return err
		}
	}

	return nil
}

func (l *Linker) convertValue(node *indexedNode, inputType reflect.Type, values []*Value) error {
	if !node.ValueOk {
		return nil
	}

	if varIndex, ok := l.isVariable(node.RawRpnNode.Name); ok {
		values[node.index] = &Value{
			ValueOk: true,
			Value: Variable{
				Index: varIndex,
				Type:  inputType,
			},
		}
		return nil
	}

	switch inputType.Kind() {
	case reflect.Int:
		integer, _ := strconv.Atoi(node.RawRpnNode.Name)
		values[node.index] = &Value{
			Value: integer,
		}
	}

	values[node.index].ValueOk = true
	return nil
}

func (l *Linker) validateOutputType(node *indexedNode, expectedType reflect.Type) error {
	if node.RawRpnNode.ValueOk {
		return nil
	}

	name := node.RawRpnNode.Name
	outputType := l.funcs[name].Outputs[0]

	if outputType != expectedType {
		return fmt.Errorf("%s returns '%v', but should return '%v'", name, outputType, expectedType)
	}

	return nil
}

func (l *Linker) isVariable(value string) (int, bool) {
	if !variableRegexp.MatchString(value) {
		return 0, false
	}

	i, err := strconv.Atoi(value[1:])
	if err != nil {
		return 0, false
	}

	return i, true
}
