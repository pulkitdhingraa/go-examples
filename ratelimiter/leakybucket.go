package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Packet struct {
	id   int
	size int
}

func NewPacket(id, size int) *Packet {
	return &Packet{
		id:   id,
		size: size,
	}
}

type LeakyBucket struct {
	capacity          int
	leakAmountPerTick int
	leakRate          time.Duration
	buffer            []Packet
	curBufferSize     int
	stopCh            chan struct{}
	mu                sync.Mutex
}

func NewLeakyBucket(cap, leakAmountPerTick int, leakRate time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:          cap,
		leakAmountPerTick: leakAmountPerTick,
		leakRate:          leakRate,
		buffer:            []Packet{},
		stopCh:            make(chan struct{}),
		curBufferSize:     0,
	}
	go lb.transmitTick()
	return lb
}

func (lb *LeakyBucket) AddPacket(p Packet) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.curBufferSize+p.size > lb.capacity {
		fmt.Printf("Bucket is full. Packet with id %d size %d is rejected\n", p.id, p.size)
		return
	}

	lb.buffer = append(lb.buffer, p)
	lb.curBufferSize += p.size
	fmt.Printf("Packet with id %d size %d added to the bucket\n", p.id, p.size)
}

func (lb *LeakyBucket) transmitTick() {
	ticker := time.NewTicker(lb.leakRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lb.mu.Lock()
			n := lb.leakAmountPerTick
			for len(lb.buffer) > 0 {
				top := lb.buffer[0]
				if top.size > n {
					break
				}
				n -= top.size
				lb.curBufferSize -= top.size
				fmt.Printf("Packet with id %d transmitted\n", top.id)
				lb.buffer = lb.buffer[1:]
			}
			if lb.curBufferSize == 0 {
				fmt.Println("No packets in the bucket.")
			}
			lb.mu.Unlock()
		case <-lb.stopCh:
			fmt.Println("Transmission Stopped.")
			return
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	bucket := NewLeakyBucket(20, 7, 500*time.Millisecond)
	for i := 1; i <= 20; i++ {
		go func(i int){
			delay := time.Duration(rand.Intn(1500)) * time.Millisecond
			time.Sleep(delay)
			size := rand.Intn(7) + 1
			p := NewPacket(i, size)
			bucket.AddPacket(*p)
		}(i)
	}
	time.Sleep(15 * time.Second)
	close(bucket.stopCh)
}
