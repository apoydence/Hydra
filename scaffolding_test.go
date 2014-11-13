package hydra

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scaffolding", func() {
	Context("Integrate", func() {
		It("with a linear path", func(done Done) {
			defer close(done)
			results := make(chan HashedData)
			wrapperConsumer := func(s SetupFunction) {
				consumer(s, results)
			}
			go setupScaffolding()(producer, filter, wrapperConsumer)

			expectedIndex := 0
			expectedData := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
			for data := range results {
				Expect(expectedData[expectedIndex]).To(Equal(data.Data().(int)))
				expectedIndex++
			}
		}, 1)

		PIt("with a non-linear path", func(done Done) {
			defer close(done)

			results1 := make(chan HashedData)
			wrapperConsumer1 := func(s SetupFunction) {
				consumer(s, results1)
			}

			results2 := make(chan HashedData)
			wrapperConsumer2 := func(s SetupFunction) {
				consumer2(s, results2)
			}
			go setupScaffolding()(producer, filter, filter2, wrapperConsumer1, wrapperConsumer2)

			go func() {
				expectedIndex := 0
				expectedData := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
				for data := range results1 {
					Expect(expectedData[expectedIndex]).To(Equal(data.Data().(int)))
					expectedIndex++
				}
			}()

			expectedIndex := 0
			expectedData := [...]int{0, 2, 4, 6, 8}
			for data := range results2 {
				Expect(expectedData[expectedIndex]).To(Equal(data.Data().(int)))
				expectedIndex++
			}
		}, 1)
	})
})

func producer(s SetupFunction) {
	out := s.AsProducer(1)
	defer close(out)
	for i := 0; i < 10; i++ {
		out <- NewHashedData(i, i)
	}
}

func filter(s SetupFunction) {
	in, out := s.AsFilter("github.com/apoydence/hydra.producer", 1)
	defer close(out)

	for data := range in {
		out <- data
	}
}

func filter2(s SetupFunction) {
	in, out := s.AsFilter("github.com/apoydence/hydra.producer", 1)
	defer close(out)

	for data := range in {
		if data.Hash()%2 == 0 {
			out <- data
		}
	}
}

func consumer(s SetupFunction, results WriteOnlyChannel) {
	defer close(results)
	in := s.AsConsumer("github.com/apoydence/hydra.filter", 1)
	for data := range in {
		results <- data
	}
}

func consumer2(s SetupFunction, results WriteOnlyChannel) {
	defer close(results)
	in := s.AsConsumer("github.com/apoydence/hydra.filter2", 1)
	for data := range in {
		results <- data
	}
}
