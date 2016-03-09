package rpn_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRpn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rpn Suite")
}
