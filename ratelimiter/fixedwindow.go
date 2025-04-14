package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type FixedWindow struct {
	windowSize time.Duration
	maxRequest int
	requestCount int
	windowStart time.Time
	mu sync.Mutex
}

func NewFixedWindow(windowSize time.Duration, maxRequest int) *FixedWindow{
	return &FixedWindow{
		windowSize: windowSize,
		maxRequest: maxRequest,
		requestCount: 0,
		windowStart: time.Now(),
	}
}

func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// If new window reset counter
	if time.Now().Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = time.Now()
		fw.requestCount = 0
	}

	if fw.requestCount < fw.maxRequest {
		fw.requestCount++
		return true
	}
	
	return false
}

func main() {
	limiter := NewFixedWindow(2*time.Second, 5)
	var wg sync.WaitGroup
	for i:=1;i<=15;i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			if limiter.Allow() {
				fmt.Printf("Id %d processed\n", i)
			} else {
				fmt.Printf("Id %d limited\n", i)
			}
		}(i)
	}
	wg.Wait()
}