package registrar

import (
	"fmt"
	"reflect"
	"regexp"
	"unsafe"

	"github.com/poy/hydra/internal/rpn"
)

type Invoker func(function interface{}) func(args []unsafe.Pointer) []unsafe.Pointer

type Registrar struct {
	funcs    map[string]rpn.Callable
	invokers []invokerInfo
}

type invokerInfo struct {
	Inputs  []reflect.Type
	Outputs []reflect.Type
	Invoker Invoker
}

var (
	funcNameReg *regexp.Regexp = regexp.MustCompile("^[A-Z][a-zA-Z0-9_]*$")
)

func New() *Registrar {
	return &Registrar{
		funcs: make(map[string]rpn.Callable),
	}
}

func (r *Registrar) RegisterInvoker(inputs, outputs []reflect.Type, invoker Invoker) bool {
	if prev := r.fetchInvoker(inputs, outputs); prev != nil {
		return false
	}

	r.invokers = append(r.invokers, invokerInfo{
		Inputs:  inputs,
		Outputs: outputs,
		Invoker: invoker,
	})

	return true
}

func (r *Registrar) Register(name string, function interface{}) {
	r.validateName(name)

	functionType := reflect.TypeOf(function)
	inputs := r.readInputs(functionType)
	outputs := r.readOutputs(functionType)
	invoker := r.fetchInvoker(inputs, outputs)
	if invoker == nil {
		panic("Invoker not registered")
	}

	r.funcs[name] = rpn.Callable{
		Inputs:   inputs,
		Outputs:  outputs,
		Function: invoker(function),
	}
}

func (r *Registrar) Funcs() map[string]rpn.Callable {
	return r.funcs
}

func (r *Registrar) readInputs(functionType reflect.Type) []reflect.Type {
	var types []reflect.Type

	for i := 0; i < functionType.NumIn(); i++ {
		types = append(types, functionType.In(i))
	}

	return types
}

func (r *Registrar) readOutputs(functionType reflect.Type) []reflect.Type {
	var types []reflect.Type

	for i := 0; i < functionType.NumOut(); i++ {
		types = append(types, functionType.Out(i))
	}

	return types
}

func (r *Registrar) validateName(name string) {
	if funcNameReg.MatchString(name) {
		return
	}

	panic(fmt.Sprintf("Invalid function name: '%s'", name))
}

func (r *Registrar) fetchInvoker(inputs, outputs []reflect.Type) Invoker {
	for _, invoker := range r.invokers {
		if reflect.DeepEqual(invoker.Inputs, inputs) && reflect.DeepEqual(invoker.Outputs, outputs) {
			return invoker.Invoker
		}
	}
	return nil
}
