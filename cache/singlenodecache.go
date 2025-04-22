package main

import (
	"container/list"
	"context"
	"sync"
	"time"
	"fmt"
)

type item struct {
	key       string
	value     string
	expiresAt time.Time
}

type LRUCache struct {
	capacity  int
	evictList *list.List
	items     map[string]*list.Element
	ttl       time.Duration
	mu        sync.RWMutex
}

func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
		ttl:       ttl,
	}
}

func (l *LRUCache) Set(key, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.items[key]; ok {
		l.evictList.MoveToFront(el)
		el.Value.(*item).value = value
		el.Value.(*item).expiresAt = time.Now().Add(l.ttl)
	}

	if l.evictList.Len() >= l.capacity {
		l.evict()
	}

	item := &item{key: key, value: value, expiresAt: time.Now().Add(l.ttl),}
	element := l.evictList.PushFront(item)
	l.items[key] = element
}

func (l *LRUCache) Get(key string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if el,ok := l.items[key]; ok {
		if time.Now().After(el.Value.(*item).expiresAt) {
			l.mu.RUnlock()
			l.mu.RLock()
			l.removeElement(el)
			l.mu.Unlock()
			l.mu.RLock()
			return "", false
		}
		l.evictList.MoveToFront(el)
		return el.Value.(*item).value, true
	}
	return "", false
}

func (l *LRUCache) evict() {
	el := l.evictList.Back()
	// To prevent tampering or unsafe manual modification
	if el != nil {
		l.removeElement(el)
	}
}

func (l *LRUCache) removeElement(el *list.Element) {
	delete(l.items, el.Value.(*item).key)
	l.evictList.Remove(el)
}

// Delete any expired keys every interval
func (l *LRUCache) TTLCollector(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.mu.Lock()
				for _,el := range l.items {
					if time.Now().After(el.Value.(*item).expiresAt) {
						l.removeElement(el)
					}
				}
				l.mu.Unlock()
			}
		}
	}()
}

func main() {
	ctx,cancel := context.WithCancel(context.Background())
	cache := NewLRUCache(3, 5*time.Second)
	cache.TTLCollector(ctx, 1*time.Second)

	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3")
	fmt.Println(cache.Get("a"))
	cache.Set("d", "4")
	fmt.Println(cache.Get("c")) // should be available
	fmt.Println(cache.Get("b")) // should be evicted due to LRU
	time.Sleep(6 * time.Second)
	fmt.Println(cache.Get("a")) // expired

	cancel()
}