package wordCount_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWordCount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WordCount Suite")
}
