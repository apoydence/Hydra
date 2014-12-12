package webCrawler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWebCrawler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebCrawler Suite")
}
