package hydra_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHydra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hydra Suite")
}
