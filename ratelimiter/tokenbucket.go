package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate 	 time.Duration
	stopRefiller chan struct{}
	mu           sync.Mutex
}

func NewTokenBucket(capacity, tokensPerInterval int, refillRate time.Duration) *TokenBucket {
	bucket := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		stopRefiller: make(chan struct{}),
	}

	go bucket.refillTokens(tokensPerInterval)
	return bucket
}

func (t *TokenBucket) refillTokens(tokensPerInterval int) {
	ticker := time.NewTicker(t.refillRate)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			t.mu.Lock()
			t.tokens = min(t.capacity, t.tokens+tokensPerInterval)
			t.mu.Unlock()
		case <- t.stopRefiller:
			return
		}
	}
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

func (t *TokenBucket) StopRefiller() {
	close(t.stopRefiller)
}

func main() {
	t := NewTokenBucket(5, 2, time.Second)
	for i:=1;i<=10;i++ {
		if t.Allow() {
			fmt.Printf("Token Taken. Remaining tokens: %d\n", t.tokens)
		} else {
			fmt.Printf("Not enough tokens.\n")
		}
		time.Sleep(150 * time.Millisecond)
	}
	t.StopRefiller()
}