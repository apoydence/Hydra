package hydra

import (
	"reflect"
	"unsafe"

	"github.com/apoydence/hydra/internal/registrar"
	"github.com/apoydence/hydra/internal/rpn"
)

type Hydra struct {
	registrar *registrar.Registrar
	parser    *rpn.Parser
}

type Invoker func(function interface{}) func(args []unsafe.Pointer) []unsafe.Pointer

func New() *Hydra {
	return &Hydra{
		registrar: registrar.New(),
	}
}

func (h *Hydra) RegisterInvoker(inputs, outputs []reflect.Type, invoker Invoker) bool {
	return h.registrar.RegisterInvoker(inputs, outputs, registrar.Invoker(invoker))
}

func (h *Hydra) Register(name string, function interface{}) {
	h.registrar.Register(name, function)
}

func (h *Hydra) BuildFunction(dsl string) (func(args []unsafe.Pointer) []unsafe.Pointer, error) {
	rpnNodes, err := h.parser.Parse(dsl)
	if err != nil {
		return nil, err
	}

	funcs := h.registrar.Funcs()
	linker := rpn.New(funcs)
	values, err := linker.Link(rpnNodes)
	if err != nil {
		return nil, err
	}

	invoker := rpn.NewInvoker(values...)

	return h.buildFunc(invoker), nil
}

func (h *Hydra) buildFunc(invoker *rpn.Invoker) func([]unsafe.Pointer) []unsafe.Pointer {
	return func(args []unsafe.Pointer) []unsafe.Pointer {
		if len(args) != 1 {
			panic("Only supports a single argument")
		}

		output := invoker.Invoke(args[0])
		return []unsafe.Pointer{output}
	}
}
