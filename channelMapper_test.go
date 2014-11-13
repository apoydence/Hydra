package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChannelMapper", func() {
	Context("With multiple instances", func() {

		setupTest := func(numOfIns, numOfOuts int) []float64{
			ins := make(chan WriteOnlyChannel)
			outs := make(chan ReadOnlyChannel)

			fa := func(sf SetupFunction){
				ins <- sf.AsProducer(numOfIns)
			}

			fb := func(sf SetupFunction){
				outs <- sf.AsConsumer("", numOfOuts)
			}

			fca := make(chan FunctionInfo)
			fcb := make(chan FunctionInfo)
			setupFa := buildSetupFunc("a", fa, fca)
			setupFb := buildSetupFunc("b", fb, fcb)

			for i:=0; i<numOfIns; i++{
				go fa(setupFa)
			}

			for i:=0; i<numOfOuts; i++{
				go fb(setupFb)
			}

			a := fetchFunctionInfos(fca, numOfIns)
			b := fetchFunctionInfos(fcb, numOfOuts)

			m := make(map[string]*distMapper)

			da := newDistMapper(a, createSlice("b"))
			m["a"] = da
			db := newDistMapper(b, createSlice())
			m["b"] = db

			var distMap distFunctionMap
			distMap = m
			channelMapper(distMap)

			return channelLoad(numOfIns, ins, numOfOuts, outs)
		}

		Context("Same instances in and out", func() {
			It("Same ins and outs", func(done Done) {
				defer close(done)

				loads := setupTest(5, 5)
				Expect(approximate(loads[0], 1, .1)).To(BeTrue())
				Expect(approximate(loads[1], 1, .1)).To(BeTrue())
	
			}, 1)

			It("More ins than outs", func(done Done) {
				defer close(done)

				loads := setupTest(7, 5)
				Expect(approximate(loads[0], 1, .1)).To(BeTrue())
				Expect(approximate(loads[1], 1, .1)).To(BeTrue())
	
			}, 1)
		})
		Context("More ins than outs", func() {
			It("Should evenly distribute", func() {
			})
		})
		Context("Less ins than outs", func() {
			It("Should evenly distribute", func() {
			})
		})
	})
})

func fetchFunctionInfos(c chan FunctionInfo, count int) []FunctionInfo{
	defer close(c)
	result := make([]FunctionInfo, 0)
	for i:=0; i<count; i++{
		result = append(result, <-c)
	}
	return result
}

func approximate(value, expected, plusOrMinus float64) bool{
	return 	value <= expected + plusOrMinus &&
		value >= expected - plusOrMinus;
}

func createSlice(names ...string)[]string{
	result := make([]string, 0)
	for _, n := range names{
		result = append(result, n)
	}
	return result
}

func channelLoad(insCount int, ins chan WriteOnlyChannel, outsCount int, outs chan ReadOnlyChannel)[]float64{
	for i:=0; i<insCount; i++{
		in := <- ins
		go func(){
			defer close(in)
			for i:=0; i<100; i++{
				in <- NewHashedData(i, i);
			}
		}()
	}

	loadCh := make(chan float64)
	loads := make([]float64, 0)

	for i:=0; i<outsCount; i++{
		out := <- outs
		go func(){
			count := 0
			for _ = range out{
				count++
			}
			loadCh <- float64(count) / 100.0
		}()
	}

	for i:=0; i<outsCount; i++{
		loads = append(loads, <-loadCh)
	}
	return loads
}
