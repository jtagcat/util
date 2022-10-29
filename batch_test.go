package simple_test

import (
	"sync"
	"testing"
	"time"

	"github.com/jtagcat/simple"
	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	stream, batched := make(chan int, 8), make(chan []int, 1)

	wait, wg := make(chan struct{}), sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		expected := []struct {
			wait  bool
			value []int
		}{
			// 4 batch + 2 timeout
			{value: []int{0, 1, 2, 3}},
			{value: []int{4, 5}},

			// overfill + close
			{wait: true, value: []int{0, 1, 2, 3}}, // in batched
			{value: []int{4, 5, 6, 7}},             // in function
			{value: []int{8, 9, 10, 11}},           // in stream queue
			{value: []int{12, 13}},                 // in stream queue (total 6)
		}

		for i := 0; i < len(expected); i++ {
			want := expected[i]
			if want.wait {
				<-wait
			}

			b := <-batched
			assert.Equal(t, want.value, b)
		}
		assert.Equal(t, 0, len(batched))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		simple.Batch(4, 100*time.Millisecond, stream, batched)
	}()

	// 4 batch + 2 timeout
	for i := 0; i < 6; i++ {
		stream <- i
	}
	time.Sleep(300 * time.Millisecond)

	// overfill + close
	for i := 0; i < 14; i++ {
		stream <- i
	}
	close(stream)
	wait <- struct{}{}

	wg.Wait()
}
