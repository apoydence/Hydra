package textDownloader_test

import (
	"encoding"
	. "github.com/apoydence/hydra/examples/textDownloader"
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/mocks"
	"github.com/apoydence/hydra/types"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextDownloader/TextDownloader", func() {
	var in chan encoding.BinaryMarshaler
	var out chan encoding.BinaryMarshaler
	var sf types.SetupFunction

	BeforeEach(func() {
		in = make(chan encoding.BinaryMarshaler)
		out = make(chan encoding.BinaryMarshaler)
		sf = mocks.NewMockSetupFunction(in, out)
	})

	Context("UrlProducer", func() {
		It("should only output new URLs", func(done Done) {
			defer close(done)
			go UrlProducer(sf, "parent")

			in <- NewStringMarshaler("url-a")
			Expect(ToString(<-out)).To(Equal("url-a"))

			in <- NewStringMarshaler("url-b")
			Expect(ToString(<-out)).To(Equal("url-b"))

			in <- NewStringMarshaler("url-a")
			Expect(timedReceive(out)).To(BeNil())
		})
	})

	Context("UrlParser", func() {
		var server *httptest.Server
		BeforeEach(func() {
			var handler serverHandler = func(w http.ResponseWriter, req *http.Request) {
				io.WriteString(w, "<html><head></head><body><a href='a'/><a href='b'/><a href='c'/></body></html>")
			}
			server = httptest.NewServer(handler)
		})

		AfterEach(func() {
			server.Close()
		})

		It("should output the links of the webpage", func(done Done) {
			defer close(done)
			go func() {
				defer GinkgoRecover()
				UrlParser(sf)
			}()
			in <- NewStringMarshaler(server.URL + "/somepage.html")
			Expect(ToString(<-out)).To(Equal(server.URL + "/a"))
			Expect(ToString(<-out)).To(Equal(server.URL + "/b"))
			Expect(ToString(<-out)).To(Equal(server.URL + "/c"))
			in <- NewStringMarshaler(server.URL + "/somefolder")
			Expect(ToString(<-out)).To(Equal(server.URL + "/somefolder/a"))
			Expect(ToString(<-out)).To(Equal(server.URL + "/somefolder/b"))
			Expect(ToString(<-out)).To(Equal(server.URL + "/somefolder/c"))
		})

		It("should timeout with bad URLs", func(done Done) {
			defer close(done)
			go func() {
				defer close(in)
				in <- NewStringMarshaler("BAD")
			}()
			UrlParser(sf)
		}, 3)
	})

	Context("MimeDetector", func() {
		var server *httptest.Server
		BeforeEach(func() {
			var handler serverHandler = func(w http.ResponseWriter, req *http.Request) {
				if req.URL.String() == "/html" {
					io.WriteString(w, "<!doctype html><html itemscope='' itemtype='http://schema.org/WebPage' lang='en'><head><meta content='Search the world's information, including webpages, images, videos and more. Google has many special features to help you find exactly what you're looking for.' name='description'><meta content='noodp' name='robots'><meta content='/images/google_favicon_128.png' itemprop='image'><title>Google</title><script>(function(){window.google={kEI:'8GSEVLzqJNbdoATki4CwDg',kEXPI:'4011559,4013920,4016824,4017578,402034")
				} else {
					io.WriteString(w, "some text")
				}
			}
			server = httptest.NewServer(handler)
		})

		AfterEach(func() {
			server.Close()
		})
		It("should encode the mime type at the beginning of the url", func(done Done) {
			defer close(done)
			go MimeDetector(sf)
			in <- NewStringMarshaler(server.URL + "/html")
			Expect(ToString(<-out)).To(Equal("text/html; charset=utf-8->" + server.URL + "/html"))
			in <- NewStringMarshaler(server.URL + "/text")
			Expect(ToString(<-out)).To(Equal("text/plain; charset=utf-8->" + server.URL + "/text"))
		})
	})

	Context("MimeSplitterHtml", func() {
		It("should only return the html URLs", func(done Done) {
			defer close(done)
			go MimeSplitterHtml(sf)
			in <- NewStringMarshaler("image/png->https://www.google.com/images/srpr/logo11w.png")
			in <- NewStringMarshaler("text/plain->https://www.textfiles.com/some.txt")
			in <- NewStringMarshaler("text/html; charset=utf-8->http://google.com")
			Expect(ToString(<-out)).To(Equal("http://google.com"))
		})
	})

	Context("MimeSplitterText", func() {
		It("should only return the text/plain URLs", func(done Done) {
			defer close(done)
			go MimeSplitterText(sf)
			in <- NewStringMarshaler("image/png->https://www.google.com/images/srpr/logo11w.png")
			in <- NewStringMarshaler("text/html; charset=utf-8->http://google.com")
			in <- NewStringMarshaler("text/plain->http://www.textfiles.com/some.txt")
			Expect(ToString(<-out)).To(Equal("http://www.textfiles.com/some.txt"))
		})
	})
})

func timedReceive(c chan encoding.BinaryMarshaler) encoding.BinaryMarshaler {
	t := time.NewTimer(time.Millisecond * 500).C
	select {
	case <-t:
		return nil
	case result := <-c:
		return result
	}
}

type serverHandler func(w http.ResponseWriter, req *http.Request)

func (sh serverHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sh(w, req)
}
