package internal

import (
	"math/rand"
	"time"
)

func init() {
	// It's such a low impact use case, there's no need for crypto/rand.
	// Famous last words...
	rand.Seed(time.Now().Unix())
}

const millis int = 1000
const minRate int = 500
const defaultJitter int = 2

// RateLimit accepts a jitter value (in seconds), uses it to create a new random
// timeout (in milliseconds) and adds a minimum rate (500ms).
// If jitter is less than one second, a default value will be used (2s).
func RateLimit(jitter int) time.Duration {
	if jitter < 1 {
		jitter = defaultJitter
	}
	i := rand.Intn(jitter*millis) + minRate
	return time.Duration(i) * time.Millisecond
}
