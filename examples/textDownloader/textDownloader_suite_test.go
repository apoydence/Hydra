package textDownloader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTextDownloader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TextDownloader Suite")
}
