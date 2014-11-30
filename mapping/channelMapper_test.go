package mapping_test

import (
	. "github.com/apoydence/hydra/mapping"
	. "github.com/apoydence/hydra/testing_helpers"
	. "github.com/apoydence/hydra/types"
	"encoding"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChannelMapper", func() {
	var fakeChanCreator ChannelCreator = func(bufferSize int) chan encoding.BinaryMarshaler{
		return make(chan encoding.BinaryMarshaler, bufferSize)
	}
	
	Context("With multiple instances", func() {

		setupTest := func(numOfIns, numOfOuts int) []float64 {
			ins := make(chan WriteOnlyChannel)
			outs := make(chan ReadOnlyChannel)

			fa := func(sf SetupFunction) {
				ins <- sf.AsProducer(numOfIns)
			}

			fb := func(sf SetupFunction) {
				outs <- sf.AsConsumer("", numOfOuts)
			}

			fca := make(chan FunctionInfo)
			fcb := make(chan FunctionInfo)
			setupFa := NewSetupFunctionBuilder("a", fa, fca)
			setupFb := NewSetupFunctionBuilder("b", fb, fcb)

			for i := 0; i < numOfIns; i++ {
				go fa(setupFa)
			}

			for i := 0; i < numOfOuts; i++ {
				go fb(setupFb)
			}

			a := fetchFunctionInfos(fca, numOfIns)
			b := fetchFunctionInfos(fcb, numOfOuts)

			distMap := NewDistributedMap()
			distMap.Add("a", a, createSlice("b"))
			distMap.Add("b", b, createSlice())

			NewChannelMapper(fakeChanCreator)(distMap)

			setupInLoads(numOfIns, 100, ins)

			return channelLoad(100, numOfIns, numOfOuts, outs)
		}

		Context("Single consumer type", func() {
			It("Same number of producers as consumers", func(done Done) {
				defer close(done)

				loads := setupTest(5, 5)
				Expect(approximate(loads[0], loads[1], .1)).To(BeTrue())

			}, 1)

			It("More producers than consumers", func(done Done) {
				defer close(done)

				loads := setupTest(5, 2)
				Expect(approximate(loads[0], loads[1], .1)).To(BeTrue())

			}, 1)

			It("Less producers than consumers", func(done Done) {
				defer close(done)

				loads := setupTest(1, 2)
				Expect(approximate(loads[0], loads[1], .1)).To(BeTrue())

			}, 1)
		})
		Context("Multiple consumer types", func() {
			setupTest := func(numOfIns, numOfOuts int) []float64 {
				ins := make(chan WriteOnlyChannel)
				outsB := make(chan ReadOnlyChannel)
				outsC := make(chan ReadOnlyChannel)

				fa := func(sf SetupFunction) {
					ins <- sf.AsProducer(numOfIns)
				}

				fb := func(sf SetupFunction) {
					outsB <- sf.AsConsumer("", numOfOuts)
				}

				fc := func(sf SetupFunction) {
					outsC <- sf.AsConsumer("", numOfOuts)
				}

				fca := make(chan FunctionInfo)
				fcb := make(chan FunctionInfo)
				fcc := make(chan FunctionInfo)
				setupFa := NewSetupFunctionBuilder("a", fa, fca)
				setupFb := NewSetupFunctionBuilder("b", fb, fcb)
				setupFc := NewSetupFunctionBuilder("c", fc, fcc)

				for i := 0; i < numOfIns; i++ {
					go fa(setupFa)
				}

				for i := 0; i < numOfOuts; i++ {
					go fb(setupFb)
					go fc(setupFc)
				}

				a := fetchFunctionInfos(fca, numOfIns)
				b := fetchFunctionInfos(fcb, numOfOuts)
				c := fetchFunctionInfos(fcc, numOfOuts)

				distMap := NewDistributedMap()

				distMap.Add("a", a, createSlice("b", "c"))
				distMap.Add("b", b, createSlice())
				distMap.Add("c", c, createSlice())

				NewChannelMapper(fakeChanCreator)(distMap)

				loadsCh := make(chan []float64)
				setupInLoads(numOfIns, 100, ins)

				fetchLoads := func(loadsCh chan []float64, outs chan ReadOnlyChannel) {
					loadsCh <- channelLoad(100, numOfIns, numOfOuts, outs)
				}

				go fetchLoads(loadsCh, outsB)
				go fetchLoads(loadsCh, outsC)

				loads := make([]float64, 0)
				for i := 0; i < 2; i++ {
					load := <-loadsCh
					for _, l := range load {
						loads = append(loads, l)
					}
				}

				return loads
			}

			It("Same number of producers as consumers", func(done Done) {
				defer close(done)

				loads := setupTest(5, 5)
				for _, l := range loads {
					Expect(approximate(l, 1, .1)).To(BeTrue())
				}
			})

			It("More producers than consumers", func(done Done) {
				defer close(done)

				loads := setupTest(10, 5)
				for _, l := range loads {
					Expect(approximate(l, 2, .1)).To(BeTrue())
				}
			})

			It("Less producers than consumers", func(done Done) {
				defer close(done)

				loads := setupTest(5, 10)
				for _, l := range loads {
					Expect(approximate(l, 0.5, .1)).To(BeTrue())
				}
			})
		})
	})
})

func fetchFunctionInfos(c chan FunctionInfo, count int) []FunctionInfo {
	defer close(c)
	result := make([]FunctionInfo, 0)
	for i := 0; i < count; i++ {
		result = append(result, <-c)
	}
	return result
}

func approximate(value, expected, plusOrMinus float64) bool {
	return value <= expected+plusOrMinus &&
		value >= expected-plusOrMinus
}

func createSlice(names ...string) []string {
	result := make([]string, 0)
	for _, n := range names {
		result = append(result, n)
	}
	return result
}

func channelLoad(count, insCount, outsCount int, outs chan ReadOnlyChannel) []float64 {
	loads := make([]float64, 0)

	outSlice := make([]ReadOnlyChannel, 0)
	for i := 0; i < outsCount; i++ {
		outSlice = append(outSlice, <-outs)
		loads = append(loads, 0)
	}

	for i := 0; i < insCount*count; i++ {
		<-outSlice[i%outsCount]
		loads[i%outsCount]++
	}

	for i, v := range loads {
		loads[i] = v / float64(count)
	}

	return loads
}

func setupInLoads(insCount, count int, ins chan WriteOnlyChannel) {
	for i := 0; i < insCount; i++ {
		in := <-ins
		go func(in WriteOnlyChannel) {
			defer close(in)
			for i := 0; i < count; i++ {
				in <- NewIntMarshaler(i)
			}
		}(in)
	}
}
