package wordCount

import (
	"bufio"
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra/types"
	"os"
	"strings"
)

func PathProducer(sf types.SetupFunction, argv []string) {
	out := sf.SetName("PathProducer").AsProducer().Build()
	defer close(out)

	for _, path := range argv {
		out <- NewStringMarshaler(path)
	}
}

func PathValidator(sf types.SetupFunction) {
	in, out := sf.AsFilter("PathProducer").Build()
	defer close(out)

	for path := range in {
		if _, err := os.Stat(ToString(path)); err == nil {
			out <- path
		}
	}
}

func WordExtractor(sf types.SetupFunction) {
	in, out := sf.AsFilter("github.com/apoydence/hydra/examples/wordCount.PathValidator").Build()
	defer close(out)

	for path := range in {
		p := ToString(path)
		f, err := os.Open(p)
		if err == nil {
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				out <- NewStringMarshaler(scanner.Text())
			}
		}

	}
}

func SymbolRemover(sf types.SetupFunction) {
	in, out := sf.AsFilter("github.com/apoydence/hydra/examples/wordCount.WordExtractor").Build()
	defer close(out)

	for word := range in {
		str := make([]byte, 0)
		bytes := []byte(strings.ToLower(ToString(word)))
		for _, x := range bytes {
			if (x >= 0x30 && x <= 0x39) || (x >= 0x61 && x <= 0x7a) {
				str = append(str, x)
			}
		}
		if len(str) > 0 {
			out <- NewStringMarshaler(string(str))
		}
	}
}

func WordCounter(sf types.SetupFunction) {
	in, out := sf.AsFilter("github.com/apoydence/hydra/examples/wordCount.SymbolRemover").Build()
	defer close(out)
	m := make(map[string]uint32)

	for word := range in {
		incMap(ToString(word), 1, m)
	}

	out <- NewWordCountMarshaler(m)
}

func WordPrinter(sf types.SetupFunction, done chan struct{}) {
	defer close(done)
	in := sf.AsConsumer("github.com/apoydence/hydra/examples/wordCount.WordCounter").Build()
	m := make(map[string]uint32)

	for wordMap := range in {
		for k, v := range ToMap(wordMap) {
			incMap(k, v, m)
		}
	}

	for k, v := range m {
		println(k, v)
	}
}

func incMap(word string, incBy uint32, m map[string]uint32) {
	var i uint32
	var ok bool
	if i, ok = m[word]; !ok {
		i = 0
	}
	i = i + incBy
	m[word] = i
}
