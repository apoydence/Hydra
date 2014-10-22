package hydra

import (
	"reflect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHydra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hydra Suite")
}

func samePointers(a interface{}, b interface{}) bool{
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
