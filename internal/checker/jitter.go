package checker

import (
	"math/rand"
	"time"
)

// jitteredInterval scales interval by a random factor in
// [1-jitter, 1+jitter]. A jitter <= 0 disables jittering and returns the
// interval unchanged.
func jitteredInterval(interval time.Duration, jitter float64) time.Duration {
	if jitter <= 0 {
		return interval
	}
	factor := 1 + (rand.Float64()*2-1)*jitter
	return time.Duration(float64(interval) * factor)
}

// startupJitter returns a random duration in [0, interval). Sleeping it once
// at process start shifts this process's polling phase relative to other
// processes launched at the same wall-clock moment.
func startupJitter(interval time.Duration) time.Duration {
	if interval <= 0 {
		return 0
	}
	return time.Duration(rand.Float64() * float64(interval))
}

// sleepWithJitter sleeps for one jittered polling round. It is the shared
// sleep point for all check modes, so jitter applies uniformly.
func sleepWithJitter(interval time.Duration, jitter float64) {
	time.Sleep(jitteredInterval(interval, jitter))
}
