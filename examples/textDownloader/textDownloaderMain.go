package main

import (
	"fmt"
	"github.com/apoydence/hydra"
	. "github.com/apoydence/hydra/examples/webCrawler"
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/types"
	"github.com/eapache/channels"
	"io"
	"os"
	"path"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("usage: %s [URL] [DOWNLOAD PATH] [# MBs]\n", path.Base(args[0]))
		os.Exit(1)
	}

	url := args[1]
	download := args[2]
	size, err := strconv.Atoi(args[3])

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Printf("usage: %s [URL] [DOWNLOAD PATH] [# MBs]\n", path.Base(args[0]))
	}

	done := make(chan struct{})
	at := types.NewAtomicBool(false)
	feeder := channels.NewInfiniteChannel()

	if !dirExists(download) {
		if err := os.MkdirAll(download, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory (%v): %v.\n", download, err.Error())
			os.Exit(1)
		}
	}

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
		textDownloader(sf, download, size, done, func() {
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

func textDownloader(sf types.SetupFunction, download string, mb int, done chan struct{}, closer func()) {
	in := sf.AsConsumer("MimeSplitterText").Build()

	var totalSize int64 = 0

	for url := range in {
		if totalSize >= int64(mb)*1024*1024 {
			closer()
			continue
		}

		u := ToString(url)
		println("download", u, path.Base(u))
		totalSize += saveToFile(Download(u), path.Join(download, path.Base(u)))
		println("Size", totalSize)
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

func dirExists(path string) bool {
	src, err := os.Stat(path)
	if err == nil {
		return src.IsDir()
	}
	return false
}
