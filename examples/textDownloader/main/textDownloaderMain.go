package main

import (
	"fmt"
	"github.com/apoydence/hydra"
	. "github.com/apoydence/hydra/examples/textDownloader"
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/types"
	"github.com/eapache/channels"
	"io"
	"os"
	"path"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("usage: %s [URL]\n", args[0])
		os.Exit(1)
	}

	url := args[1]
	done := make(chan struct{})
	at := types.NewAtomicBool(false)
	feeder := channels.NewInfiniteChannel()

	feeder.In() <- url

	uFeeder := func(sf types.SetupFunction) {
		urlFeeder(sf, feeder.Out(), at)
	}

	producer := func(sf types.SetupFunction) {
		UrlProducer(sf, "UrlFeeder")
	}

	looper := func(sf types.SetupFunction) {
		urlLooper(sf, feeder.In(), at)
	}

	var cancel types.Canceller

	downloader := func(sf types.SetupFunction) {
		textDownloader(sf, done, func() {
			if !at.Get() {
				feeder.Close()
				cancel()
			}
			at.Set(true)
		})
	}

	cancel = hydra.NewSetupScaffolding()(uFeeder, looper, downloader, producer, UrlParser, MimeDetector, MimeSplitterHtml, MimeSplitterText)

	<-done
}

func urlFeeder(sf types.SetupFunction, feeder <-chan interface{}, done types.AtomicBool) {
	out := sf.SetName("UrlFeeder").AsProducer().Build()
	defer close(out)

	for url := range feeder {
		if !done.Get() {
			out <- NewStringMarshaler(url.(string))
		}
	}
}

func urlLooper(sf types.SetupFunction, feeder chan<- interface{}, done types.AtomicBool) {
	in := sf.AsConsumer("MimeSplitterHtml").Build()

	for url := range in {
		if !done.Get() {
			u := ToString(url)
			feeder <- u
		}
	}
}

func textDownloader(sf types.SetupFunction, done chan struct{}, closer func()) {
	in := sf.AsConsumer("MimeSplitterText").Build()

	var totalSize int64 = 0

	for url := range in {
		if totalSize >= 1*1024*1024 {
			closer()
			continue
		}

		u := ToString(url)
		println("download", u, path.Base(u))
		totalSize += saveToFile(Download(u), path.Join("/tmp/textFiles", path.Base(u)))
	}

	close(done)
}

func saveToFile(body io.ReadCloser, path string) int64 {
	if body == nil {
		return 0
	}

	defer body.Close()

	f, err := os.Create(path)

	if err != nil {
		return 0
	}

	defer f.Close()

	buffer := make([]byte, 1024)
	var count int64 = 0

	var n int

	for {
		n, err = body.Read(buffer)
		if n <= 0 || (err != nil && err != io.EOF) {
			break
		}

		count += int64(n)
		f.Write(buffer[:n])
	}

	return count
}
