package end2end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEnd2endtests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "End2End Test Suite")
}
