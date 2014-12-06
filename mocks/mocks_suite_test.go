package mocks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mocks Suite")
}
