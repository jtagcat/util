package wakeup_test

import (
	"context"
	"testing"
	"time"

	"github.com/jtagcat/util/wakeup"
	"github.com/stretchr/testify/assert"
)

const testSleep = 10 * time.Millisecond

func TestWakeup(t *testing.T) {
	testCounter, testState := 0, false
	waitProceed, mayProceed := false, make(chan struct{})

	wctx, wake := wakeup.WithWakeup(context.Background())
	go wakeup.Wait(wctx, func() (goToSleep bool) {
		testCounter++

		if waitProceed {
			waitProceed = false
			<-mayProceed
		}

		if testState {
			testState = false
			return true
		}
		testState = true
		return false
	})

	// Test initial run

	time.Sleep(testSleep)
	// ran 2 times before sleeping
	assert.Equal(t, 2, testCounter)

	// Test single Wakeup

	wake.Wakeup()
	time.Sleep(testSleep)
	// ran 2 times before sleeping
	assert.Equal(t, 4, testCounter)

	// Test multiple Wakeups

	waitProceed = true
	wake.Wakeup() // longer-running

	wake.Wakeup() // queue next wakeup
	wake.Wakeup() // this should be ignored, as wakeup is already queued

	mayProceed <- struct{}{}
	// ran 4 (not 6) times before sleeping
	time.Sleep(testSleep)
	assert.Equal(t, 8, testCounter)
}
