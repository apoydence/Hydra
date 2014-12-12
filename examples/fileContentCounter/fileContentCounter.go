package main

import (
	"fmt"
	"github.com/apoydence/hydra"
	"github.com/apoydence/hydra/examples/wordCount"
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/types"
	"os"
	"path"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("usage: %s [Paths...]\n", path.Base(args[0]))
		os.Exit(1)
	}

	validator := func(sf types.SetupFunction) {
		wordCount.PathValidator(sf, "PathProducer")
	}

	producer := func(sf types.SetupFunction) {
		pathProducer(sf, args[1:])
	}

	done := make(chan struct{})
	cp := func(sf types.SetupFunction) {
		wordPrinter(sf, done)
	}

	hydra.NewSetupScaffolding()(producer, validator, wordCount.WordExtractor, cp, wordCount.SymbolRemover, wordCount.WordCounter, wordCount.FinalWordCounter)

	<-done
}

func pathProducer(sf types.SetupFunction, argv []string) {
	out := sf.SetName("PathProducer").AsProducer().Build()
	defer close(out)

	for _, path := range argv {
		out <- NewStringMarshaler(path)
	}
}

func wordPrinter(sf types.SetupFunction, done chan struct{}) {
	defer close(done)
	in := sf.AsConsumer("github.com/apoydence/hydra/examples/wordCount.FinalWordCounter").Build()

	for wordMap := range in {
		for k, v := range ToMap(wordMap) {
			println(k, v)
		}
	}
}
