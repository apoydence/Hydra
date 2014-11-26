package functionHandlers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFunctionHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FunctionHandlers Suite")
}
