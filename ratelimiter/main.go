package ratelimiter

import (
	"sync"
	"time"
)

var ratelimiters map[string]*ratelimiter
var mutex sync.Mutex

type ratelimiter struct {
	ticker      *time.Ticker
	subscribers int
}

func Get(key string, limit int) <-chan time.Time {
	mutex.Lock()
	defer mutex.Unlock()
	if ratelimiters == nil {
		ratelimiters = make(map[string]*ratelimiter)
	}
	rl := ratelimiters[key]
	if rl == nil {
		rl = &ratelimiter{
			ticker: time.NewTicker(time.Second / time.Duration(limit)),
		}
		ratelimiters[key] = rl
	}
	rl.subscribers++
	return rl.ticker.C
}

func Stop(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	if ratelimiters == nil {
		return
	}
	rl := ratelimiters[key]
	if rl == nil {
		return
	}
	rl.subscribers--
	if rl.subscribers <= 0 {
		rl.ticker.Stop()
		delete(ratelimiters, key)
	}
}
