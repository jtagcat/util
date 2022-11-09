package util

import "sync"

// go <func> with single waitgroup
//
//	doGoWait= util.GoWg(func() {
//		exampleFunc(foo, bar)
//	})
//	defer done()
func GoWg(fn func()) (done func()) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()

	return wg.Done
}
