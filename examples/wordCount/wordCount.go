package main

import(
	. "github.com/apoydence/hydra/examples/wordCount/types"
	"github.com/apoydence/hydra"
	"github.com/apoydence/hydra/types"
	"os"
	"fmt"
	"bufio"
	"strings"
)

func main(){
	args := os.Args
	if len(args) == 0{
		fmt.Printf("usage: %s [Paths...]", args[0])
		os.Exit(1)
	}

	producer := func(sf types.SetupFunction){
		pathProducer(sf, args[1:])
	}

	done := make(chan struct{})
	cp := func(sf types.SetupFunction){
		wordPrinter(sf, done)
	}

	hydra.NewSetupScaffolding()(producer, pathValidator, wordExtractor, cp, puncRemover, wordCounter)

	<- done
}

func pathProducer(sf types.SetupFunction, argv []string){
	out := sf.SetName("pathProducer").AsProducer().Build()
	defer close(out)
	
	for _, path := range argv{
		out <- NewStringMarshaler(path)
	}
}

func pathValidator(sf types.SetupFunction){
	in, out := sf.AsFilter("pathProducer").Build()
	defer close(out)
	
	for path := range in{
		if _, err := os.Stat(ToString(path)); err == nil{
			out <- path
		} 
	}
}

func wordExtractor(sf types.SetupFunction){
	in, out := sf.AsFilter("main.pathValidator").Build()
	defer close(out)

	for path := range in{
		p := ToString(path)
		f, err := os.Open(p)
		if err == nil{
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanWords)
			for scanner.Scan(){
				out <- NewStringMarshaler(scanner.Text())
			}
		}
		
	}	
}

func puncRemover(sf types.SetupFunction){
	in, out := sf.AsFilter("main.wordExtractor").Build()
	defer close(out)

	replacer := strings.NewReplacer(replaceString(".", ",", "?", "!", ";", ":", "'", "\"", "$", "(", ")", "-", "/", "\\", "`", "[", "]", "{", "}", "<", ">", "*", "+", "%", "^", "&", "#", "=")...)
	
	for word := range in{
		w := strings.ToLower(ToString(word))
		if len(w) > 0{
			out <- NewStringMarshaler(replacer.Replace(w))
		}
	}
}

func replaceString(punc ...string) []string{
	results := make([]string, 0)
	for _, p := range punc{
		results = append(results, p, "")
	}
	return results
}

func wordCounter(sf types.SetupFunction){
	in, out := sf.AsFilter("main.puncRemover").Build()
	defer close(out)
	m := make(map[string]uint32)

	for word := range in{
		incMap(ToString(word), 1, m)
	}

	out <- NewWordCountMarshaler(m)
}

func wordPrinter(sf types.SetupFunction, done chan struct{}) {
	defer close(done)
	in := sf.AsConsumer("main.wordCounter").Build()
	m := make(map[string]uint32)
	
	for wordMap:= range in{
		for k, v := range ToMap(wordMap){
			incMap(k, v, m)
		}
	}

	for k, v := range m{
		println(k, v)
	}
}

func incMap(word string, incBy uint32, m map[string]uint32){
	var i uint32
	var ok bool
	if i, ok = m[word]; !ok{
		i = 0
	}
	i = i + incBy
	m[word] = i
}

