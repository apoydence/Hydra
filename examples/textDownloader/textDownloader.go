package textDownloader

import (
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/types"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	urlParser "net/url"
	"path"
	"strings"
	"time"
)

func UrlProducer(sf types.SetupFunction, parent string) {
	in, out := sf.SetName("UrlProducer").AsFilter(parent).Build()
	defer close(out)

	visitedUrls := make(map[string]interface{})

	for urlBm := range in {
		if sf.Cancelled() {
			continue
		}
		url := ToString(urlBm)
		if _, visited := visitedUrls[url]; !visited {
			visitedUrls[url] = nil
			out <- NewStringMarshaler(url)
		}
	}
}

func UrlParser(sf types.SetupFunction) {
	in, out := sf.SetName("UrlParser").AsFilter("UrlProducer").Build()
	defer close(out)

	for urlBM := range in {
		if sf.Cancelled() {
			continue
		}
		urlStr := ToString(urlBM)
		if path.Ext(urlStr) == "" && urlStr[len(urlStr)-1] != '/' {
			urlStr += "/"
		}
		url, err := urlParser.Parse(urlStr)
		if err != nil {
			continue
		}

		body := Download(url.String())
		for link := range fetchLinks(body) {
			joinedLink, err := url.Parse(link)
			if err != nil {
				continue
			}

			out <- NewStringMarshaler(joinedLink.String())
		}
	}
}

func MimeDetector(sf types.SetupFunction) {
	in, out := sf.SetName("MimeDetector").AsFilter("UrlParser").Build()
	defer close(out)

	buffer := make([]byte, 512)

	for urlBM := range in {
		if sf.Cancelled() {
			continue
		}
		url := ToString(urlBM)
		body := Download(url)
		if body == nil {
			continue
		}
		func(body io.ReadCloser) {
			defer body.Close()
			n, err := body.Read(buffer)

			if err == nil || err == io.EOF {
				encoded := http.DetectContentType(buffer[:n]) + "->" + url
				out <- NewStringMarshaler(encoded)
			}
		}(body)
	}
}

func MimeSplitterHtml(sf types.SetupFunction) {
	in, out := sf.SetName("MimeSplitterHtml").AsFilter("MimeDetector").Build()
	defer close(out)
	for urlBM := range in {
		if sf.Cancelled() {
			continue
		}
		url := ToString(urlBM)
		mime, u := decodeMimeUrl(url)
		if strings.Contains(mime, "html") {
			out <- NewStringMarshaler(u)
		}
	}
}

func MimeSplitterText(sf types.SetupFunction) {
	in, out := sf.SetName("MimeSplitterText").AsFilter("MimeDetector").Build()
	defer close(out)
	for urlBM := range in {
		if sf.Cancelled() {
			continue
		}
		url := ToString(urlBM)
		mime, u := decodeMimeUrl(url)
		if strings.Contains(mime, "text/plain") {
			out <- NewStringMarshaler(u)
		}
	}
}

func decodeMimeUrl(u string) (mime string, url string) {
	decoded := strings.SplitN(u, "->", 2)
	mime = decoded[0]
	url = decoded[1]
	return
}

func Download(url string) io.ReadCloser {
	u, err := urlParser.Parse(url)
	if err != nil {
		return nil
	}

	dialTimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(time.Second*2))
	}

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(u.String())
	if err == nil {
		return resp.Body
	}

	return nil
}

func fetchLinks(body io.ReadCloser) chan string {
	links := make(chan string)
	if body != nil {
		defer body.Close()
		doc, err := html.Parse(body)
		if err == nil {
			go func() {
				defer close(links)
				parseHtml(doc, links)
			}()
		}
	} else {
		close(links)
	}
	return links
}

func parseHtml(node *html.Node, links chan string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				links <- a.Val
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseHtml(c, links)
	}
}
