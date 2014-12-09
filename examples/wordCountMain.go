package main

import (
	"fmt"
	"github.com/apoydence/hydra"
	"github.com/apoydence/hydra/examples/wordCount"
	"github.com/apoydence/hydra/types"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 0 {
		fmt.Printf("usage: %s [Paths...]\n", args[0])
		os.Exit(1)
	}

	producer := func(sf types.SetupFunction) {
		wordCount.PathProducer(sf, args[1:])
	}

	done := make(chan struct{})
	cp := func(sf types.SetupFunction) {
		wordCount.WordPrinter(sf, done)
	}

	hydra.NewSetupScaffolding()(producer, wordCount.PathValidator, wordCount.WordExtractor, cp, wordCount.SymbolRemover, wordCount.WordCounter)

	<-done
}
